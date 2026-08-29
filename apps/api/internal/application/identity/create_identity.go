package identity

import (
	"context"
	"errors"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

var ErrIdentityAlreadyExists = errors.New("identity already exists")

type Repository interface {
	Create(ctx context.Context, identity domainidentity.Identity) error
}

type CreateIdentity struct {
	repository Repository
}

func NewCreateIdentity(repository Repository) CreateIdentity {
	return CreateIdentity{repository: repository}
}

type CreateIdentityCommand struct {
	ID string
}

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
