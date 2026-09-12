package composition

import (
	"context"

	applicationfriendpulse "github.com/kinrelay/kin/apps/api/internal/application/friendpulse"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

// InMemoryFriendPulseRuntime is the provider-neutral runnable composition root
// for Friend Pulse. Concrete delivery mechanisms (HTTP, workers, etc.) can call
// FriendPulse without owning application policy.
type InMemoryFriendPulseRuntime struct {
	FriendPulse FriendPulseFlow
}

// NewInMemoryFriendPulseRuntime constructs a runnable Friend Pulse composition
// using deterministic in-memory ports. The ports intentionally start empty;
// production delivery can populate or replace them without changing the use case.
func NewInMemoryFriendPulseRuntime() InMemoryFriendPulseRuntime {
	friendships := &inMemoryActiveFriendships{}
	candidates := &inMemoryPulseCandidates{}
	projector := &inMemoryContextProjector{}

	return InMemoryFriendPulseRuntime{
		FriendPulse: NewFriendPulseFlow(friendships, candidates, projector),
	}
}

type inMemoryActiveFriendships struct{}

func (*inMemoryActiveFriendships) IsActiveBetween(context.Context, domainidentity.ID, domainidentity.ID) (bool, error) {
	return false, nil
}

type inMemoryPulseCandidates struct{}

func (*inMemoryPulseCandidates) ListForOwner(context.Context, domainidentity.ID) ([]applicationfriendpulse.Candidate, error) {
	return nil, nil
}

type inMemoryContextProjector struct{}

func (*inMemoryContextProjector) Project(context.Context, domainidentity.ID, domainidentity.ID, string) (domainprivacy.ContextProjection, bool, error) {
	return domainprivacy.ContextProjection{}, false, nil
}
