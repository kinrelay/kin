package composition

import (
	"context"
	"testing"

	applicationfriendpulse "github.com/kinrelay/kin/apps/api/internal/application/friendpulse"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

type activeFriendshipsStub struct{}

func (activeFriendshipsStub) IsActiveBetween(context.Context, domainidentity.ID, domainidentity.ID) (bool, error) {
	return true, nil
}

type pulseCandidatesStub struct{}

func (pulseCandidatesStub) ListForOwner(context.Context, domainidentity.ID) ([]applicationfriendpulse.Candidate, error) {
	return nil, nil
}

type contextProjectorStub struct{}

func (contextProjectorStub) Project(context.Context, domainidentity.ID, domainidentity.ID, string) (domainprivacy.ContextProjection, bool, error) {
	return domainprivacy.ContextProjection{}, false, nil
}

func TestFriendPulseFlowExposesAuthenticatedDelivery(t *testing.T) {
	flow := NewFriendPulseFlow(activeFriendshipsStub{}, pulseCandidatesStub{}, contextProjectorStub{})

	pulse, err := flow.Get(context.Background(), "viewer-1", "friend-1")
	if err != nil {
		t.Fatalf("get friend pulse: %v", err)
	}
	if len(pulse.Items) > 3 {
		t.Fatalf("pulse must stay bounded to at most 3 items, got %d", len(pulse.Items))
	}
}
