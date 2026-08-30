package activity

import (
	"errors"
	"testing"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func activityIdentityID(t *testing.T, value string) domainidentity.ID {
	t.Helper()
	id, err := domainidentity.NewID(value)
	if err != nil {
		t.Fatalf("NewID(%q): %v", value, err)
	}
	return id
}

func TestNewContentNormalizesAndRejectsBlankMeaning(t *testing.T) {
	content, err := NewContent("  learning Go domain modeling  ")
	if err != nil {
		t.Fatalf("NewContent() error = %v", err)
	}
	if got := content.String(); got != "learning Go domain modeling" {
		t.Fatalf("Content.String() = %q, want normalized content", got)
	}

	for _, input := range []string{"", "   ", "\t\n"} {
		_, err := NewContent(input)
		if !errors.Is(err, ErrEmptyContent) {
			t.Fatalf("NewContent(%q) error = %v, want %v", input, err, ErrEmptyContent)
		}
	}
}

func TestNewManualActivityRecordsOwnerProvenanceAndPrivateDefault(t *testing.T) {
	owner := activityIdentityID(t, "alice")
	content, err := NewContent("researching local-first apps")
	if err != nil {
		t.Fatalf("NewContent() error = %v", err)
	}
	occurredAt := time.Date(2026, time.August, 29, 9, 30, 0, 0, time.UTC)
	contributedAt := time.Date(2026, time.August, 29, 10, 0, 0, 0, time.UTC)

	created, err := NewManual("activity-1", owner, content, occurredAt, contributedAt)
	if err != nil {
		t.Fatalf("NewManual() error = %v", err)
	}
	if created.ID() != ID("activity-1") {
		t.Fatalf("ID() = %q, want activity-1", created.ID())
	}
	if created.OwnerID() != owner {
		t.Fatalf("OwnerID() = %q, want %q", created.OwnerID(), owner)
	}
	if created.Content() != content {
		t.Fatalf("Content() = %q, want %q", created.Content().String(), content.String())
	}
	if created.Provenance() != ProvenanceManual {
		t.Fatalf("Provenance() = %q, want %q", created.Provenance(), ProvenanceManual)
	}
	if !created.IsPrivate() {
		t.Fatal("new manual activity must be private by default")
	}
	if !created.OccurredAt().Equal(occurredAt) {
		t.Fatalf("OccurredAt() = %v, want %v", created.OccurredAt(), occurredAt)
	}
	if !created.ContributedAt().Equal(contributedAt) {
		t.Fatalf("ContributedAt() = %v, want %v", created.ContributedAt(), contributedAt)
	}
}

func TestNewManualActivityRejectsInvalidIdentityAndMetadata(t *testing.T) {
	content, err := NewContent("reading about CRDTs")
	if err != nil {
		t.Fatalf("NewContent() error = %v", err)
	}
	now := time.Date(2026, time.August, 29, 10, 0, 0, 0, time.UTC)
	validOwner := activityIdentityID(t, "alice")

	tests := []struct {
		name        string
		id          string
		owner       domainidentity.ID
		occurredAt  time.Time
		contributed time.Time
		wantErr     error
	}{
		{name: "blank activity id", id: "  ", owner: validOwner, occurredAt: now, contributed: now, wantErr: ErrInvalidID},
		{name: "invalid owner", id: "activity-1", owner: domainidentity.ID("   "), occurredAt: now, contributed: now, wantErr: domainidentity.ErrInvalidID},
		{name: "missing occurred time", id: "activity-1", owner: validOwner, contributed: now, wantErr: ErrInvalidTimestamp},
		{name: "missing contributed time", id: "activity-1", owner: validOwner, occurredAt: now, wantErr: ErrInvalidTimestamp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewManual(tt.id, tt.owner, content, tt.occurredAt, tt.contributed)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewManual() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
