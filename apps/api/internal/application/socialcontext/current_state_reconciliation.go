package socialcontext

import (
	"context"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

// CurrentStateRepository owns the atomic consistency boundary for one owner's
// semantic Social Context lifecycle. All mutations for one application command
// cross the persistence boundary together.
type CurrentStateRepository interface {
	ReconcileOwnerCurrentState(
		ctx context.Context,
		ownerID domainidentity.ID,
		mutations []domainsocialcontext.CurrentStateMutation,
	) (int, error)
}
