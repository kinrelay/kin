package privacy

import (
	"errors"
	"strings"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

var (
	ErrBlankDeliverySocialContextID = errors.New("delivery social context id is blank")
	ErrInvalidPrivacyPolicyRevision = errors.New("privacy policy revision must be positive")
	ErrInvalidRelationshipRevision  = errors.New("relationship revision must be positive")
	ErrDeliveryIntentNotPending     = errors.New("delivery intent is not pending")
)

type PrivacyPolicyRevision uint64

type RelationshipRevision uint64

type DeliveryIntentState string

const (
	DeliveryIntentPending      DeliveryIntentState = "pending"
	DeliveryIntentDispatchable DeliveryIntentState = "dispatchable"
	DeliveryIntentCancelled    DeliveryIntentState = "cancelled"
)

// PendingDeliveryIntent is a Kin-owned privacy snapshot waiting for dispatch authorization.
// The projection is never sufficient authorization by itself; the bound revisions exist so
// application orchestration can detect stale privacy or relationship state before dispatch.
type PendingDeliveryIntent struct {
	ownerID             domainidentity.ID
	viewerID            domainidentity.ID
	socialContextID      string
	projection           ContextProjection
	privacyRevision      PrivacyPolicyRevision
	relationshipRevision RelationshipRevision
	state                DeliveryIntentState
}

func NewPendingDeliveryIntent(
	ownerID domainidentity.ID,
	viewerID domainidentity.ID,
	socialContextID string,
	projection ContextProjection,
	privacyRevision PrivacyPolicyRevision,
	relationshipRevision RelationshipRevision,
) (PendingDeliveryIntent, error) {
	normalizedOwnerID, err := domainidentity.NewID(string(ownerID))
	if err != nil {
		return PendingDeliveryIntent{}, err
	}
	normalizedViewerID, err := domainidentity.NewID(string(viewerID))
	if err != nil {
		return PendingDeliveryIntent{}, err
	}
	contextID := strings.TrimSpace(socialContextID)
	if contextID == "" {
		return PendingDeliveryIntent{}, ErrBlankDeliverySocialContextID
	}
	if err := validateDeliveryRevisions(privacyRevision, relationshipRevision); err != nil {
		return PendingDeliveryIntent{}, err
	}

	return PendingDeliveryIntent{
		ownerID:             normalizedOwnerID,
		viewerID:            normalizedViewerID,
		socialContextID:      contextID,
		projection:           projection,
		privacyRevision:      privacyRevision,
		relationshipRevision: relationshipRevision,
		state:                DeliveryIntentPending,
	}, nil
}

// RefreshProjection replaces the still-pending snapshot with a projection calculated from
// current privacy and relationship state. It never authorizes dispatch by itself.
func (i PendingDeliveryIntent) RefreshProjection(
	projection ContextProjection,
	privacyRevision PrivacyPolicyRevision,
	relationshipRevision RelationshipRevision,
) (PendingDeliveryIntent, error) {
	if i.state != DeliveryIntentPending {
		return PendingDeliveryIntent{}, ErrDeliveryIntentNotPending
	}
	if err := validateDeliveryRevisions(privacyRevision, relationshipRevision); err != nil {
		return PendingDeliveryIntent{}, err
	}
	i.projection = projection
	i.privacyRevision = privacyRevision
	i.relationshipRevision = relationshipRevision
	return i, nil
}

// MarkDispatchable is deliberately invoked only from inside an application-owned atomic
// revision guard after the latest privacy and relationship revisions have been matched.
func (i PendingDeliveryIntent) MarkDispatchable() (PendingDeliveryIntent, error) {
	if i.state != DeliveryIntentPending {
		return PendingDeliveryIntent{}, ErrDeliveryIntentNotPending
	}
	i.state = DeliveryIntentDispatchable
	return i, nil
}

// Cancel prevents a pending intent from becoming dispatchable after authorization is denied.
func (i PendingDeliveryIntent) Cancel() PendingDeliveryIntent {
	if i.state == DeliveryIntentPending {
		i.state = DeliveryIntentCancelled
	}
	return i
}

func (i PendingDeliveryIntent) OwnerID() domainidentity.ID { return i.ownerID }
func (i PendingDeliveryIntent) ViewerID() domainidentity.ID { return i.viewerID }
func (i PendingDeliveryIntent) SocialContextID() string { return i.socialContextID }
func (i PendingDeliveryIntent) Projection() ContextProjection { return i.projection }
func (i PendingDeliveryIntent) PrivacyRevision() PrivacyPolicyRevision { return i.privacyRevision }
func (i PendingDeliveryIntent) RelationshipRevision() RelationshipRevision { return i.relationshipRevision }
func (i PendingDeliveryIntent) State() DeliveryIntentState { return i.state }

func validateDeliveryRevisions(privacyRevision PrivacyPolicyRevision, relationshipRevision RelationshipRevision) error {
	if privacyRevision == 0 {
		return ErrInvalidPrivacyPolicyRevision
	}
	if relationshipRevision == 0 {
		return ErrInvalidRelationshipRevision
	}
	return nil
}
