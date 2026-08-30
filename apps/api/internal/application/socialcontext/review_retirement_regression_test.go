package socialcontext

import (
	"context"
	"reflect"
	"testing"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

type retirementRecordingRepository struct {
	retired [][]string
	saved   int
}

func (r *retirementRecordingRepository) SaveIfAbsent(context.Context, domainidentity.ID, domainsocialcontext.SocialContext) (bool, error) {
	r.saved++
	return true, nil
}

func (r *retirementRecordingRepository) RetireByProvenance(_ context.Context, _ domainidentity.ID, activityIDs []string) (int, error) {
	r.retired = append(r.retired, append([]string(nil), activityIDs...))
	return 1, nil
}

func TestDeriveContextCandidateRetiresPreviouslyPersistedContextWhenLaterReversalWins(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	base := time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC)
	generator := &fakeContextGenerator{out: GeneratedContext{
		RetiredProvenance: []string{"activity-start"},
	}}
	repository := &retirementRecordingRepository{}
	uc := NewDeriveContextCandidate(fakeContextActivityReader{activities: []ActivityForContext{
		{ID: "activity-start", OwnerID: ownerID, Content: "最近開始研究分散式系統與可靠性工程取捨", OccurredAt: base},
		{ID: "activity-stop", OwnerID: ownerID, Content: "不再研究分散式系統", OccurredAt: base.Add(time.Hour)},
	}}, generator, repository)

	outcome, err := uc.Execute(context.Background(), DeriveContextCandidateCommand{
		RequesterID: "owner-1",
		ActivityIDs: []string{"activity-start", "activity-stop"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if outcome.Status != DerivationSuppressed {
		t.Fatalf("outcome = %#v, want suppressed after retirement leaves no current context", outcome)
	}
	if repository.saved != 0 {
		t.Fatalf("saved = %d, want no new context", repository.saved)
	}
	if want := [][]string{{"activity-start"}}; !reflect.DeepEqual(repository.retired, want) {
		t.Fatalf("retired = %#v, want %#v", repository.retired, want)
	}
}
