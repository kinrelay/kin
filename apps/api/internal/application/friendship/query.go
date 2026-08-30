package friendship

import (
	"context"
	"errors"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

var ErrFriendshipQueryUnauthorized = errors.New("requester is not a friendship participant")

// FriendshipReadModel is the minimal participant-facing projection for MVP 0.
type FriendshipReadModel struct {
	FirstParticipantID  string
	SecondParticipantID string
	Active              bool
}

// ReadRepository exposes read-only active friendship projections.
type ReadRepository interface {
	FindActiveBetween(ctx context.Context, first, second domainidentity.ID) (FriendshipReadModel, bool, error)
}

// GetFriendshipQuery identifies the relationship pair and the authenticated requester.
type GetFriendshipQuery struct {
	RequesterID         string
	FirstParticipantID  string
	SecondParticipantID string
}

// GetFriendship reads one active close-friend relationship without mutating domain state.
type GetFriendship struct {
	repository ReadRepository
}

// NewGetFriendship constructs the active relationship query use case.
func NewGetFriendship(repository ReadRepository) GetFriendship {
	return GetFriendship{repository: repository}
}

// Execute returns the active relationship projection only to one of its participants.
func (uc GetFriendship) Execute(ctx context.Context, query GetFriendshipQuery) (FriendshipReadModel, bool, error) {
	requesterID, err := domainidentity.NewID(query.RequesterID)
	if err != nil {
		return FriendshipReadModel{}, false, err
	}
	firstID, err := domainidentity.NewID(query.FirstParticipantID)
	if err != nil {
		return FriendshipReadModel{}, false, err
	}
	secondID, err := domainidentity.NewID(query.SecondParticipantID)
	if err != nil {
		return FriendshipReadModel{}, false, err
	}
	if requesterID != firstID && requesterID != secondID {
		return FriendshipReadModel{}, false, ErrFriendshipQueryUnauthorized
	}

	return uc.repository.FindActiveBetween(ctx, firstID, secondID)
}
