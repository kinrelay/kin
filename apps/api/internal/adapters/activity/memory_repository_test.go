package activity

import (
	"context"
	"testing"
	"time"

	domainactivity "github.com/kinrelay/kin/apps/api/internal/domain/activity"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func adapterActivity(t *testing.T, idValue string) domainactivity.Activity {
	t.Helper()
	owner, err := domainidentity.NewID("alice")
	if err != nil {
		t.Fatalf("NewID() error = %v", err)
	}
	content, err := domainactivity.NewContent("reading about CRDTs")
	if err != nil {
		t.Fatalf("NewContent() error = %v", err)
	}
	now := time.Date(2026, time.August, 29, 10, 0, 0, 0, time.UTC)
	created, err := domainactivity.NewManual(idValue, owner, content, now, now)
	if err != nil {
		t.Fatalf("NewManual() error = %v", err)
	}
	return created
}

func TestMemoryRepositorySavesNormalizedActivityByStableID(t *testing.T) {
	repository := NewMemoryRepository()
	created := adapterActivity(t, "activity-1")

	if err := repository.Save(context.Background(), created); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	repository.mu.RLock()
	stored, ok := repository.activities[created.ID()]
	repository.mu.RUnlock()
	if !ok {
		t.Fatalf("activity %q was not stored", created.ID())
	}
	if stored.ID() != created.ID() || stored.OwnerID() != created.OwnerID() || stored.Content() != created.Content() {
		t.Fatalf("stored activity = %#v, want %#v", stored, created)
	}
	if !stored.IsPrivate() || stored.Provenance() != domainactivity.ProvenanceManual {
		t.Fatal("repository must preserve private/manual domain state")
	}
}
