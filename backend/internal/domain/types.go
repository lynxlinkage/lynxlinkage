package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// StringSlice is a []string that serialises to a JSON array in the database
// column. SQLite has no native array type so we store JSON; this also keeps
// the schema portable to Postgres later (where we could swap to a TEXT[] or
// jsonb column without changing the Go API).
type StringSlice []string

// Value implements driver.Valuer.
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan implements sql.Scanner.
func (s *StringSlice) Scan(src any) error {
	if src == nil {
		*s = nil
		return nil
	}
	var raw string
	switch v := src.(type) {
	case string:
		raw = v
	case []byte:
		raw = string(v)
	default:
		return fmt.Errorf("StringSlice: cannot scan %T", src)
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		*s = nil
		return nil
	}
	if !strings.HasPrefix(raw, "[") {
		return errors.New("StringSlice: expected JSON array")
	}
	return json.Unmarshal([]byte(raw), s)
}
