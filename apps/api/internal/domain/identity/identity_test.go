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

	id, err := NewID("user-123")
	if err != nil {
		t.Fatalf("NewID() error = %v", err)
	}

	created := New(id)
	if created.ID() != id {
		t.Fatalf("Identity.ID() = %q, want %q", created.ID(), id)
	}
}
