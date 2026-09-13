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

// InMemoryFriendPulseFixture supplies deterministic state for the runnable
// in-memory composition while keeping the default constructor deny-by-default.
type InMemoryFriendPulseFixture struct {
	ActiveFriendships [][2]string
	CandidatesByOwner map[string][]applicationfriendpulse.Candidate
	Projections       map[string]domainprivacy.ContextProjection
}

// NewInMemoryFriendPulseRuntime constructs an empty runnable Friend Pulse
// composition. Empty state is intentionally default-deny.
func NewInMemoryFriendPulseRuntime() InMemoryFriendPulseRuntime {
	return NewInMemoryFriendPulseRuntimeWithFixture(InMemoryFriendPulseFixture{})
}

// NewInMemoryFriendPulseRuntimeWithFixture constructs a runnable Friend Pulse
// composition with deterministic state for active friendships, owner candidates,
// and relationship-safe projections.
func NewInMemoryFriendPulseRuntimeWithFixture(fixture InMemoryFriendPulseFixture) InMemoryFriendPulseRuntime {
	friendships := &inMemoryActiveFriendships{active: make(map[[2]domainidentity.ID]bool)}
	for _, pair := range fixture.ActiveFriendships {
		first, firstErr := domainidentity.NewID(pair[0])
		second, secondErr := domainidentity.NewID(pair[1])
		if firstErr == nil && secondErr == nil {
			friendships.active[orderedIdentityPair(first, second)] = true
		}
	}

	candidates := &inMemoryPulseCandidates{byOwner: fixture.CandidatesByOwner}
	projector := &inMemoryContextProjector{projections: fixture.Projections}

	return InMemoryFriendPulseRuntime{
		FriendPulse: NewFriendPulseFlow(friendships, candidates, projector),
	}
}

type inMemoryActiveFriendships struct {
	active map[[2]domainidentity.ID]bool
}

func (f *inMemoryActiveFriendships) IsActiveBetween(_ context.Context, first, second domainidentity.ID) (bool, error) {
	return f.active[orderedIdentityPair(first, second)], nil
}

func orderedIdentityPair(first, second domainidentity.ID) [2]domainidentity.ID {
	if first > second {
		first, second = second, first
	}
	return [2]domainidentity.ID{first, second}
}

type inMemoryPulseCandidates struct {
	byOwner map[string][]applicationfriendpulse.Candidate
}

func (c *inMemoryPulseCandidates) ListForOwner(_ context.Context, ownerID domainidentity.ID) ([]applicationfriendpulse.Candidate, error) {
	return append([]applicationfriendpulse.Candidate(nil), c.byOwner[string(ownerID)]...), nil
}

type inMemoryContextProjector struct {
	projections map[string]domainprivacy.ContextProjection
}

func (p *inMemoryContextProjector) Project(_ context.Context, _ domainidentity.ID, _ domainidentity.ID, socialContextID string) (domainprivacy.ContextProjection, bool, error) {
	projection, visible := p.projections[socialContextID]
	return projection, visible, nil
}
