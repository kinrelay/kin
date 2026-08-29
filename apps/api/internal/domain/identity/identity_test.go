package identity

import (
	"errors"
	"testing"
)

func TestNewID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		want    ID
		wantErr error
	}{
		{name: "valid", value: "user-123", want: ID("user-123")},
		{name: "trims surrounding whitespace", value: "  user-123  ", want: ID("user-123")},
		{name: "empty", value: "   ", wantErr: ErrInvalidID},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewID(tt.value)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewID() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("NewID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIdentityKeepsStableID(t *testing.T) {
	t.Parallel()

	created, err := New(ID("user-123"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if created.ID() != ID("user-123") {
		t.Fatalf("Identity.ID() = %q, want %q", created.ID(), ID("user-123"))
	}
}

func TestNewIdentityRejectsForgedInvalidID(t *testing.T) {
	t.Parallel()

	for _, raw := range []ID{"", "   "} {
		created, err := New(raw)
		if !errors.Is(err, ErrInvalidID) {
			t.Fatalf("New(%q) error = %v, want %v", raw, err, ErrInvalidID)
		}
		if created != (Identity{}) {
			t.Fatalf("New(%q) = %#v, want zero Identity", raw, created)
		}
	}
}
