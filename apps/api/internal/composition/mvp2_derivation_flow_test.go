package composition

import (
	"context"
	"sync"
	"testing"
	"time"

	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
	domainactivity "github.com/kinrelay/kin/apps/api/internal/domain/activity"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func TestMVP2DerivationFlowPromotesExistingPrivateActivityAndExposesOwnerReadModel(t *testing.T) {
	ctx := context.Background()
	flow := NewMVP2DerivationFlow()
	ownerID, err := domainidentity.NewID("alice")
	if err != nil {
		t.Fatalf("NewID(): %v", err)
	}
	content, err := domainactivity.NewContent("最近開始深入研究分散式系統設計與一致性取捨")
	if err != nil {
		t.Fatalf("NewContent(): %v", err)
	}
	now := time.Date(2026, time.August, 30, 10, 30, 0, 0, time.UTC)
	activity, err := domainactivity.NewManual("activity-1", ownerID, content, now, now)
	if err != nil {
		t.Fatalf("NewManual(): %v", err)
	}
	if err := flow.Activities.Save(ctx, activity); err != nil {
		t.Fatalf("Save activity: %v", err)
	}

	outcome, err := flow.Derive.Execute(ctx, applicationsocialcontext.DeriveContextCandidateCommand{
		RequesterID: "alice",
		ActivityIDs: []string{"activity-1"},
	})
	if err != nil {
		t.Fatalf("Derive.Execute() error = %v", err)
	}
	if outcome.Status != applicationsocialcontext.DerivationPromoted {
		t.Fatalf("derivation status = %q, want promoted; reason=%v", outcome.Status, outcome.Reason)
	}

	contexts, err := flow.List.Execute(ctx, applicationsocialcontext.ListMySocialContextsQuery{RequesterID: "alice"})
	if err != nil {
		t.Fatalf("List.Execute() error = %v", err)
	}
	if len(contexts) != 1 {
		t.Fatalf("context count = %d, want 1: %#v", len(contexts), contexts)
	}
	if contexts[0].Meaning == content.String() {
		t.Fatalf("derived meaning replayed raw Activity: %q", contexts[0].Meaning)
	}
	if len(contexts[0].Provenance) != 1 || contexts[0].Provenance[0] != "activity-1" {
		t.Fatalf("provenance = %#v, want activity-1", contexts[0].Provenance)
	}
}

func TestMVP2DerivationFlowSuppressesLowSignalWithoutPersistingContext(t *testing.T) {
	ctx := context.Background()
	flow := NewMVP2DerivationFlow()
	ownerID, _ := domainidentity.NewID("alice")
	content, _ := domainactivity.NewContent("看文章")
	now := time.Date(2026, time.August, 30, 10, 30, 0, 0, time.UTC)
	activity, err := domainactivity.NewManual("activity-low", ownerID, content, now, now)
	if err != nil {
		t.Fatalf("NewManual(): %v", err)
	}
	if err := flow.Activities.Save(ctx, activity); err != nil {
		t.Fatalf("Save activity: %v", err)
	}

	outcome, err := flow.Derive.Execute(ctx, applicationsocialcontext.DeriveContextCandidateCommand{
		RequesterID: "alice",
		ActivityIDs: []string{"activity-low"},
	})
	if err != nil {
		t.Fatalf("Derive.Execute() error = %v", err)
	}
	if outcome.Status != applicationsocialcontext.DerivationSuppressed {
		t.Fatalf("derivation status = %q, want suppressed", outcome.Status)
	}
	contexts, err := flow.List.Execute(ctx, applicationsocialcontext.ListMySocialContextsQuery{RequesterID: "alice"})
	if err != nil {
		t.Fatalf("List.Execute() error = %v", err)
	}
	if len(contexts) != 0 {
		t.Fatalf("context count = %d, want 0", len(contexts))
	}
}

func TestMVP2DerivationFlowSuppressesRetryWithoutPersistingDuplicateContext(t *testing.T) {
	ctx := context.Background()
	flow := NewMVP2DerivationFlow()
	ownerID, _ := domainidentity.NewID("alice")
	content, _ := domainactivity.NewContent("最近開始深入研究分散式系統設計與一致性取捨")
	now := time.Date(2026, time.August, 30, 10, 30, 0, 0, time.UTC)
	activity, err := domainactivity.NewManual("activity-retry", ownerID, content, now, now)
	if err != nil {
		t.Fatalf("NewManual(): %v", err)
	}
	if err := flow.Activities.Save(ctx, activity); err != nil {
		t.Fatalf("Save activity: %v", err)
	}

	command := applicationsocialcontext.DeriveContextCandidateCommand{
		RequesterID: "alice",
		ActivityIDs: []string{"activity-retry"},
	}
	first, err := flow.Derive.Execute(ctx, command)
	if err != nil {
		t.Fatalf("first Derive.Execute() error = %v", err)
	}
	if first.Status != applicationsocialcontext.DerivationPromoted {
		t.Fatalf("first derivation status = %q, want promoted", first.Status)
	}
	second, err := flow.Derive.Execute(ctx, command)
	if err != nil {
		t.Fatalf("second Derive.Execute() error = %v", err)
	}
	if second.Status != applicationsocialcontext.DerivationSuppressed {
		t.Fatalf("second derivation status = %q, want suppressed duplicate", second.Status)
	}

	contexts, err := flow.List.Execute(ctx, applicationsocialcontext.ListMySocialContextsQuery{RequesterID: "alice"})
	if err != nil {
		t.Fatalf("List.Execute() error = %v", err)
	}
	if len(contexts) != 1 {
		t.Fatalf("context count = %d, want exactly 1 after retry", len(contexts))
	}
}

func TestMVP2DerivationFlowAtomicallySuppressesConcurrentDuplicateDerivations(t *testing.T) {
	ctx := context.Background()
	flow := NewMVP2DerivationFlow()
	ownerID, _ := domainidentity.NewID("alice")
	content, _ := domainactivity.NewContent("最近開始深入研究分散式系統設計與一致性取捨")
	now := time.Date(2026, time.August, 30, 10, 30, 0, 0, time.UTC)
	activity, err := domainactivity.NewManual("activity-concurrent", ownerID, content, now, now)
	if err != nil {
		t.Fatalf("NewManual(): %v", err)
	}
	if err := flow.Activities.Save(ctx, activity); err != nil {
		t.Fatalf("Save activity: %v", err)
	}

	command := applicationsocialcontext.DeriveContextCandidateCommand{
		RequesterID: "alice",
		ActivityIDs: []string{"activity-concurrent"},
	}
	const workers = 16
	outcomes := make(chan applicationsocialcontext.DerivationOutcome, workers)
	errs := make(chan error, workers)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			outcome, err := flow.Derive.Execute(ctx, command)
			if err != nil {
				errs <- err
				return
			}
			outcomes <- outcome
		}()
	}
	close(start)
	wg.Wait()
	close(outcomes)
	close(errs)

	for err := range errs {
		t.Fatalf("concurrent Derive.Execute() error = %v", err)
	}
	promoted := 0
	suppressed := 0
	for outcome := range outcomes {
		switch outcome.Status {
		case applicationsocialcontext.DerivationPromoted:
			promoted++
		case applicationsocialcontext.DerivationSuppressed:
			suppressed++
		default:
			t.Fatalf("unexpected concurrent derivation status = %q, reason=%v", outcome.Status, outcome.Reason)
		}
	}
	if promoted != 1 || suppressed != workers-1 {
		t.Fatalf("promoted=%d suppressed=%d, want 1 promoted and %d suppressed", promoted, suppressed, workers-1)
	}

	contexts, err := flow.List.Execute(ctx, applicationsocialcontext.ListMySocialContextsQuery{RequesterID: "alice"})
	if err != nil {
		t.Fatalf("List.Execute() error = %v", err)
	}
	if len(contexts) != 1 {
		t.Fatalf("context count = %d, want exactly 1 after concurrent derivation", len(contexts))
	}
}
