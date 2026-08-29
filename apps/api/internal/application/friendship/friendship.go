package friendship

import (
	"context"
	"errors"

	domainfriendship "github.com/kinrelay/kin/apps/api/internal/domain/friendship"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

var (
	ErrIdentityNotFound = errors.New("identity not found")
	ErrFriendshipAlreadyExists = errors.New("friendship already exists between participants")
	ErrFriendshipNotFound = errors.New("friendship invitation not found")
)

type IdentityDirectory interface {
	Exists(ctx context.Context, id domainidentity.ID) (bool, error)
}

type Repository interface {
	FindBetween(ctx context.Context, first, second domainidentity.ID) (domainfriendship.Friendship, bool, error)
	CreateIfAbsent(ctx context.Context, friendship domainfriendship.Friendship) (bool, error)
	UpdateBetween(ctx context.Context, first, second domainidentity.ID, update func(*domainfriendship.Friendship) error) (domainfriendship.Friendship, bool, error)
}

type InviteFriendCommand struct {
	InviterID string
	InviteeID string
}

type InviteFriend struct {
	identities IdentityDirectory
	repository Repository
}

func NewInviteFriend(identities IdentityDirectory, repository Repository) InviteFriend {
	return InviteFriend{identities: identities, repository: repository}
}

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
	created, err := uc.repository.CreateIfAbsent(ctx, invitation)
	if err != nil {
		return domainfriendship.Friendship{}, err
	}
	if !created {
		return domainfriendship.Friendship{}, ErrFriendshipAlreadyExists
	}
	return invitation, nil
}

type AcceptFriendshipCommand struct {
	InviterID string
	InviteeID string
	ActorID   string
}

type AcceptFriendship struct {
	identities IdentityDirectory
	repository Repository
}

func NewAcceptFriendship(identities IdentityDirectory, repository Repository) AcceptFriendship {
	return AcceptFriendship{identities: identities, repository: repository}
}

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
	updated, found, err := uc.repository.UpdateBetween(ctx, inviterID, inviteeID, func(friendship *domainfriendship.Friendship) error {
		return friendship.Accept(actorID)
	})
	if err != nil {
		return domainfriendship.Friendship{}, err
	}
	if !found {
		return domainfriendship.Friendship{}, ErrFriendshipNotFound
	}
	return updated, nil
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
