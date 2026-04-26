// Command seed loads YAML files from backend/seed/ into the database. It is
// idempotent: each entity is upserted by primary key (or unique name for
// partners) so re-running the command updates rather than duplicates.
//
// Usage:
//
//	go run ./cmd/seed -dir ./seed
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/lynxlinkage/lynxlinkage/backend/internal/config"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/domain"
	"github.com/lynxlinkage/lynxlinkage/backend/internal/store"
	"gopkg.in/yaml.v3"
)

type seedFile struct {
	Researches []researchSeed `yaml:"researches"`
	Jobs       []jobSeed      `yaml:"jobs"`
	Partners   []partnerSeed  `yaml:"partners"`
}

type researchSeed struct {
	ID            int64    `yaml:"id"`
	Title         string   `yaml:"title"`
	Summary       string   `yaml:"summary"`
	Tags          []string `yaml:"tags"`
	CoverImageURL string   `yaml:"cover_image_url"`
	ExternalURL   string   `yaml:"external_url"`
	Source        string   `yaml:"source"`
	PublishedAt   string   `yaml:"published_at"`
	DisplayOrder  int      `yaml:"display_order"`
}

type jobSeed struct {
	ID              int64  `yaml:"id"`
	Title           string `yaml:"title"`
	Team            string `yaml:"team"`
	Location        string `yaml:"location"`
	EmploymentType  string `yaml:"employment_type"`
	DescriptionMD   string `yaml:"description_md"`
	ApplyURLOrEmail string `yaml:"apply_url_or_email"`
	PostedAt        string `yaml:"posted_at"`
	IsActive        *bool  `yaml:"is_active"`
}

type partnerSeed struct {
	Name         string `yaml:"name"`
	LogoURL      string `yaml:"logo_url"`
	WebsiteURL   string `yaml:"website_url"`
	Tier         string `yaml:"tier"`
	Description  string `yaml:"description"`
	DisplayOrder int    `yaml:"display_order"`
}

func main() {
	dir := flag.String("dir", "./seed", "directory containing seed YAML files")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	if err := run(*dir, logger); err != nil {
		logger.Error("seed failed", "err", err)
		os.Exit(1)
	}
}

func run(dir string, logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx := context.Background()
	db, err := store.Open(ctx, cfg.DatabaseURL, logger)
	if err != nil {
		return err
	}
	defer db.Close()

	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return err
	}
	more, _ := filepath.Glob(filepath.Join(dir, "*.yml"))
	files = append(files, more...)
	sort.Strings(files)

	if len(files) == 0 {
		logger.Warn("no seed files found", "dir", dir)
		return nil
	}

	rRepo := store.NewResearchRepo(db)
	jRepo := store.NewJobRepo(db)
	pRepo := store.NewPartnerRepo(db)

	totals := struct{ R, J, P int }{}
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}
		var s seedFile
		if err := yaml.Unmarshal(raw, &s); err != nil {
			return fmt.Errorf("yaml %s: %w", f, err)
		}

		for _, r := range s.Researches {
			t, err := parseDate(r.PublishedAt)
			if err != nil {
				return fmt.Errorf("%s: research %q published_at: %w", f, r.Title, err)
			}
			source := r.Source
			if source == "" {
				source = string(domain.SourceInternal)
			}
			card := &domain.ResearchCard{
				ID:            r.ID,
				Title:         r.Title,
				Summary:       r.Summary,
				Tags:          domain.StringSlice(r.Tags),
				CoverImageURL: r.CoverImageURL,
				ExternalURL:   r.ExternalURL,
				Source:        domain.ResearchSource(source),
				PublishedAt:   t,
				DisplayOrder:  r.DisplayOrder,
			}
			if err := rRepo.Upsert(ctx, card); err != nil {
				return fmt.Errorf("upsert research %q: %w", r.Title, err)
			}
			totals.R++
		}

		for _, j := range s.Jobs {
			t, err := parseDate(j.PostedAt)
			if err != nil {
				return fmt.Errorf("%s: job %q posted_at: %w", f, j.Title, err)
			}
			active := true
			if j.IsActive != nil {
				active = *j.IsActive
			}
			et := j.EmploymentType
			if et == "" {
				et = string(domain.EmploymentFullTime)
			}
			job := &domain.JobPosting{
				ID:              j.ID,
				Title:           j.Title,
				Team:            j.Team,
				Location:        j.Location,
				EmploymentType:  domain.EmploymentType(et),
				DescriptionMD:   j.DescriptionMD,
				ApplyURLOrEmail: j.ApplyURLOrEmail,
				PostedAt:        t,
				IsActive:        active,
			}
			if err := jRepo.Upsert(ctx, job); err != nil {
				return fmt.Errorf("upsert job %q: %w", j.Title, err)
			}
			totals.J++
		}

		for _, p := range s.Partners {
			tier := p.Tier
			if tier == "" {
				tier = string(domain.TierStrategic)
			}
			pr := &domain.Partner{
				Name:         p.Name,
				LogoURL:      p.LogoURL,
				WebsiteURL:   p.WebsiteURL,
				Tier:         domain.PartnerTier(tier),
				Description:  p.Description,
				DisplayOrder: p.DisplayOrder,
			}
			if err := pRepo.Upsert(ctx, pr); err != nil {
				return fmt.Errorf("upsert partner %q: %w", p.Name, err)
			}
			totals.P++
		}

		logger.Info("loaded", "file", filepath.Base(f),
			"researches", len(s.Researches),
			"jobs", len(s.Jobs),
			"partners", len(s.Partners))
	}

	logger.Info("seed complete", "researches", totals.R, "jobs", totals.J, "partners", totals.P)
	return nil
}

// parseDate accepts RFC3339 ("2026-04-01T09:00:00Z") or just a date
// ("2026-04-01"); the latter is anchored to UTC midnight.
func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty date")
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("unrecognised date %q", s)
}
