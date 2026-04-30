package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/mailout"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/store"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/uploads"
)

// allowedUploadTypes is the set of MIME types candidates are permitted
// to attach. The type is detected from the file's magic bytes on the
// server side, so the client-supplied Content-Type header is ignored.
var allowedUploadTypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
}

const reapplyWindow = 7 * 24 * time.Hour

// handleSubmitApplication handles a multipart application submission
// against a job posting. The request body must be multipart/form-data
// and is bounded by Content-Length (server rejects > MaxUploadTotalBytes).
//
// Form fields:
//
//	name      (required, string)
//	email     (required, RFC 5322 address)
//	message   (optional, ≤ 4 KB)
//	files[]   (0..3 files, each ≤ MaxUploadFileBytes)
func (s *Server) handleSubmitApplication(c *gin.Context) {
	jobID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || jobID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	if c.Request.ContentLength > 0 && c.Request.ContentLength > s.MaxUploadTotalBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": fmt.Sprintf("payload exceeds %d bytes", s.MaxUploadTotalBytes),
		})
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, s.MaxUploadTotalBytes)

	// Confirm the target job exists and is currently open. Closed roles
	// shouldn't accept new applications.
	job, err := s.Jobs.Get(c.Request.Context(), jobID)
	if err != nil {
		if store.IsNoRows(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
			return
		}
		s.Logger.Error("apply: load job", "err", err, "id", jobID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if !job.IsActive {
		c.JSON(http.StatusGone, gin.H{"error": "this role is no longer accepting applications"})
		return
	}

	if err := c.Request.ParseMultipartForm(8 << 20); err != nil {
		s.Logger.Info("apply: parse multipart", "err", err, "ip", c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form"})
		return
	}

	name := strings.TrimSpace(c.Request.FormValue("name"))
	email := strings.TrimSpace(c.Request.FormValue("email"))
	message := strings.TrimSpace(c.Request.FormValue("message"))

	switch {
	case name == "" || len(name) > 200:
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required (max 200 chars)"})
		return
	case email == "" || len(email) > 320:
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	case len(message) > 4096:
		c.JSON(http.StatusBadRequest, gin.H{"error": "message must be ≤ 4 KB"})
		return
	}
	if _, err := mail.ParseAddress(email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not valid"})
		return
	}

	// Prevent the same candidate from applying to the same role within 7 days.
	dup, err := s.Applications.ExistsRecentApplication(c.Request.Context(), jobID, email, reapplyWindow)
	if err != nil {
		s.Logger.Error("apply: dedup check", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if dup {
		c.JSON(http.StatusConflict, gin.H{
			"error": "You have already applied for this role in the last 7 days. Please wait before reapplying.",
		})
		return
	}

	headers := c.Request.MultipartForm.File["files"]
	if len(headers) > s.MaxUploadFiles {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("at most %d files are allowed", s.MaxUploadFiles),
		})
		return
	}
	for _, h := range headers {
		if h.Size > s.MaxUploadFileBytes {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": fmt.Sprintf("file %q exceeds %d bytes", h.Filename, s.MaxUploadFileBytes),
			})
			return
		}
	}

	// Verify each file's real type by sniffing the first 512 bytes.
	// This is done before writing to the database so no rollback is needed
	// when a disallowed type is detected.
	for _, h := range headers {
		fh, err := h.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not read uploaded file"})
			return
		}
		sniff := make([]byte, 512)
		n, _ := io.ReadFull(fh, sniff)
		_ = fh.Close()
		detected := strings.SplitN(http.DetectContentType(sniff[:n]), ";", 2)[0]
		if !allowedUploadTypes[detected] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("file %q is not allowed; only PDF and image files (JPEG, PNG, GIF, WebP) are accepted", h.Filename),
			})
			return
		}
	}

	app := &domain.Application{
		JobID:     jobID,
		Name:      name,
		Email:     email,
		Message:   message,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}

	// Stamp the row with the current default status so it appears in
	// HR's "unread" bucket immediately. If the lookup fails we still
	// accept the submission — better to have a status-less row than
	// to drop a candidate's application on the floor.
	var defaultStatusID int64
	if def, err := s.Statuses.GetDefault(c.Request.Context()); err == nil {
		defaultStatusID = def.ID
	} else if !store.IsNoRows(err) {
		s.Logger.Warn("apply: load default status", "err", err)
	}

	id, err := s.Applications.Create(c.Request.Context(), app, defaultStatusID)
	if err != nil {
		s.Logger.Error("apply: create row", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	saved := make([]domain.ApplicationFile, 0, len(headers))
	rollback := func() {
		for _, f := range saved {
			_ = s.Uploads.Remove(f.StoredPath)
		}
		_ = s.Applications.Delete(c.Request.Context(), id)
	}
	for _, h := range headers {
		fh, err := h.Open()
		if err != nil {
			rollback()
			s.Logger.Error("apply: open part", "err", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not read uploaded file"})
			return
		}
		// Re-sniff to reconstruct a full reader; Open() is safe to call again.
		sniff := make([]byte, 512)
		n, _ := io.ReadFull(fh, sniff)
		fullReader := io.MultiReader(bytes.NewReader(sniff[:n]), fh)
		ct := strings.SplitN(http.DetectContentType(sniff[:n]), ";", 2)[0]
		rel, written, err := s.Uploads.Save(id, h.Filename, fullReader, s.MaxUploadFileBytes)
		_ = fh.Close()
		if err != nil {
			rollback()
			if errors.Is(err, uploads.ErrFileTooLarge) {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": fmt.Sprintf("file %q exceeds %d bytes", h.Filename, s.MaxUploadFileBytes),
				})
				return
			}
			s.Logger.Error("apply: save file", "err", err, "name", h.Filename)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not store file"})
			return
		}
		f := domain.ApplicationFile{
			ApplicationID: id,
			OriginalName:  h.Filename,
			StoredPath:    rel,
			ContentType:   ct,
			SizeBytes:     written,
		}
		if _, err := s.Applications.AddFile(c.Request.Context(), &f); err != nil {
			_ = s.Uploads.Remove(rel)
			rollback()
			s.Logger.Error("apply: insert file row", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		saved = append(saved, f)
	}

	app.Files = saved
	s.Logger.Info("application received",
		"app_id", id, "job_id", jobID, "files", len(saved), "ip", c.ClientIP())
	s.sendApplicationAck(name, email, job.Title)
	c.JSON(http.StatusCreated, gin.H{"id": id, "files": len(saved)})
}

// sendApplicationAck emails the candidate a receipt after a successful
// application. Failures are logged only; the HTTP response already succeeded.
func (s *Server) sendApplicationAck(candidateName, to, jobTitle string) {
	if s.Mail == nil || !s.Mail.Ready() {
		return
	}
	subject := mailout.ApplicationAckSubject(jobTitle)
	plain := mailout.ApplicationAckBody(candidateName, jobTitle, s.AppName)
	html := mailout.ApplicationAckHTML(candidateName, jobTitle, s.AppName, s.SiteURL)
	recipient := to
	go func() {
		if err := s.Mail.SendAlternative(recipient, subject, plain, html); err != nil {
			s.Logger.Error("apply: ack email", "err", err, "to", recipient)
			return
		}
		s.Logger.Info("apply: ack email sent", "to", recipient)
	}()
}

func (s *Server) handleAdminListApplications(c *gin.Context) {
	var filter store.ListFilter
	if v := strings.TrimSpace(c.Query("jobId")); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid jobId"})
			return
		}
		filter.JobID = n
	}
	if v := strings.TrimSpace(c.Query("statusId")); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid statusId"})
			return
		}
		filter.StatusID = n
	}
	filter.Sort = strings.TrimSpace(c.Query("sort"))
	if v := strings.TrimSpace(c.Query("limit")); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
		filter.Limit = n
	}

	items, err := s.Applications.List(c.Request.Context(), filter)
	if err != nil {
		s.Logger.Error("admin list applications", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (s *Server) handleAdminGetApplication(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	app, err := s.Applications.Get(c.Request.Context(), id)
	if err != nil {
		if store.IsNoRows(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		s.Logger.Error("admin get application", "err", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	files, err := s.Applications.ListFiles(c.Request.Context(), id)
	if err != nil {
		s.Logger.Error("admin list app files", "err", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	app.Files = files

	history, err := s.Applications.ListEvents(c.Request.Context(), id)
	if err != nil {
		s.Logger.Error("admin list app history", "err", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	app.History = history

	c.JSON(http.StatusOK, app)
}

type updateStatusRequest struct {
	StatusID int64  `json:"statusId" binding:"required"`
	Note     string `json:"note"`
}

func (s *Server) handleAdminUpdateApplicationStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		return
	}
	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if req.StatusID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "statusId is required"})
		return
	}
	if len(req.Note) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "note must be ≤ 500 chars"})
		return
	}

	actor := actorID(c)
	event, err := s.Applications.UpdateStatus(c.Request.Context(), id, req.StatusID, actor, strings.TrimSpace(req.Note))
	switch {
	case err == nil:
	case errors.Is(err, store.ErrSameStatus):
		// Idempotent no-op: the row is already in the requested state.
		// Returning the current application keeps the UI's "save"
		// behaviour identical to the success path.
	case store.IsNoRows(err):
		c.JSON(http.StatusNotFound, gin.H{"error": "application or status not found"})
		return
	case errors.Is(err, store.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "application or status not found"})
		return
	default:
		s.Logger.Error("admin update status", "err", err, "app_id", id, "status_id", req.StatusID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	app, err := s.Applications.Get(c.Request.Context(), id)
	if err != nil {
		s.Logger.Error("admin reload after status", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	files, err := s.Applications.ListFiles(c.Request.Context(), id)
	if err == nil {
		app.Files = files
	}
	history, err := s.Applications.ListEvents(c.Request.Context(), id)
	if err == nil {
		app.History = history
	}
	if event != nil {
		s.Logger.Info("application status changed",
			"app_id", id, "to_status_id", req.StatusID, "actor_id", actor)
	}
	c.JSON(http.StatusOK, app)
}

func (s *Server) handleAdminDownloadApplicationFile(c *gin.Context) {
	appID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || appID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		return
	}
	fileID, err := strconv.ParseInt(c.Param("fileId"), 10, 64)
	if err != nil || fileID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file id"})
		return
	}

	f, err := s.Applications.GetFile(c.Request.Context(), fileID)
	if err != nil {
		if store.IsNoRows(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		s.Logger.Error("admin download: lookup", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if f.ApplicationID != appID {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	fh, err := s.Uploads.Open(f.StoredPath)
	if err != nil {
		s.Logger.Error("admin download: open", "err", err, "path", f.StoredPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "file unavailable"})
		return
	}
	defer fh.Close()

	ct := f.ContentType
	if ct == "" {
		ct = "application/octet-stream"
	}
	c.Header("Content-Type", ct)
	c.Header("Content-Length", strconv.FormatInt(f.SizeBytes, 10))
	c.Header("Content-Disposition",
		`attachment; filename="`+sanitiseHeader(f.OriginalName)+`"`)
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Cache-Control", "private, no-store")
	c.Status(http.StatusOK)
	if _, err := io.Copy(c.Writer, fh); err != nil {
		s.Logger.Warn("admin download: stream", "err", err, "id", fileID)
	}
}

// sanitiseHeader strips characters that would break a quoted-string
// HTTP header value. Anything risky is replaced with an underscore.
func sanitiseHeader(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if r < 0x20 || r == 0x7f || r == '"' || r == '\\' {
			out = append(out, '_')
			continue
		}
		out = append(out, r)
	}
	return string(out)
}
