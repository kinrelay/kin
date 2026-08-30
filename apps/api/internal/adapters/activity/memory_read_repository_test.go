package activity

import (
	"context"
	"testing"
	"time"

	applicationactivity "github.com/kinrelay/kin/apps/api/internal/application/activity"
	domainactivity "github.com/kinrelay/kin/apps/api/internal/domain/activity"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func readAdapterActivity(t *testing.T, idValue, ownerValue, contentValue string, occurredAt, contributedAt time.Time) domainactivity.Activity {
	t.Helper()
	owner, err := domainidentity.NewID(ownerValue)
	if err != nil {
		t.Fatalf("NewID(%q): %v", ownerValue, err)
	}
	content, err := domainactivity.NewContent(contentValue)
	if err != nil {
		t.Fatalf("NewContent(%q): %v", contentValue, err)
	}
	created, err := domainactivity.NewManual(idValue, owner, content, occurredAt, contributedAt)
	if err != nil {
		t.Fatalf("NewManual(%q): %v", idValue, err)
	}
	return created
}

func TestMemoryReadRepositoryProjectsOnlyRequestedOwnerWithoutExposingAggregate(t *testing.T) {
	ctx := context.Background()
	writeRepository := NewMemoryRepository()
	readRepository := NewMemoryReadRepository(writeRepository)
	aliceID, err := domainidentity.NewID("alice")
	if err != nil {
		t.Fatalf("NewID(alice): %v", err)
	}
	baseTime := time.Date(2026, time.August, 30, 1, 0, 0, 0, time.UTC)
	aliceActivity := readAdapterActivity(t, "alice-1", "alice", "reading CRDT papers", baseTime, baseTime.Add(time.Minute))
	bobActivity := readAdapterActivity(t, "bob-1", "bob", "private bob activity", baseTime, baseTime)
	for _, value := range []domainactivity.Activity{aliceActivity, bobActivity} {
		if err := writeRepository.Save(ctx, value); err != nil {
			t.Fatalf("Save(%q) error = %v", value.ID(), err)
		}
	}

	got, err := readRepository.ListByOwner(ctx, aliceID)
	if err != nil {
		t.Fatalf("ListByOwner() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("result count = %d, want 1: %#v", len(got), got)
	}
	want := applicationactivity.ActivityReadModel{
		ID:            "alice-1",
		OwnerID:       aliceID,
		Content:       "reading CRDT papers",
		Provenance:    "manual",
		OccurredAt:    aliceActivity.OccurredAt(),
		ContributedAt: aliceActivity.ContributedAt(),
	}
	if got[0] != want {
		t.Fatalf("projection = %#v, want %#v", got[0], want)
	}
}

func TestMemoryReadRepositoryReturnsExplicitEmptyCollectionForUnknownOwner(t *testing.T) {
	readRepository := NewMemoryReadRepository(NewMemoryRepository())
	owner, err := domainidentity.NewID("alice")
	if err != nil {
		t.Fatalf("NewID(): %v", err)
	}

	got, err := readRepository.ListByOwner(context.Background(), owner)
	if err != nil {
		t.Fatalf("ListByOwner() error = %v", err)
	}
	if got == nil {
		t.Fatal("ListByOwner() result = nil, want explicit empty collection")
	}
	if len(got) != 0 {
		t.Fatalf("result count = %d, want 0", len(got))
	}
}
