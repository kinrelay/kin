package socialcontext

import (
	"context"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

// CurrentStateRepository owns the atomic consistency boundary for one owner's
// semantic Social Context lifecycle. A nil replacement means retire; an older or
// equal occurrence is a deterministic no-op.
type CurrentStateRepository interface {
	ReconcileCurrentState(
		ctx context.Context,
		ownerID domainidentity.ID,
		semanticID domainsocialcontext.SemanticIdentity,
		occurredAt time.Time,
		replacement *domainsocialcontext.SocialContext,
	) (bool, error)
}
