package socialcontext

import (
	"context"
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

func TestDeriveContextCandidateRejectsMalformedBlankGeneratorResult(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	generator := &fakeContextGenerator{out: GeneratedContext{
		Meaning:    "",
		Provenance: []string{"activity-1"},
	}}
	repository := &fakeSocialContextRepository{}
	uc := NewDeriveContextCandidate(fakeContextActivityReader{activities: []ActivityForContext{{
		ID:      "activity-1",
		OwnerID: ownerID,
		Content: "最近開始深入研究分散式系統設計",
	}}}, generator, repository)

	outcome, err := uc.Execute(context.Background(), DeriveContextCandidateCommand{
		RequesterID: "owner-1",
		ActivityIDs: []string{"activity-1"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if outcome.Status != DerivationRejected || outcome.Reason != domainsocialcontext.ErrBlankContextMeaning {
		t.Fatalf("outcome = %#v, want rejected blank generator result", outcome)
	}
	if repository.saveIfAbsentCalls != 0 {
		t.Fatalf("SaveIfAbsent calls = %d, want 0", repository.saveIfAbsentCalls)
	}
}

func TestDeriveContextCandidateSuppressesExplicitNoCandidateGeneratorResult(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	generator := &fakeContextGenerator{out: GeneratedContext{}}
	repository := &fakeSocialContextRepository{}
	uc := NewDeriveContextCandidate(fakeContextActivityReader{activities: []ActivityForContext{{
		ID:      "activity-1",
		OwnerID: ownerID,
		Content: "最近開始深入研究分散式系統設計",
	}}}, generator, repository)

	outcome, err := uc.Execute(context.Background(), DeriveContextCandidateCommand{
		RequesterID: "owner-1",
		ActivityIDs: []string{"activity-1"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if outcome.Status != DerivationSuppressed || outcome.Reason != nil {
		t.Fatalf("outcome = %#v, want explicit no-candidate suppression", outcome)
	}
	if repository.saveIfAbsentCalls != 0 {
		t.Fatalf("SaveIfAbsent calls = %d, want 0", repository.saveIfAbsentCalls)
	}
}
