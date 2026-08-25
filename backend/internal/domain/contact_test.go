package domain

import (
	"strings"
	"testing"
)

func TestValidateContactMFTIKFields(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		kind    ContactKind
		slug    string
		message string
		wantErr string
	}{
		{
			name:    "general without slug unchanged",
			kind:    ContactGeneral,
			message: "Just saying hello about partnership ideas.",
		},
		{
			name:    "general mentioning reserved host still ok",
			kind:    ContactGeneral,
			message: "Curious about www.mftik.app historically.",
		},
		{
			name:    "mftik valid slug",
			kind:    ContactMFTIK,
			slug:    "desk-01",
			message: "Please set up desk-01.mftik.app",
		},
		{
			name:    "mftik omitted slug ok",
			kind:    ContactMFTIK,
			message: "I'd like an instance for our team.",
		},
		{
			name:    "mftik reserved slug field",
			kind:    ContactMFTIK,
			slug:    "www",
			message: "Requesting access",
			wantErr: "reserved",
		},
		{
			name:    "mftik too short slug",
			kind:    ContactMFTIK,
			slug:    "a",
			message: "Requesting access",
			wantErr: "between 3 and 32",
		},
		{
			name:    "slug present on general still validated",
			kind:    ContactGeneral,
			slug:    "www",
			message: "hello there friend",
			wantErr: "reserved",
		},
		{
			name:    "mftik reserved host in message only",
			kind:    ContactMFTIK,
			message: "Please give us www.mftik.app",
			wantErr: "reserved",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateContactMFTIKFields(tc.kind, tc.slug, tc.message)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("got %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("got %v, want error containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestValidateMFTIKSlug(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		slug    string
		wantErr string
	}{
		{name: "omitted empty", slug: ""},
		{name: "omitted whitespace", slug: "   "},
		{name: "valid simple", slug: "desk"},
		{name: "valid with hyphen and digits", slug: "desk-01"},
		{name: "valid max length", slug: strings.Repeat("a", 32)},
		{name: "too short", slug: "a", wantErr: "between 3 and 32"},
		{name: "too short two", slug: "ab", wantErr: "between 3 and 32"},
		{name: "too long", slug: strings.Repeat("a", 33), wantErr: "between 3 and 32"},
		{name: "uppercase", slug: "Desk", wantErr: "lowercase"},
		{name: "leading hyphen", slug: "-desk", wantErr: "hyphen"},
		{name: "trailing hyphen", slug: "desk-", wantErr: "hyphen"},
		{name: "underscore", slug: "desk_01", wantErr: "lowercase"},
		{name: "reserved www", slug: "www", wantErr: "reserved"},
		{name: "reserved admin", slug: "admin", wantErr: "reserved"},
		{name: "reserved api", slug: "api", wantErr: "reserved"},
		{name: "reserved mftik", slug: "mftik", wantErr: "reserved"},
		{name: "reserved support", slug: "support", wantErr: "reserved"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateMFTIKSlug(tc.slug)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("ValidateMFTIKSlug(%q) = %v, want nil", tc.slug, err)
				}
				return
			}
			if err == nil {
				t.Fatalf("ValidateMFTIKSlug(%q) = nil, want error containing %q", tc.slug, tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("ValidateMFTIKSlug(%q) = %q, want substring %q", tc.slug, err.Error(), tc.wantErr)
			}
		})
	}
}

func TestReservedMFTIKHostInMessage(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "reserved www host",
			message: "Please provision www.mftik.app for our team.",
			want:    "www",
		},
		{
			name:    "reserved admin host case insensitive",
			message: "Need Admin.MFTIK.APP access",
			want:    "admin",
		},
		{
			name:    "valid host ignored",
			message: "Looking for desk-01.mftik.app please.",
			want:    "",
		},
		{
			name:    "no host mention",
			message: "Hello, I would like access to MFTIK for my company.",
			want:    "",
		},
		{
			name:    "substring not a host",
			message: "email me about notwww.mftik.appish things",
			want:    "",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ReservedMFTIKHostInMessage(tc.message)
			if got != tc.want {
				t.Fatalf("ReservedMFTIKHostInMessage(...) = %q, want %q", got, tc.want)
			}
		})
	}
}
