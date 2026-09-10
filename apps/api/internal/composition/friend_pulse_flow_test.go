package composition

import (
	"context"
	"reflect"
	"testing"
	"time"

	applicationfriendpulse "github.com/kinrelay/kin/apps/api/internal/application/friendpulse"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

type activeFriendshipsStub struct {
	first  domainidentity.ID
	second domainidentity.ID
}

func (s *activeFriendshipsStub) IsActiveBetween(_ context.Context, first, second domainidentity.ID) (bool, error) {
	s.first = first
	s.second = second
	return true, nil
}

type pulseCandidatesStub struct {
	owner domainidentity.ID
}

func (s *pulseCandidatesStub) ListForOwner(_ context.Context, owner domainidentity.ID) ([]applicationfriendpulse.Candidate, error) {
	s.owner = owner
	observedAt := time.Date(2026, time.September, 9, 0, 0, 0, 0, time.UTC)
	return []applicationfriendpulse.Candidate{
		{SocialContextID: "context-1", SignalScore: 40, ObservedAt: observedAt},
		{SocialContextID: "context-2", SignalScore: 30, ObservedAt: observedAt},
		{SocialContextID: "context-3", SignalScore: 20, ObservedAt: observedAt},
		{SocialContextID: "context-4", SignalScore: 10, ObservedAt: observedAt},
	}, nil
}

type projectionCall struct {
	viewer          domainidentity.ID
	owner           domainidentity.ID
	socialContextID string
}

type contextProjectorStub struct {
	calls []projectionCall
}

func (s *contextProjectorStub) Project(_ context.Context, viewer, owner domainidentity.ID, socialContextID string) (domainprivacy.ContextProjection, bool, error) {
	s.calls = append(s.calls, projectionCall{viewer: viewer, owner: owner, socialContextID: socialContextID})
	return domainprivacy.ContextProjection{Meaning: socialContextID}, true, nil
}

func TestFriendPulseFlowExposesAuthenticatedDelivery(t *testing.T) {
	friendships := &activeFriendshipsStub{}
	candidates := &pulseCandidatesStub{}
	projector := &contextProjectorStub{}
	flow := NewFriendPulseFlow(friendships, candidates, projector)

	pulse, err := flow.Get(context.Background(), "viewer-1", "friend-1")
	if err != nil {
		t.Fatalf("get friend pulse: %v", err)
	}
	if len(pulse.Items) != 3 {
		t.Fatalf("pulse must contain exactly the three highest-signal visible items, got %d", len(pulse.Items))
	}

	viewerID, err := domainidentity.NewID("viewer-1")
	if err != nil {
		t.Fatalf("viewer id: %v", err)
	}
	friendID, err := domainidentity.NewID("friend-1")
	if err != nil {
		t.Fatalf("friend id: %v", err)
	}
	if friendships.first != viewerID || friendships.second != friendID {
		t.Fatalf("friendship lookup received %q/%q, want %q/%q", friendships.first, friendships.second, viewerID, friendID)
	}
	if candidates.owner != friendID {
		t.Fatalf("candidate lookup owner = %q, want %q", candidates.owner, friendID)
	}

	// Privacy projection must run before relevance/ranking and bounding, so every
	// candidate is projected even though the delivered pulse is capped at three.
	wantCalls := []projectionCall{
		{viewer: viewerID, owner: friendID, socialContextID: "context-1"},
		{viewer: viewerID, owner: friendID, socialContextID: "context-2"},
		{viewer: viewerID, owner: friendID, socialContextID: "context-3"},
		{viewer: viewerID, owner: friendID, socialContextID: "context-4"},
	}
	if !reflect.DeepEqual(projector.calls, wantCalls) {
		t.Fatalf("projection calls = %#v, want %#v", projector.calls, wantCalls)
	}

	wantMeanings := []string{"context-1", "context-2", "context-3"}
	for i, item := range pulse.Items {
		if item.Projection.Meaning != wantMeanings[i] {
			t.Fatalf("pulse item %d meaning = %q, want %q", i, item.Projection.Meaning, wantMeanings[i])
		}
	}
}
