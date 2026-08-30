package socialcontext

import (
	"context"
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

type fakeContextActivityReader struct {
	activities []ActivityForContext
}

func (f fakeContextActivityReader) ListOwnerPrivateNormalized(context.Context, domainidentity.ID, []string) ([]ActivityForContext, error) {
	return append([]ActivityForContext(nil), f.activities...), nil
}

type fakeContextGenerator struct {
	calls int
	input ContextGenerationInput
	out   GeneratedContext
}

func (f *fakeContextGenerator) Generate(_ context.Context, input ContextGenerationInput) (GeneratedContext, error) {
	f.calls++
	f.input = input
	return f.out, nil
}

type fakeSocialContextRepository struct {
	saved []domainsocialcontext.SocialContext
}

func (f *fakeSocialContextRepository) Save(_ context.Context, context domainsocialcontext.SocialContext) error {
	f.saved = append(f.saved, context)
	return nil
}

func TestDeriveContextCandidateOnlySendsEligibleAuthorizedActivityToGenerator(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	generator := &fakeContextGenerator{out: GeneratedContext{
		Meaning:    "最近對分散式系統的可靠性與取捨特別有興趣",
		Provenance: []string{"activity-eligible"},
	}}
	repository := &fakeSocialContextRepository{}
	uc := NewDeriveContextCandidate(fakeContextActivityReader{activities: []ActivityForContext{
		{ID: "activity-eligible", OwnerID: ownerID, Content: "最近開始深入研究分散式系統設計"},
		{ID: "activity-low", OwnerID: ownerID, Content: "看文章"},
	}}, generator, repository)

	outcome, err := uc.Execute(context.Background(), DeriveContextCandidateCommand{
		RequesterID: "owner-1",
		ActivityIDs: []string{"activity-eligible", "activity-low"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if generator.calls != 1 {
		t.Fatalf("generator calls = %d, want 1", generator.calls)
	}
	if len(generator.input.Activities) != 1 || generator.input.Activities[0].ID != "activity-eligible" {
		t.Fatalf("generator activities = %#v, want only eligible activity", generator.input.Activities)
	}
	if outcome.Status != DerivationPromoted || len(repository.saved) != 1 {
		t.Fatalf("outcome = %#v, saved = %d; want promoted and persisted", outcome, len(repository.saved))
	}
}

func TestDeriveContextCandidateSuppressesWhenNoActivityIsEligible(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	generator := &fakeContextGenerator{}
	repository := &fakeSocialContextRepository{}
	uc := NewDeriveContextCandidate(fakeContextActivityReader{activities: []ActivityForContext{
		{ID: "activity-low", OwnerID: ownerID, Content: "看文章"},
	}}, generator, repository)

	outcome, err := uc.Execute(context.Background(), DeriveContextCandidateCommand{RequesterID: "owner-1", ActivityIDs: []string{"activity-low"}})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if outcome.Status != DerivationSuppressed || generator.calls != 0 || len(repository.saved) != 0 {
		t.Fatalf("outcome = %#v, generator calls = %d, saved = %d", outcome, generator.calls, len(repository.saved))
	}
}

func TestDeriveContextCandidateRejectsInvalidGeneratorCandidateWithoutPersistence(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	generator := &fakeContextGenerator{out: GeneratedContext{
		Meaning:    "最近開始深入研究分散式系統設計",
		Provenance: []string{"activity-1"},
	}}
	repository := &fakeSocialContextRepository{}
	uc := NewDeriveContextCandidate(fakeContextActivityReader{activities: []ActivityForContext{
		{ID: "activity-1", OwnerID: ownerID, Content: "最近開始深入研究分散式系統設計"},
	}}, generator, repository)

	outcome, err := uc.Execute(context.Background(), DeriveContextCandidateCommand{RequesterID: "owner-1", ActivityIDs: []string{"activity-1"}})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if outcome.Status != DerivationRejected || outcome.Reason != domainsocialcontext.ErrSourceReplay || len(repository.saved) != 0 {
		t.Fatalf("outcome = %#v, saved = %d; want rejected source replay without persistence", outcome, len(repository.saved))
	}
}

func TestDeriveContextCandidateRejectsGeneratorProvenanceOutsideEligibleSources(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	generator := &fakeContextGenerator{out: GeneratedContext{
		Meaning:    "最近對分散式系統的可靠性與取捨特別有興趣",
		Provenance: []string{"activity-unknown"},
	}}
	repository := &fakeSocialContextRepository{}
	uc := NewDeriveContextCandidate(fakeContextActivityReader{activities: []ActivityForContext{
		{ID: "activity-1", OwnerID: ownerID, Content: "最近開始深入研究分散式系統設計"},
	}}, generator, repository)

	outcome, err := uc.Execute(context.Background(), DeriveContextCandidateCommand{RequesterID: "owner-1", ActivityIDs: []string{"activity-1"}})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if outcome.Status != DerivationRejected || outcome.Reason != domainsocialcontext.ErrMissingContextProvenance || len(repository.saved) != 0 {
		t.Fatalf("outcome = %#v, saved = %d; want rejected provenance without persistence", outcome, len(repository.saved))
	}
}
