package activity

import (
	"context"
	"errors"
	"testing"
	"time"

	domainactivity "github.com/kinrelay/kin/apps/api/internal/domain/activity"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func adapterActivity(t *testing.T, idValue string) domainactivity.Activity {
	t.Helper()
	return adapterActivityWithContent(t, idValue, "reading about CRDTs")
}

func adapterActivityWithContent(t *testing.T, idValue, contentValue string) domainactivity.Activity {
	t.Helper()
	owner, err := domainidentity.NewID("alice")
	if err != nil {
		t.Fatalf("NewID() error = %v", err)
	}
	content, err := domainactivity.NewContent(contentValue)
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

func TestMemoryRepositoryTreatsIdenticalSaveAsIdempotentRetry(t *testing.T) {
	repository := NewMemoryRepository()
	created := adapterActivity(t, "activity-1")
	ctx := context.Background()

	if err := repository.Save(ctx, created); err != nil {
		t.Fatalf("first Save() error = %v", err)
	}
	if err := repository.Save(ctx, created); err != nil {
		t.Fatalf("identical retry Save() error = %v", err)
	}
}

func TestMemoryRepositoryRejectsConflictingActivityIDWithoutOverwrite(t *testing.T) {
	repository := NewMemoryRepository()
	original := adapterActivityWithContent(t, "activity-1", "reading about CRDTs")
	conflict := adapterActivityWithContent(t, "activity-1", "learning event sourcing")
	ctx := context.Background()

	if err := repository.Save(ctx, original); err != nil {
		t.Fatalf("first Save() error = %v", err)
	}
	if err := repository.Save(ctx, conflict); !errors.Is(err, ErrActivityIDConflict) {
		t.Fatalf("conflicting Save() error = %v, want %v", err, ErrActivityIDConflict)
	}

	repository.mu.RLock()
	stored := repository.activities[original.ID()]
	repository.mu.RUnlock()
	if stored.Content() != original.Content() {
		t.Fatalf("stored content = %q, want original %q", stored.Content().String(), original.Content().String())
	}
}
