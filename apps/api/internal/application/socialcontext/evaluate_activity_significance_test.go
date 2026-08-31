package socialcontext

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

type significanceActivityReaderFake struct {
	items       []ActivityForSignificance
	err         error
	calls       int
	ownerAsked  domainidentity.ID
	activityIDs []string
}

func (f *significanceActivityReaderFake) ListOwnerPrivateNormalized(_ context.Context, ownerID domainidentity.ID, activityIDs []string) ([]ActivityForSignificance, error) {
	f.calls++
	f.ownerAsked = ownerID
	f.activityIDs = append([]string(nil), activityIDs...)
	if f.err != nil {
		return nil, f.err
	}
	return append([]ActivityForSignificance(nil), f.items...), nil
}

func significanceOwner(t *testing.T, value string) domainidentity.ID {
	t.Helper()
	id, err := domainidentity.NewID(value)
	if err != nil {
		t.Fatalf("NewID(%q): %v", value, err)
	}
	return id
}

func TestEvaluateActivitySignificanceReadsOnlyRequestedBatchOfRequesterPrivateNormalizedActivities(t *testing.T) {
	alice := significanceOwner(t, "alice")
	base := time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC)
	reader := &significanceActivityReaderFake{items: []ActivityForSignificance{
		// Deliberately out of occurrence order: application orchestration owns the
		// chronology contract rather than trusting an adapter's return order.
		{ID: "activity-3", OwnerID: alice, Content: " 最近持續研究 distributed systems 的 consistency trade-offs ", OccurredAt: base.Add(2 * time.Hour)},
		{ID: "activity-2", OwnerID: alice, Content: "看影片", OccurredAt: base.Add(time.Hour)},
		{ID: "activity-1", OwnerID: alice, Content: "最近持續研究 distributed systems 的 consistency trade-offs", OccurredAt: base},
	}}
	useCase := NewEvaluateActivitySignificance(reader)
	batch := []string{"activity-1", "activity-2", "activity-3"}

	decisions, err := useCase.Execute(context.Background(), EvaluateActivitySignificanceQuery{RequesterID: " alice ", ActivityIDs: batch})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if reader.calls != 1 || reader.ownerAsked != alice {
		t.Fatalf("reader calls/owner = %d/%q, want 1/%q", reader.calls, reader.ownerAsked, alice)
	}
	if !reflect.DeepEqual(reader.activityIDs, batch) {
		t.Fatalf("reader activity ids = %#v, want %#v", reader.activityIDs, batch)
	}
	if len(decisions) != 3 {
		t.Fatalf("decision count = %d, want 3", len(decisions))
	}
	if decisions[0].ActivityID != "activity-1" || decisions[0].Status != domainsocialcontext.SignificanceSuppressed || decisions[0].Reason != domainsocialcontext.SuppressionDuplicate {
		t.Fatalf("decision[0] = %#v, want older duplicate suppressed", decisions[0])
	}
	if decisions[1].ActivityID != "activity-2" || decisions[1].Status != domainsocialcontext.SignificanceSuppressed || decisions[1].Reason != domainsocialcontext.SuppressionLowSignal {
		t.Fatalf("decision[1] = %#v, want low-signal suppression", decisions[1])
	}
	if decisions[2].ActivityID != "activity-3" || decisions[2].Status != domainsocialcontext.SignificanceEligible {
		t.Fatalf("decision[2] = %#v, want newest equivalent signal eligible", decisions[2])
	}
}

func TestEvaluateActivitySignificanceRejectsInvalidRequesterBeforeRead(t *testing.T) {
	reader := &significanceActivityReaderFake{}
	useCase := NewEvaluateActivitySignificance(reader)

	_, err := useCase.Execute(context.Background(), EvaluateActivitySignificanceQuery{RequesterID: "  ", ActivityIDs: []string{"activity-1"}})
	if !errors.Is(err, domainidentity.ErrInvalidID) {
		t.Fatalf("Execute() error = %v, want %v", err, domainidentity.ErrInvalidID)
	}
	if reader.calls != 0 {
		t.Fatalf("reader calls = %d, want 0", reader.calls)
	}
}

func TestEvaluateActivitySignificanceRejectsReaderContractOwnerLeak(t *testing.T) {
	alice := significanceOwner(t, "alice")
	bob := significanceOwner(t, "bob")
	reader := &significanceActivityReaderFake{items: []ActivityForSignificance{
		{ID: "activity-bob", OwnerID: bob, Content: "Bob private signal must not enter Alice significance evaluation"},
	}}
	useCase := NewEvaluateActivitySignificance(reader)

	decisions, err := useCase.Execute(context.Background(), EvaluateActivitySignificanceQuery{RequesterID: string(alice), ActivityIDs: []string{"activity-bob"}})
	if !errors.Is(err, ErrActivityOwnerMismatch) {
		t.Fatalf("Execute() error = %v, want %v", err, ErrActivityOwnerMismatch)
	}
	if decisions != nil {
		t.Fatalf("decisions = %#v, want nil on owner-boundary violation", decisions)
	}
}

func TestEvaluateActivitySignificancePropagatesActivityReadFailure(t *testing.T) {
	readErr := errors.New("activity read unavailable")
	reader := &significanceActivityReaderFake{err: readErr}
	useCase := NewEvaluateActivitySignificance(reader)

	_, err := useCase.Execute(context.Background(), EvaluateActivitySignificanceQuery{RequesterID: "alice", ActivityIDs: []string{"activity-1"}})
	if !errors.Is(err, readErr) {
		t.Fatalf("Execute() error = %v, want %v", err, readErr)
	}
}

func TestEvaluateActivitySignificanceReturnsExplicitEmptyDecisionSetForEmptyBatch(t *testing.T) {
	reader := &significanceActivityReaderFake{}
	useCase := NewEvaluateActivitySignificance(reader)

	decisions, err := useCase.Execute(context.Background(), EvaluateActivitySignificanceQuery{RequesterID: "alice"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if reader.calls != 0 {
		t.Fatalf("reader calls = %d, want 0 for empty batch", reader.calls)
	}
	if decisions == nil || len(decisions) != 0 {
		t.Fatalf("decisions = %#v, want explicit empty collection", decisions)
	}
}
