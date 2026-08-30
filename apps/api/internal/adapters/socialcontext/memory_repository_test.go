package socialcontext

import (
	"context"
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

func TestMemoryRepositorySavesPromotedSocialContext(t *testing.T) {
	candidate, err := domainsocialcontext.NewContextCandidate("最近對分散式系統的可靠性與取捨特別有興趣", []string{"activity-1"})
	if err != nil {
		t.Fatalf("NewContextCandidate() error = %v", err)
	}
	socialContext, err := domainsocialcontext.PromoteContextCandidate(candidate, []domainsocialcontext.SourceActivity{{ID: "activity-1", Content: "最近開始深入研究分散式系統設計"}})
	if err != nil {
		t.Fatalf("PromoteContextCandidate() error = %v", err)
	}
	ownerID, _ := domainidentity.NewID("owner-1")

	repository := NewMemoryRepository()
	if err := repository.Save(context.Background(), ownerID, socialContext); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	stored := repository.All()
	if len(stored) != 1 || stored[0].Meaning() != socialContext.Meaning() || !stored[0].IsPrivate() {
		t.Fatalf("All() = %#v, want one private promoted SocialContext", stored)
	}
}
