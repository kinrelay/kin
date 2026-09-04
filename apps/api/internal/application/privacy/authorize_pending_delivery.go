package privacy

import (
	"context"
	"errors"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

var ErrDeliveryAuthorizationUnstable = errors.New("delivery authorization revisions kept changing")

// DeliveryAuthorizationSnapshot contains the current Kin-owned privacy inputs used to
// reauthorize one pending friend-facing delivery. The projection is recalculated from these
// inputs; callers must not trust the projection stored on the pending intent.
type DeliveryAuthorizationSnapshot struct {
	Meaning              string
	Decision             *domainprivacy.DisclosureDecision
	PrivacyRevision      domainprivacy.PrivacyPolicyRevision
	RelationshipActive   bool
	RelationshipRevision domainprivacy.RelationshipRevision
}

// DeliveryAuthorizationReader reads the current privacy decision and relationship state.
type DeliveryAuthorizationReader interface {
	Read(
		ctx context.Context,
		ownerID domainidentity.ID,
		viewerID domainidentity.ID,
		socialContextID string,
	) (DeliveryAuthorizationSnapshot, error)
}

// AtomicRevisionGuard is the Kin-owned capability boundary that couples revision validation
// with the transition to a dispatchable state. Concrete adapters may implement this with a
// transaction, CAS, lock, or another mechanism without leaking that technology inward.
type AtomicRevisionGuard interface {
	WithMatchingRevisions(
		ctx context.Context,
		ownerID domainidentity.ID,
		socialContextID string,
		viewerID domainidentity.ID,
		privacyRevision domainprivacy.PrivacyPolicyRevision,
		relationshipRevision domainprivacy.RelationshipRevision,
		commit func() error,
	) (bool, error)
}

// AuthorizePendingDelivery re-evaluates friend visibility immediately before a pending intent
// may become dispatchable. Revision races cause a fresh read and re-projection, never a stale
// commit.
type AuthorizePendingDelivery struct {
	reader DeliveryAuthorizationReader
	guard  AtomicRevisionGuard
}

func NewAuthorizePendingDelivery(reader DeliveryAuthorizationReader, guard AtomicRevisionGuard) AuthorizePendingDelivery {
	return AuthorizePendingDelivery{reader: reader, guard: guard}
}

func (uc AuthorizePendingDelivery) Execute(
	ctx context.Context,
	intent domainprivacy.PendingDeliveryIntent,
) (domainprivacy.PendingDeliveryIntent, error) {
	current := intent

	// A bounded retry keeps the application contract deterministic while allowing a revision
	// race to refresh once or twice without risking an unbounded loop under continuous churn.
	for attempt := 0; attempt < 3; attempt++ {
		snapshot, err := uc.reader.Read(ctx, intent.OwnerID(), intent.ViewerID(), intent.SocialContextID())
		if err != nil {
			return current, err
		}

		if !snapshot.RelationshipActive || !domainprivacy.AllowsDisclosure(snapshot.Decision) {
			return current.Cancel(), nil
		}

		projection, visible := domainprivacy.ProjectContext(snapshot.Meaning, snapshot.Decision)
		if !visible {
			return current.Cancel(), nil
		}

		current, err = current.RefreshProjection(
			projection,
			snapshot.PrivacyRevision,
			snapshot.RelationshipRevision,
		)
		if err != nil {
			return current, err
		}

		matched, err := uc.guard.WithMatchingRevisions(
			ctx,
			intent.OwnerID(),
			intent.SocialContextID(),
			intent.ViewerID(),
			snapshot.PrivacyRevision,
			snapshot.RelationshipRevision,
			func() error {
				dispatchable, transitionErr := current.MarkDispatchable()
				if transitionErr != nil {
					return transitionErr
				}
				current = dispatchable
				return nil
			},
		)
		if err != nil {
			return current, err
		}
		if matched {
			return current, nil
		}
	}

	return current.Cancel(), ErrDeliveryAuthorizationUnstable
}
