package friendship

import (
	"context"
	"errors"

	domainfriendship "github.com/kinrelay/kin/apps/api/internal/domain/friendship"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

var (
	// ErrIdentityNotFound indicates that a participant or actor is not a known Kin identity.
	ErrIdentityNotFound = errors.New("identity not found")
	// ErrFriendshipAlreadyExists indicates that the participant pair already has a friendship aggregate.
	ErrFriendshipAlreadyExists = errors.New("friendship already exists between participants")
	// ErrFriendshipNotFound indicates that no invitation exists for the requested participant pair.
	ErrFriendshipNotFound = errors.New("friendship invitation not found")
)

// IdentityDirectory answers whether identities required by friendship commands exist.
type IdentityDirectory interface {
	Exists(ctx context.Context, id domainidentity.ID) (bool, error)
}

// Repository persists and retrieves friendship aggregates by their participant pair.
type Repository interface {
	FindBetween(ctx context.Context, first, second domainidentity.ID) (domainfriendship.Friendship, bool, error)
	Save(ctx context.Context, friendship domainfriendship.Friendship) error
}

// InviteFriendCommand expresses the intent to invite one distinct identity.
type InviteFriendCommand struct {
	InviterID string
	InviteeID string
}

// InviteFriend creates a pending friendship invitation.
type InviteFriend struct {
	identities IdentityDirectory
	repository Repository
}

// NewInviteFriend constructs the invite-friend use case.
func NewInviteFriend(identities IdentityDirectory, repository Repository) InviteFriend {
	return InviteFriend{identities: identities, repository: repository}
}

// Execute validates participants, enforces pair uniqueness, and persists a pending invitation.
func (uc InviteFriend) Execute(ctx context.Context, command InviteFriendCommand) (domainfriendship.Friendship, error) {
	inviterID, err := domainidentity.NewID(command.InviterID)
	if err != nil {
		return domainfriendship.Friendship{}, err
	}
	inviteeID, err := domainidentity.NewID(command.InviteeID)
	if err != nil {
		return domainfriendship.Friendship{}, err
	}

	invitation, err := domainfriendship.Invite(inviterID, inviteeID)
	if err != nil {
		return domainfriendship.Friendship{}, err
	}

	if err := requireIdentity(ctx, uc.identities, inviterID); err != nil {
		return domainfriendship.Friendship{}, err
	}
	if err := requireIdentity(ctx, uc.identities, inviteeID); err != nil {
		return domainfriendship.Friendship{}, err
	}

	_, found, err := uc.repository.FindBetween(ctx, inviterID, inviteeID)
	if err != nil {
		return domainfriendship.Friendship{}, err
	}
	if found {
		return domainfriendship.Friendship{}, ErrFriendshipAlreadyExists
	}
	if err := uc.repository.Save(ctx, invitation); err != nil {
		return domainfriendship.Friendship{}, err
	}

	return invitation, nil
}

// AcceptFriendshipCommand expresses who is accepting which invitation.
type AcceptFriendshipCommand struct {
	InviterID string
	InviteeID string
	ActorID   string
}

// AcceptFriendship activates an existing invitation through the domain invariant.
type AcceptFriendship struct {
	identities IdentityDirectory
	repository Repository
}

// NewAcceptFriendship constructs the accept-friendship use case.
func NewAcceptFriendship(identities IdentityDirectory, repository Repository) AcceptFriendship {
	return AcceptFriendship{identities: identities, repository: repository}
}

// Execute verifies identities, loads the single participant-pair aggregate, and asks it to accept.
func (uc AcceptFriendship) Execute(ctx context.Context, command AcceptFriendshipCommand) (domainfriendship.Friendship, error) {
	inviterID, err := domainidentity.NewID(command.InviterID)
	if err != nil {
		return domainfriendship.Friendship{}, err
	}
	inviteeID, err := domainidentity.NewID(command.InviteeID)
	if err != nil {
		return domainfriendship.Friendship{}, err
	}
	actorID, err := domainidentity.NewID(command.ActorID)
	if err != nil {
		return domainfriendship.Friendship{}, err
	}

	for _, id := range []domainidentity.ID{inviterID, inviteeID, actorID} {
		if err := requireIdentity(ctx, uc.identities, id); err != nil {
			return domainfriendship.Friendship{}, err
		}
	}

	existing, found, err := uc.repository.FindBetween(ctx, inviterID, inviteeID)
	if err != nil {
		return domainfriendship.Friendship{}, err
	}
	if !found {
		return domainfriendship.Friendship{}, ErrFriendshipNotFound
	}
	if err := existing.Accept(actorID); err != nil {
		return domainfriendship.Friendship{}, err
	}
	if err := uc.repository.Save(ctx, existing); err != nil {
		return domainfriendship.Friendship{}, err
	}

	return existing, nil
}

func requireIdentity(ctx context.Context, identities IdentityDirectory, id domainidentity.ID) error {
	exists, err := identities.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrIdentityNotFound
	}
	return nil
}
