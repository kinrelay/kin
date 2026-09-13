package composition

import (
	"context"
	"errors"
	"testing"

	applicationfriendpulse "github.com/kinrelay/kin/apps/api/internal/application/friendpulse"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

func TestNewInMemoryFriendPulseRuntimeProvidesRunnableDelivery(t *testing.T) {
	runtime := NewInMemoryFriendPulseRuntime()

	_, err := runtime.FriendPulse.Get(context.Background(), "viewer-1", "friend-1")
	if !errors.Is(err, applicationfriendpulse.ErrFriendPulseUnauthorized) {
		t.Fatalf("expected wired runtime to enforce active friendship, got %v", err)
	}
}

func TestNewInMemoryFriendPulseRuntimeSupportsAuthorizedDelivery(t *testing.T) {
	fixture := InMemoryFriendPulseFixture{
		ActiveFriendships: [][2]string{{"viewer-1", "friend-1"}},
		CandidatesByOwner: map[string][]applicationfriendpulse.Candidate{
			"friend-1": []applicationfriendpulse.Candidate{
				{SocialContextID: "context-1", SignalScore: 10},
			},
		},
		Projections: map[string]domainprivacy.ContextProjection{
			"context-1": domainprivacy.ContextProjection{Meaning: "friend-visible context"},
		},
	}
	runtime := NewInMemoryFriendPulseRuntimeWithFixture(fixture)

	pulse, err := runtime.FriendPulse.Get(context.Background(), "viewer-1", "friend-1")
	if err != nil {
		t.Fatalf("expected authorized Friend Pulse delivery, got %v", err)
	}
	if len(pulse.Items) != 1 || pulse.Items[0].Projection.Meaning != "friend-visible context" {
		t.Fatalf("expected seeded friend-visible projection, got %#v", pulse.Items)
	}
}
