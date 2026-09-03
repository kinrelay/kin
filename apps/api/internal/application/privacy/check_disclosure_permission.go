package privacy

import (
	"context"
	"strings"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

// CheckDisclosurePermissionQuery identifies the disclosure decision to evaluate.
type CheckDisclosurePermissionQuery struct {
	OwnerID         string
	SocialContextID string
	ViewerID        string
}

// CheckDisclosurePermission evaluates a stored disclosure decision with default-deny semantics.
type CheckDisclosurePermission struct {
	repository DisclosureDecisionRepository
}

// NewCheckDisclosurePermission constructs the disclosure permission query use case.
func NewCheckDisclosurePermission(repository DisclosureDecisionRepository) CheckDisclosurePermission {
	return CheckDisclosurePermission{repository: repository}
}

// Execute returns false when no decision exists and otherwise returns the owner's stored decision.
func (uc CheckDisclosurePermission) Execute(
	ctx context.Context,
	query CheckDisclosurePermissionQuery,
) (bool, error) {
	ownerID, err := domainidentity.NewID(query.OwnerID)
	if err != nil {
		return false, err
	}
	viewerID, err := domainidentity.NewID(query.ViewerID)
	if err != nil {
		return false, err
	}
	socialContextID := strings.TrimSpace(query.SocialContextID)
	if socialContextID == "" {
		return false, domainprivacy.ErrBlankSocialContextID
	}

	decision, found, err := uc.repository.Find(ctx, ownerID, socialContextID, viewerID)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}
	return domainprivacy.AllowsDisclosure(&decision), nil
}
