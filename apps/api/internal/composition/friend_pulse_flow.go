package composition

import (
	"context"

	applicationfriendpulse "github.com/kinrelay/kin/apps/api/internal/application/friendpulse"
)

// FriendPulseFlow exposes the existing privacy-safe Friend Pulse query through
// the API composition boundary. Authentication identity is supplied by the
// delivery caller and is never accepted as separate viewer-controlled data.
type FriendPulseFlow struct {
	get applicationfriendpulse.GetFriendPulse
}

// NewFriendPulseFlow wires the concrete outer adapters for the three ports into
// the existing Friend Pulse application use case.
func NewFriendPulseFlow(
	friendships applicationfriendpulse.ActiveFriendshipReader,
	candidates applicationfriendpulse.CandidateReader,
	projector applicationfriendpulse.ContextProjector,
) FriendPulseFlow {
	return FriendPulseFlow{
		get: applicationfriendpulse.NewGetFriendPulse(friendships, candidates, projector),
	}
}

// Get is the authenticated delivery entry point for one friend's Pulse.
func (f FriendPulseFlow) Get(ctx context.Context, authenticatedViewerID, friendID string) (applicationfriendpulse.Pulse, error) {
	return f.get.Execute(ctx, applicationfriendpulse.Query{
		AuthenticatedViewerID: authenticatedViewerID,
		FriendID:              friendID,
	})
}
