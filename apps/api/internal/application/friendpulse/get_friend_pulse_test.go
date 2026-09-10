package friendpulse

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

type fakeFriendshipReader struct {
	active bool
}

func (f fakeFriendshipReader) IsActiveBetween(_ context.Context, _, _ domainidentity.ID) (bool, error) {
	return f.active, nil
}

type fakeCandidateReader struct {
	candidates []Candidate
}

func (f fakeCandidateReader) ListForOwner(_ context.Context, _ domainidentity.ID) ([]Candidate, error) {
	return append([]Candidate(nil), f.candidates...), nil
}

type fakeProjector struct {
	projections map[string]domainprivacy.ContextProjection
	notVisible  map[string]bool
	calls       *[]string
}

func (f fakeProjector) Project(
	_ context.Context,
	_, _ domainidentity.ID,
	socialContextID string,
) (domainprivacy.ContextProjection, bool, error) {
	if f.calls != nil {
		*f.calls = append(*f.calls, socialContextID)
	}
	if f.notVisible[socialContextID] {
		return domainprivacy.ContextProjection{}, false, nil
	}
	projection, ok := f.projections[socialContextID]
	return projection, ok, nil
}

func TestGetFriendPulseRejectsNonActiveFriendship(t *testing.T) {
	uc := NewGetFriendPulse(
		fakeFriendshipReader{active: false},
		fakeCandidateReader{},
		fakeProjector{},
	)

	_, err := uc.Execute(context.Background(), Query{
		AuthenticatedViewerID: "viewer",
		FriendID:              "friend",
	})
	if !errors.Is(err, ErrFriendPulseUnauthorized) {
		t.Fatalf("expected ErrFriendPulseUnauthorized, got %v", err)
	}
}

func TestGetFriendPulseProjectsCandidatesBeforeAnyRelevanceOrdering(t *testing.T) {
	base := time.Date(2026, time.September, 9, 8, 0, 0, 0, time.UTC)
	candidates := []Candidate{
		{SocialContextID: "first-source", SignalScore: 10, ObservedAt: base},
		{SocialContextID: "second-source", SignalScore: 100, ObservedAt: base.Add(time.Hour)},
	}
	calls := []string{}
	uc := NewGetFriendPulse(
		fakeFriendshipReader{active: true},
		fakeCandidateReader{candidates: candidates},
		fakeProjector{
			projections: map[string]domainprivacy.ContextProjection{
				"first-source":  {Meaning: "first"},
				"second-source": {Meaning: "second"},
			},
			calls: &calls,
		},
	)

	if _, err := uc.Execute(context.Background(), Query{AuthenticatedViewerID: "viewer", FriendID: "friend"}); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	want := []string{"first-source", "second-source"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("privacy projection call order = %#v, want source order %#v; candidate relevance metadata must not be consumed before projection", calls, want)
	}
}

func TestGetFriendPulseFiltersExpiredAndRepetitionSuppressedCandidatesBeforeProjection(t *testing.T) {
	now := time.Date(2026, time.September, 10, 1, 0, 0, 0, time.UTC)
	calls := []string{}
	uc := NewGetFriendPulse(
		fakeFriendshipReader{active: true},
		fakeCandidateReader{candidates: []Candidate{
			{SocialContextID: "expired", SignalScore: 100, EligibleUntil: now.Add(-time.Minute)},
			{SocialContextID: "suppressed", SignalScore: 100, RepetitionSuppressed: true},
			{SocialContextID: "eligible", SignalScore: 50, EligibleUntil: now.Add(time.Hour)},
		}},
		fakeProjector{
			projections: map[string]domainprivacy.ContextProjection{"eligible": {Meaning: "eligible"}},
			calls:       &calls,
		},
	)
	uc.now = func() time.Time { return now }

	pulse, err := uc.Execute(context.Background(), Query{AuthenticatedViewerID: "viewer", FriendID: "friend"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(pulse.Items) != 1 || pulse.Items[0].Projection.Meaning != "eligible" {
		t.Fatalf("pulse = %#v, want only eligible candidate", pulse)
	}
	if want := []string{"eligible"}; !reflect.DeepEqual(calls, want) {
		t.Fatalf("projection calls = %#v, want %#v; expired/suppressed candidates must be excluded by production eligibility behavior", calls, want)
	}
}

func TestGetFriendPulseReturnsAtMostThreeVisibleProjectionsInDeterministicSignalOrder(t *testing.T) {
	base := time.Date(2026, time.September, 5, 1, 0, 0, 0, time.UTC)
	candidates := []Candidate{
		{SocialContextID: "revoked", SignalScore: 100, ObservedAt: base.Add(40 * time.Minute)},
		{SocialContextID: "high-new", SignalScore: 90, ObservedAt: base.Add(20 * time.Minute)},
		{SocialContextID: "high-old", SignalScore: 90, ObservedAt: base.Add(10 * time.Minute)},
		{SocialContextID: "medium", SignalScore: 50, ObservedAt: base.Add(30 * time.Minute)},
		{SocialContextID: "low", SignalScore: 20, ObservedAt: base.Add(5 * time.Minute)},
	}
	projector := fakeProjector{
		projections: map[string]domainprivacy.ContextProjection{
			"low":      {Meaning: "low signal"},
			"high-new": {Meaning: "high signal newer"},
			"high-old": {Meaning: "high signal older"},
			"medium":   {Meaning: "medium signal"},
		},
		notVisible: map[string]bool{"revoked": true},
	}
	uc := NewGetFriendPulse(fakeFriendshipReader{active: true}, fakeCandidateReader{candidates: candidates}, projector)

	pulse, err := uc.Execute(context.Background(), Query{
		AuthenticatedViewerID: "viewer",
		FriendID:              "friend",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := make([]string, 0, len(pulse.Items))
	for _, item := range pulse.Items {
		got = append(got, item.Projection.Meaning)
	}
	want := []string{"high signal newer", "high signal older", "medium signal"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("pulse meanings = %#v, want %#v", got, want)
	}
	if len(pulse.Items) > 3 {
		t.Fatalf("pulse returned %d items, want <= 3", len(pulse.Items))
	}
}

func TestGetFriendPulseUsesViewerSpecificProjectionInsteadOfRawCandidateData(t *testing.T) {
	uc := NewGetFriendPulse(
		fakeFriendshipReader{active: true},
		fakeCandidateReader{candidates: []Candidate{{SocialContextID: "ctx-1", SignalScore: 80}}},
		fakeProjector{projections: map[string]domainprivacy.ContextProjection{
			"ctx-1": {Meaning: "safe relationship-specific projection"},
		}},
	)

	pulse, err := uc.Execute(context.Background(), Query{AuthenticatedViewerID: "viewer", FriendID: "friend"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(pulse.Items) != 1 || pulse.Items[0].Projection.Meaning != "safe relationship-specific projection" {
		t.Fatalf("unexpected pulse: %#v", pulse)
	}
}
