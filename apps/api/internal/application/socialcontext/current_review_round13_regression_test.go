package socialcontext

import (
	"context"
	"testing"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func TestDeriveContextCandidateOmitsBlockedAbandonmentFromCompoundGeneratorInput(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	baseTime := time.Date(2026, 8, 31, 20, 0, 0, 0, time.UTC)
	generator := &fakeContextGenerator{out: GeneratedContext{
		Meaning:    "近期持續關注分散式系統，也開始投入耐力運動準備",
		Provenance: []string{"activity-distributed", "activity-compound"},
	}}
	repository := &fakeSocialContextRepository{}
	uc := NewDeriveContextCandidate(fakeContextActivityReader{activities: []ActivityForContext{
		{ID: "activity-distributed", OwnerID: ownerID, Content: "最近開始深入研究分散式系統設計", OccurredAt: baseTime},
		{ID: "activity-french", OwnerID: ownerID, Content: "最近開始學習法文", OccurredAt: baseTime.Add(time.Minute)},
		{ID: "activity-compound", OwnerID: ownerID, Content: "後來放棄了，最近開始準備馬拉松比賽", OccurredAt: baseTime.Add(2 * time.Minute)},
	}}, generator, repository)

	outcome, err := uc.Execute(context.Background(), DeriveContextCandidateCommand{
		RequesterID: "owner-1",
		ActivityIDs: []string{"activity-distributed", "activity-french", "activity-compound"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if outcome.Status != DerivationPromoted {
		t.Fatalf("outcome status = %q, want promoted", outcome.Status)
	}
	if len(generator.input.Activities) != 2 {
		t.Fatalf("generator activities = %#v, want distributed and sanitized compound activities", generator.input.Activities)
	}
	if got := generator.input.Activities[1].Content; got != "最近開始準備馬拉松比賽" {
		t.Fatalf("compound generator content = %q, want blocked abandonment removed", got)
	}
}
