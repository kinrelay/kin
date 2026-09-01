package socialcontext

import (
	"context"
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

func TestMemoryReadRepositoryProjectsOnlyRequestedOwnerValidatedContexts(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	otherOwnerID, _ := domainidentity.NewID("owner-2")

	ownerContext := mustPromotedSocialContext(t, "最近對可靠性與系統取捨特別有興趣", "activity-1", "最近開始深入研究分散式系統設計")
	otherContext := mustPromotedSocialContext(t, "最近對資料視覺化特別有興趣", "activity-2", "最近開始研究資料視覺化工具")

	source := NewMemoryRepository()
	inserted, err := source.SaveIfAbsent(context.Background(), ownerID, ownerContext)
	if err != nil || !inserted {
		t.Fatalf("SaveIfAbsent(owner) = %v, %v; want true, nil", inserted, err)
	}
	inserted, err = source.SaveIfAbsent(context.Background(), otherOwnerID, otherContext)
	if err != nil || !inserted {
		t.Fatalf("SaveIfAbsent(other owner) = %v, %v; want true, nil", inserted, err)
	}

	reader := NewMemoryReadRepository(source)
	items, err := reader.ListByOwner(context.Background(), ownerID)
	if err != nil {
		t.Fatalf("ListByOwner() error = %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	item := items[0]
	if item.ID == "" || item.OwnerID != ownerID || item.Meaning != ownerContext.Meaning() || item.PromotedAt.IsZero() {
		t.Fatalf("item = %#v, want owner projection with stable metadata", item)
	}
	if len(item.Provenance) != 1 || item.Provenance[0] != "activity-1" {
		t.Fatalf("provenance = %#v, want abstract Activity identifier", item.Provenance)
	}
}

func TestMemoryReadRepositoryReturnsExplicitEmptyCollection(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	reader := NewMemoryReadRepository(NewMemoryRepository())

	items, err := reader.ListByOwner(context.Background(), ownerID)
	if err != nil {
		t.Fatalf("ListByOwner() error = %v", err)
	}
	if items == nil || len(items) != 0 {
		t.Fatalf("items = %#v, want non-nil empty collection", items)
	}
}

func mustPromotedSocialContext(t *testing.T, meaning, activityID, activityContent string) domainsocialcontext.SocialContext {
	t.Helper()
	candidate, err := domainsocialcontext.NewContextCandidate(meaning, []string{activityID})
	if err != nil {
		t.Fatalf("NewContextCandidate() error = %v", err)
	}
	context, err := domainsocialcontext.PromoteContextCandidate(candidate, []domainsocialcontext.SourceActivity{{ID: activityID, Content: activityContent}})
	if err != nil {
		t.Fatalf("PromoteContextCandidate() error = %v", err)
	}
	return context
}
