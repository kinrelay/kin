package socialcontext

import (
	"context"
	"reflect"
	"testing"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

type equalTimeContextReader struct {
	activities []ActivityForContext
}

func (r equalTimeContextReader) ListOwnerPrivateNormalized(context.Context, domainidentity.ID, []string) ([]ActivityForContext, error) {
	return append([]ActivityForContext(nil), r.activities...), nil
}

type capturingContextGenerator struct {
	ids []string
}

func (g *capturingContextGenerator) Generate(_ context.Context, input ContextGenerationInput) (GeneratedContext, error) {
	g.ids = g.ids[:0]
	for _, activity := range input.Activities {
		g.ids = append(g.ids, activity.ID)
	}
	return GeneratedContext{}, nil
}

type noopSocialContextRepository struct{}

func (noopSocialContextRepository) SaveIfAbsent(context.Context, domainidentity.ID, domainsocialcontext.SocialContext) (bool, error) {
	return true, nil
}

func TestDeriveContextCandidateUsesContributionChronologyWhenOccurredAtTies(t *testing.T) {
	ownerID, err := domainidentity.NewID("owner-1")
	if err != nil {
		t.Fatalf("NewID() error = %v", err)
	}
	occurredAt := time.Date(2026, 8, 31, 1, 0, 0, 0, time.UTC)
	olderContribution := occurredAt.Add(time.Minute)
	newerContribution := occurredAt.Add(2 * time.Minute)

	reader := equalTimeContextReader{activities: []ActivityForContext{
		{ID: "activity-newer", OwnerID: ownerID, Content: "最近開始準備第一次全程馬拉松訓練", OccurredAt: occurredAt, ContributedAt: newerContribution},
		{ID: "activity-older", OwnerID: ownerID, Content: "最近開始深入研究分散式系統設計", OccurredAt: occurredAt, ContributedAt: olderContribution},
	}}
	generator := &capturingContextGenerator{}
	useCase := NewDeriveContextCandidate(reader, generator, noopSocialContextRepository{})

	_, err = useCase.Execute(context.Background(), DeriveContextCandidateCommand{
		RequesterID: "owner-1",
		ActivityIDs: []string{"activity-newer", "activity-older"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if want := []string{"activity-older", "activity-newer"}; !reflect.DeepEqual(generator.ids, want) {
		t.Fatalf("generator activity order = %#v, want contribution chronology %#v", generator.ids, want)
	}
}
