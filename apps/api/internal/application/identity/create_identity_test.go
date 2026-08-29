package identity

import (
	"context"
	"errors"
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

type fakeRepository struct {
	created map[domainidentity.ID]domainidentity.Identity
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{created: make(map[domainidentity.ID]domainidentity.Identity)}
}

func (r *fakeRepository) Create(_ context.Context, item domainidentity.Identity) error {
	if _, exists := r.created[item.ID()]; exists {
		return ErrIdentityAlreadyExists
	}

	r.created[item.ID()] = item
	return nil
}

func TestCreateIdentityExecute(t *testing.T) {
	t.Parallel()

	repository := newFakeRepository()
	useCase := NewCreateIdentity(repository)

	created, err := useCase.Execute(context.Background(), CreateIdentityCommand{ID: "user-123"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if created.ID() != domainidentity.ID("user-123") {
		t.Fatalf("created ID = %q, want %q", created.ID(), "user-123")
	}
	if _, exists := repository.created[created.ID()]; !exists {
		t.Fatal("repository did not receive created identity")
	}
}

func TestCreateIdentityRejectsInvalidID(t *testing.T) {
	t.Parallel()

	useCase := NewCreateIdentity(newFakeRepository())

	_, err := useCase.Execute(context.Background(), CreateIdentityCommand{ID: "   "})
	if !errors.Is(err, domainidentity.ErrInvalidID) {
		t.Fatalf("Execute() error = %v, want %v", err, domainidentity.ErrInvalidID)
	}
}

func TestCreateIdentityRejectsDuplicateID(t *testing.T) {
	t.Parallel()

	useCase := NewCreateIdentity(newFakeRepository())
	command := CreateIdentityCommand{ID: "user-123"}

	if _, err := useCase.Execute(context.Background(), command); err != nil {
		t.Fatalf("first Execute() error = %v", err)
	}

	_, err := useCase.Execute(context.Background(), command)
	if !errors.Is(err, ErrIdentityAlreadyExists) {
		t.Fatalf("second Execute() error = %v, want %v", err, ErrIdentityAlreadyExists)
	}
}
