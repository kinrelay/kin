package activity

import (
	"context"
	"testing"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

type activityReadPortFake struct {
	items          []ActivityReadModel
	requestedOwner domainidentity.ID
	err            error
}

func (f *activityReadPortFake) ListByOwner(_ context.Context, ownerID domainidentity.ID) ([]ActivityReadModel, error) {
	f.requestedOwner = ownerID
	if f.err != nil {
		return nil, f.err
	}
	return append([]ActivityReadModel(nil), f.items...), nil
}

func TestListMyActivitiesUsesRequesterAsOwnerAndProjectsOnlyOwnActivities(t *testing.T) {
	aliceID, err := domainidentity.NewID("alice")
	if err != nil {
		t.Fatalf("NewID(alice): %v", err)
	}
	bobID, err := domainidentity.NewID("bob")
	if err != nil {
		t.Fatalf("NewID(bob): %v", err)
	}
	baseTime := time.Date(2026, time.August, 30, 1, 0, 0, 0, time.UTC)
	reader := &activityReadPortFake{items: []ActivityReadModel{
		{ID: "activity-bob", OwnerID: bobID, Content: "private bob activity", Provenance: "manual", OccurredAt: baseTime, ContributedAt: baseTime},
		{ID: "activity-alice", OwnerID: aliceID, Content: "reading CRDT papers", Provenance: "manual", OccurredAt: baseTime, ContributedAt: baseTime},
	}}
	query := NewListMyActivities(reader)

	got, err := query.Execute(context.Background(), ListMyActivitiesQuery{RequesterID: " alice "})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if reader.requestedOwner != aliceID {
		t.Fatalf("read port owner = %q, want %q", reader.requestedOwner, aliceID)
	}
	if len(got) != 1 {
		t.Fatalf("result count = %d, want 1 owner-only item: %#v", len(got), got)
	}
	if got[0].ID != "activity-alice" || got[0].OwnerID != aliceID || got[0].Content != "reading CRDT papers" || got[0].Provenance != "manual" {
		t.Fatalf("result = %#v, want purpose-built alice projection", got[0])
	}
}

func TestListMyActivitiesHasDeterministicNewestFirstOrdering(t *testing.T) {
	owner, err := domainidentity.NewID("alice")
	if err != nil {
		t.Fatalf("NewID(): %v", err)
	}
	old := time.Date(2026, time.August, 29, 10, 0, 0, 0, time.UTC)
	newer := time.Date(2026, time.August, 30, 9, 0, 0, 0, time.UTC)
	reader := &activityReadPortFake{items: []ActivityReadModel{
		{ID: "z-tie", OwnerID: owner, Content: "tie z", Provenance: "manual", OccurredAt: newer, ContributedAt: newer},
		{ID: "old", OwnerID: owner, Content: "old", Provenance: "manual", OccurredAt: old, ContributedAt: old},
		{ID: "a-tie", OwnerID: owner, Content: "tie a", Provenance: "manual", OccurredAt: newer, ContributedAt: newer},
	}}
	query := NewListMyActivities(reader)

	got, err := query.Execute(context.Background(), ListMyActivitiesQuery{RequesterID: "alice"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	want := []string{"a-tie", "z-tie", "old"}
	if len(got) != len(want) {
		t.Fatalf("result count = %d, want %d", len(got), len(want))
	}
	for i, id := range want {
		if got[i].ID != id {
			t.Fatalf("result[%d].ID = %q, want %q; full result: %#v", i, got[i].ID, id, got)
		}
	}
}

func TestListMyActivitiesReturnsEmptyCollectionForNoActivities(t *testing.T) {
	query := NewListMyActivities(&activityReadPortFake{})

	got, err := query.Execute(context.Background(), ListMyActivitiesQuery{RequesterID: "alice"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got == nil {
		t.Fatal("Execute() result = nil, want explicit empty collection")
	}
	if len(got) != 0 {
		t.Fatalf("result count = %d, want 0", len(got))
	}
}

func TestListMyActivitiesRejectsInvalidRequesterBeforeRead(t *testing.T) {
	reader := &activityReadPortFake{}
	query := NewListMyActivities(reader)

	_, err := query.Execute(context.Background(), ListMyActivitiesQuery{RequesterID: "   "})
	if err == nil {
		t.Fatal("Execute() error = nil, want invalid Identity ID")
	}
	if reader.requestedOwner != "" {
		t.Fatalf("read port was called with %q for invalid requester", reader.requestedOwner)
	}
}
