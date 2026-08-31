package socialcontext

import (
	"context"
	"testing"
	"time"

	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

func TestEvaluateActivitySignificanceUsesContributionChronologyWhenOccurrencesTie(t *testing.T) {
	alice := significanceOwner(t, "alice")
	occurredAt := time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC)
	reader := &significanceActivityReaderFake{items: []ActivityForSignificance{
		{ID: "activity-new", OwnerID: alice, Content: "最近持續研究 distributed systems 的 consistency trade-offs", OccurredAt: occurredAt, ContributedAt: occurredAt.Add(2 * time.Minute)},
		{ID: "activity-old", OwnerID: alice, Content: "最近持續研究 distributed systems 的 consistency trade-offs", OccurredAt: occurredAt, ContributedAt: occurredAt.Add(time.Minute)},
	}}
	useCase := NewEvaluateActivitySignificance(reader)

	decisions, err := useCase.Execute(context.Background(), EvaluateActivitySignificanceQuery{RequesterID: "alice", ActivityIDs: []string{"activity-old", "activity-new"}})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(decisions) != 2 {
		t.Fatalf("decision count = %d, want 2", len(decisions))
	}
	if decisions[0].ActivityID != "activity-old" || decisions[0].Status != domainsocialcontext.SignificanceSuppressed || decisions[0].Reason != domainsocialcontext.SuppressionDuplicate {
		t.Fatalf("decision[0] = %#v, want older contribution suppressed", decisions[0])
	}
	if decisions[1].ActivityID != "activity-new" || decisions[1].Status != domainsocialcontext.SignificanceEligible {
		t.Fatalf("decision[1] = %#v, want newer contribution eligible", decisions[1])
	}
}
