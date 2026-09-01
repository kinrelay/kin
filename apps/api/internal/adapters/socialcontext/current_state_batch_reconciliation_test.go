package socialcontext

import (
	"context"
	"testing"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

func TestMemoryRepositoryReconcileOwnerCurrentStateIsAllOrNothing(t *testing.T) {
	ctx := context.Background()
	ownerID, _ := domainidentity.NewID("owner-1")
	repository := NewMemoryRepository()

	marathon := promotedContext(t,
		"近期持續投入馬拉松訓練與耐力運動",
		"activity-marathon",
		"最近開始準備第一次全程馬拉松訓練",
	)
	distributed := promotedContext(t,
		"近期關注分散式系統的一致性模型與可靠性取捨",
		"activity-distributed",
		"最近開始深入研究分散式系統設計與一致性模型",
	)
	occurredAt := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)

	_, err := repository.ReconcileOwnerCurrentState(ctx, ownerID, []domainsocialcontext.CurrentStateMutation{
		{SemanticID: "marathon", OccurredAt: occurredAt, Replacement: &marathon},
		{SemanticID: "", OccurredAt: occurredAt, Replacement: &distributed},
	})
	if err == nil {
		t.Fatal("ReconcileOwnerCurrentState() error = nil, want invalid semantic identity")
	}
	if got := repository.All(); len(got) != 0 {
		t.Fatalf("All() = %#v, want no partial state after rejected batch", got)
	}
}

func TestMemoryRepositoryReconcileOwnerCurrentStateRetiresOnlyTargetSemanticComponent(t *testing.T) {
	ctx := context.Background()
	ownerID, _ := domainidentity.NewID("owner-1")
	repository := NewMemoryRepository()

	marathon := promotedContext(t,
		"近期持續投入馬拉松訓練與耐力運動",
		"activity-marathon",
		"最近開始準備第一次全程馬拉松訓練",
	)
	distributed := promotedContext(t,
		"近期關注分散式系統的一致性模型與可靠性取捨",
		"activity-distributed",
		"最近開始深入研究分散式系統設計與一致性模型",
	)
	t1 := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	t2 := t1.Add(time.Hour)

	changed, err := repository.ReconcileOwnerCurrentState(ctx, ownerID, []domainsocialcontext.CurrentStateMutation{
		{SemanticID: "marathon", OccurredAt: t1, Replacement: &marathon},
		{SemanticID: "distributed-systems", OccurredAt: t1, Replacement: &distributed},
	})
	if err != nil || changed != 2 {
		t.Fatalf("ReconcileOwnerCurrentState(initial) = %d, %v; want 2, nil", changed, err)
	}

	changed, err = repository.ReconcileOwnerCurrentState(ctx, ownerID, []domainsocialcontext.CurrentStateMutation{
		{SemanticID: "marathon", OccurredAt: t2, Replacement: nil},
	})
	if err != nil || changed != 1 {
		t.Fatalf("ReconcileOwnerCurrentState(reverse marathon) = %d, %v; want 1, nil", changed, err)
	}

	stored := repository.All()
	if len(stored) != 1 {
		t.Fatalf("All() count = %d, want unrelated semantic component preserved", len(stored))
	}
	if stored[0].Meaning() != distributed.Meaning() {
		t.Fatalf("All()[0].Meaning() = %q, want %q", stored[0].Meaning(), distributed.Meaning())
	}
}

func TestMemoryRepositoryReconcileOwnerCurrentStateRefreshesProvenanceAcrossReaffirmation(t *testing.T) {
	ctx := context.Background()
	ownerID, _ := domainidentity.NewID("owner-1")
	repository := NewMemoryRepository()

	first := promotedContext(t,
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

	for _, mutation := range []domainsocialcontext.CurrentStateMutation{
		{SemanticID: "marathon", OccurredAt: t1, Replacement: &first},
		{SemanticID: "marathon", OccurredAt: t2, Replacement: nil},
		{SemanticID: "marathon", OccurredAt: t3, Replacement: &reaffirmed},
	} {
		if _, err := repository.ReconcileOwnerCurrentState(ctx, ownerID, []domainsocialcontext.CurrentStateMutation{mutation}); err != nil {
			t.Fatalf("ReconcileOwnerCurrentState() error = %v", err)
		}
	}

	stored := repository.All()
	if len(stored) != 1 {
		t.Fatalf("All() count = %d, want one re-established context", len(stored))
	}
	provenance := stored[0].Provenance()
	if len(provenance) != 1 || provenance[0] != "activity-c" {
		t.Fatalf("Provenance() = %#v, want refreshed provenance [activity-c]", provenance)
	}
}

func TestMemoryRepositoryReconcileOwnerCurrentStateABCDChronologyEndsRetired(t *testing.T) {
	ctx := context.Background()
	ownerID, _ := domainidentity.NewID("owner-1")
	repository := NewMemoryRepository()

	affirmA := promotedContext(t, "近期持續投入馬拉松訓練與耐力運動", "activity-a", "最近開始準備第一次全程馬拉松訓練")
	reaffirmC := promotedContext(t, "近期重新投入馬拉松訓練與耐力運動", "activity-c", "最近又開始規律準備馬拉松訓練")
	base := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)

	mutations := []domainsocialcontext.CurrentStateMutation{
		{SemanticID: "marathon", OccurredAt: base, Replacement: &affirmA},
		{SemanticID: "marathon", OccurredAt: base.Add(time.Hour), Replacement: nil},
		{SemanticID: "marathon", OccurredAt: base.Add(2 * time.Hour), Replacement: &reaffirmC},
		{SemanticID: "marathon", OccurredAt: base.Add(3 * time.Hour), Replacement: nil},
	}
	for _, mutation := range mutations {
		if _, err := repository.ReconcileOwnerCurrentState(ctx, ownerID, []domainsocialcontext.CurrentStateMutation{mutation}); err != nil {
			t.Fatalf("ReconcileOwnerCurrentState() error = %v", err)
		}
	}
	if got := repository.All(); len(got) != 0 {
		t.Fatalf("All() = %#v, want A→B→C→D chronology to end retired", got)
	}
}
