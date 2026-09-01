package socialcontext

import (
	"context"
	"sync"
	"testing"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

func TestMemoryRepositoryReconcileCurrentStateRejectsDelayedStaleAffirmative(t *testing.T) {
	ctx := context.Background()
	ownerID, _ := domainidentity.NewID("owner-1")
	repository := NewMemoryRepository()

	affirmed := promotedContext(t,
		"近期持續投入馬拉松訓練與耐力運動",
		"activity-a",
		"最近開始準備第一次全程馬拉松訓練",
	)
	reaffirmed := promotedContext(t,
		"近期重新投入馬拉松訓練與耐力運動",
		"activity-c",
		"最近又開始規律準備馬拉松訓練",
	)

	t1 := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	t2 := t1.Add(time.Hour)
	t3 := t2.Add(time.Hour)

	changed, err := repository.ReconcileCurrentState(ctx, ownerID, "marathon-training", t1, &affirmed)
	if err != nil || !changed {
		t.Fatalf("ReconcileCurrentState(affirm) = %v, %v; want true, nil", changed, err)
	}
	changed, err = repository.ReconcileCurrentState(ctx, ownerID, "marathon-training", t3, nil)
	if err != nil || !changed {
		t.Fatalf("ReconcileCurrentState(reverse) = %v, %v; want true, nil", changed, err)
	}

	changed, err = repository.ReconcileCurrentState(ctx, ownerID, "marathon-training", t2, &reaffirmed)
	if err != nil {
		t.Fatalf("ReconcileCurrentState(stale affirm) error = %v", err)
	}
	if changed {
		t.Fatal("ReconcileCurrentState(stale affirm) changed = true, want deterministic no-op")
	}
	if got := repository.All(); len(got) != 0 {
		t.Fatalf("All() = %#v, want no stale context after newer reversal", got)
	}
}

func TestMemoryRepositoryReconcileCurrentStateIsAtomicAcrossInterleavedReversalAndStaleRetry(t *testing.T) {
	ctx := context.Background()
	ownerID, _ := domainidentity.NewID("owner-1")
	repository := NewMemoryRepository()

	initial := promotedContext(t,
		"近期持續投入馬拉松訓練與耐力運動",
		"activity-a",
		"最近開始準備第一次全程馬拉松訓練",
	)
	staleRetry := promotedContext(t,
		"近期仍持續投入馬拉松訓練與耐力運動",
		"activity-b",
		"這週繼續安排馬拉松長距離訓練",
	)

	t1 := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	t2 := t1.Add(time.Hour)
	t3 := t2.Add(time.Hour)

	if changed, err := repository.ReconcileCurrentState(ctx, ownerID, "marathon-training", t1, &initial); err != nil || !changed {
		t.Fatalf("ReconcileCurrentState(initial) = %v, %v; want true, nil", changed, err)
	}

	start := make(chan struct{})
	var wg sync.WaitGroup
	for _, mutation := range []struct {
		occurredAt  time.Time
		replacement *domainsocialcontext.SocialContext
	}{
		{occurredAt: t3, replacement: nil},
		{occurredAt: t2, replacement: &staleRetry},
	} {
		mutation := mutation
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			if _, err := repository.ReconcileCurrentState(ctx, ownerID, "marathon-training", mutation.occurredAt, mutation.replacement); err != nil {
				t.Errorf("ReconcileCurrentState() error = %v", err)
			}
		}()
	}
	close(start)
	wg.Wait()

	if got := repository.All(); len(got) != 0 {
		t.Fatalf("All() = %#v, want newer reversal to win regardless of command interleaving", got)
	}
}

func promotedContext(t *testing.T, meaning, activityID, content string) domainsocialcontext.SocialContext {
	t.Helper()
	candidate, err := domainsocialcontext.NewContextCandidate(meaning, []string{activityID})
	if err != nil {
		t.Fatalf("NewContextCandidate() error = %v", err)
	}
	result, err := domainsocialcontext.PromoteContextCandidate(candidate, []domainsocialcontext.SourceActivity{{ID: activityID, Content: content}})
	if err != nil {
		t.Fatalf("PromoteContextCandidate() error = %v", err)
	}
	return result
}
