package socialcontext

import (
	"context"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

type ReconcileCurrentSocialContextCommand struct {
	OwnerID   string
	Mutations []domainsocialcontext.CurrentStateMutation
}

type ReconcileCurrentSocialContext struct {
	repository CurrentStateRepository
}

func NewReconcileCurrentSocialContext(repository CurrentStateRepository) ReconcileCurrentSocialContext {
	return ReconcileCurrentSocialContext{repository: repository}
}

func (uc ReconcileCurrentSocialContext) Execute(ctx context.Context, command ReconcileCurrentSocialContextCommand) (int, error) {
	ownerID, err := domainidentity.NewID(command.OwnerID)
	if err != nil {
		return 0, err
	}
	return uc.repository.ReconcileOwnerCurrentState(ctx, ownerID, append([]domainsocialcontext.CurrentStateMutation(nil), command.Mutations...))
}
