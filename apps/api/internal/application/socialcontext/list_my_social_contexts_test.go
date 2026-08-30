package socialcontext

import (
	"context"
	"testing"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

type fakeSocialContextReader struct {
	items []SocialContextReadModel
}

func (f fakeSocialContextReader) ListByOwner(context.Context, domainidentity.ID) ([]SocialContextReadModel, error) {
	return append([]SocialContextReadModel(nil), f.items...), nil
}

func TestListMySocialContextsReturnsOnlyRequesterContextsNewestFirst(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	otherOwnerID, _ := domainidentity.NewID("owner-2")
	older := time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC)
	newer := older.Add(time.Hour)

	query := NewListMySocialContexts(fakeSocialContextReader{items: []SocialContextReadModel{
		{ID: "context-b", OwnerID: ownerID, Meaning: "較新的脈絡", Provenance: []string{"activity-2"}, PromotedAt: newer},
		{ID: "context-other", OwnerID: otherOwnerID, Meaning: "別人的脈絡", Provenance: []string{"activity-x"}, PromotedAt: newer.Add(time.Hour)},
		{ID: "context-a", OwnerID: ownerID, Meaning: "較舊的脈絡", Provenance: []string{"activity-1"}, PromotedAt: older},
	}})

	items, err := query.Execute(context.Background(), ListMySocialContextsQuery{RequesterID: "owner-1"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}
	if items[0].ID != "context-b" || items[1].ID != "context-a" {
		t.Fatalf("items order = %#v, want newest first", items)
	}
	if items[0].OwnerID != ownerID || items[1].OwnerID != ownerID {
		t.Fatalf("items = %#v, want requester-only contexts", items)
	}
	if len(items[0].Provenance) != 1 || items[0].Provenance[0] != "activity-2" {
		t.Fatalf("provenance = %#v, want abstract source identifiers", items[0].Provenance)
	}
}

func TestListMySocialContextsUsesIDAsDeterministicTieBreaker(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	promotedAt := time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC)
	query := NewListMySocialContexts(fakeSocialContextReader{items: []SocialContextReadModel{
		{ID: "context-b", OwnerID: ownerID, Meaning: "B", Provenance: []string{"activity-2"}, PromotedAt: promotedAt},
		{ID: "context-a", OwnerID: ownerID, Meaning: "A", Provenance: []string{"activity-1"}, PromotedAt: promotedAt},
	}})

	items, err := query.Execute(context.Background(), ListMySocialContextsQuery{RequesterID: "owner-1"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(items) != 2 || items[0].ID != "context-a" || items[1].ID != "context-b" {
		t.Fatalf("items = %#v, want ID tie-break ordering", items)
	}
}

func TestListMySocialContextsReturnsExplicitEmptyCollection(t *testing.T) {
	query := NewListMySocialContexts(fakeSocialContextReader{})

	items, err := query.Execute(context.Background(), ListMySocialContextsQuery{RequesterID: "owner-1"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if items == nil || len(items) != 0 {
		t.Fatalf("items = %#v, want non-nil empty collection", items)
	}
}
