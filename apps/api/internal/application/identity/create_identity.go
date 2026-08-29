package identity

import (
	"context"
	"errors"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

// ErrIdentityAlreadyExists indicates that the requested identity ID is already in use.
var ErrIdentityAlreadyExists = errors.New("identity already exists")

// Repository persists identities for the create-identity use case.
type Repository interface {
	Create(ctx context.Context, identity domainidentity.Identity) error
}

// CreateIdentity creates and persists validated identities.
type CreateIdentity struct {
	repository Repository
}

// NewCreateIdentity constructs the create-identity use case.
func NewCreateIdentity(repository Repository) CreateIdentity {
	return CreateIdentity{repository: repository}
}

// CreateIdentityCommand contains the input required to create an identity.
type CreateIdentityCommand struct {
	ID string
}

// Execute validates the command, constructs the domain identity, and persists it.
func (uc CreateIdentity) Execute(ctx context.Context, command CreateIdentityCommand) (domainidentity.Identity, error) {
	id, err := domainidentity.NewID(command.ID)
	if err != nil {
		return domainidentity.Identity{}, err
	}

	created, err := domainidentity.New(id)
	if err != nil {
		return domainidentity.Identity{}, err
	}
	if err := uc.repository.Create(ctx, created); err != nil {
		return domainidentity.Identity{}, err
	}

	return created, nil
}
