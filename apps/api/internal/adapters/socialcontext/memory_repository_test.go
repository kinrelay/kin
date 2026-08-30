package socialcontext

import (
	"context"
	"sync"
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

func TestMemoryRepositoryEquivalentContextIsOwnerScopedAndMeaningSpecific(t *testing.T) {
	ctx := context.Background()
	ownerID, _ := domainidentity.NewID("owner-1")
	otherOwnerID, _ := domainidentity.NewID("owner-2")

	candidate, _ := domainsocialcontext.NewContextCandidate("近期關注分散式系統的一致性模型、可靠性與工程取捨", []string{"activity-1"})
	storedContext, err := domainsocialcontext.PromoteContextCandidate(candidate, []domainsocialcontext.SourceActivity{{ID: "activity-1", Content: "最近開始深入研究分散式系統設計"}})
	if err != nil {
		t.Fatalf("PromoteContextCandidate(stored) error = %v", err)
	}
	equivalentCandidate, _ := domainsocialcontext.NewContextCandidate("近期關注分散式系統的一致性模型、可靠性與工程取捨", []string{"activity-2"})
	equivalentContext, err := domainsocialcontext.PromoteContextCandidate(equivalentCandidate, []domainsocialcontext.SourceActivity{{ID: "activity-2", Content: "持續比較不同一致性模型的工程取捨"}})
	if err != nil {
		t.Fatalf("PromoteContextCandidate(equivalent) error = %v", err)
	}
	differentCandidate, _ := domainsocialcontext.NewContextCandidate("近期關注耐力運動與長距離訓練", []string{"activity-3"})
	differentContext, err := domainsocialcontext.PromoteContextCandidate(differentCandidate, []domainsocialcontext.SourceActivity{{ID: "activity-3", Content: "最近開始準備第一次全程馬拉松訓練"}})
	if err != nil {
		t.Fatalf("PromoteContextCandidate(different) error = %v", err)
	}

	repository := NewMemoryRepository()
	if err := repository.Save(ctx, ownerID, storedContext); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := repository.ExistsEquivalent(ctx, ownerID, equivalentContext)
	if err != nil || !got {
		t.Fatalf("ExistsEquivalent(same owner, same meaning) = %v, %v; want true, nil", got, err)
	}
	got, err = repository.ExistsEquivalent(ctx, otherOwnerID, equivalentContext)
	if err != nil || got {
		t.Fatalf("ExistsEquivalent(other owner) = %v, %v; want false, nil", got, err)
	}
	got, err = repository.ExistsEquivalent(ctx, ownerID, differentContext)
	if err != nil || got {
		t.Fatalf("ExistsEquivalent(different meaning) = %v, %v; want false, nil", got, err)
	}
}

func TestMemoryRepositorySaveIfAbsentIsAtomicAcrossConcurrentCalls(t *testing.T) {
	ctx := context.Background()
	ownerID, _ := domainidentity.NewID("owner-1")
	candidate, _ := domainsocialcontext.NewContextCandidate("近期關注分散式系統的一致性模型、可靠性與工程取捨", []string{"activity-1"})
	socialContext, err := domainsocialcontext.PromoteContextCandidate(candidate, []domainsocialcontext.SourceActivity{{ID: "activity-1", Content: "最近開始深入研究分散式系統設計"}})
	if err != nil {
		t.Fatalf("PromoteContextCandidate() error = %v", err)
	}

	repository := NewMemoryRepository()
	const workers = 16
	inserted := make(chan bool, workers)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			ok, err := repository.SaveIfAbsent(ctx, ownerID, socialContext)
			if err != nil {
				t.Errorf("SaveIfAbsent() error = %v", err)
				return
			}
			inserted <- ok
		}()
	}
	close(start)
	wg.Wait()
	close(inserted)

	insertCount := 0
	for ok := range inserted {
		if ok {
			insertCount++
		}
	}
	if insertCount != 1 {
		t.Fatalf("insert count = %d, want exactly 1", insertCount)
	}
	if got := len(repository.All()); got != 1 {
		t.Fatalf("stored context count = %d, want exactly 1", got)
	}
}
