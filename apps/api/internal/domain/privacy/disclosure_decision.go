package privacy

import (
	"errors"
	"strings"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

var (
	ErrOnlyContextOwnerCanChangeDisclosure = errors.New("only the context owner can change disclosure")
	ErrBlankSocialContextID                 = errors.New("social context id is blank")
	ErrInvalidVisibility                    = errors.New("disclosure visibility is invalid")
)

type Visibility string

const (
	VisibilityHidden  Visibility = "hidden"
	VisibilityVisible Visibility = "visible"
)

type DisclosureDecision struct {
	ownerID         domainidentity.ID
	socialContextID string
	viewerID        domainidentity.ID
	visibility      Visibility
}

func NewDisclosureDecision(
	actorID domainidentity.ID,
	ownerID domainidentity.ID,
	socialContextID string,
	viewerID domainidentity.ID,
	visibility Visibility,
) (DisclosureDecision, error) {
	normalizedActorID, err := domainidentity.NewID(string(actorID))
	if err != nil {
		return DisclosureDecision{}, err
	}
	normalizedOwnerID, err := domainidentity.NewID(string(ownerID))
	if err != nil {
		return DisclosureDecision{}, err
	}
	normalizedViewerID, err := domainidentity.NewID(string(viewerID))
	if err != nil {
		return DisclosureDecision{}, err
	}
	if normalizedActorID != normalizedOwnerID {
		return DisclosureDecision{}, ErrOnlyContextOwnerCanChangeDisclosure
	}

	normalizedContextID := strings.TrimSpace(socialContextID)
	if normalizedContextID == "" {
		return DisclosureDecision{}, ErrBlankSocialContextID
	}
	if err := validateVisibility(visibility); err != nil {
		return DisclosureDecision{}, err
	}

	return DisclosureDecision{
		ownerID:         normalizedOwnerID,
		socialContextID: normalizedContextID,
		viewerID:        normalizedViewerID,
		visibility:      visibility,
	}, nil
}

func (d *DisclosureDecision) ChangeVisibility(actorID domainidentity.ID, visibility Visibility) error {
	normalizedActorID, err := domainidentity.NewID(string(actorID))
	if err != nil {
		return err
	}
	if normalizedActorID != d.ownerID {
		return ErrOnlyContextOwnerCanChangeDisclosure
	}
	if err := validateVisibility(visibility); err != nil {
		return err
	}
	d.visibility = visibility
	return nil
}

func (d DisclosureDecision) OwnerID() domainidentity.ID { return d.ownerID }
func (d DisclosureDecision) SocialContextID() string { return d.socialContextID }
func (d DisclosureDecision) ViewerID() domainidentity.ID { return d.viewerID }
func (d DisclosureDecision) Visibility() Visibility { return d.visibility }
func (d DisclosureDecision) AllowsDisclosure() bool { return d.visibility == VisibilityVisible }

func AllowsDisclosure(decision *DisclosureDecision) bool {
	return decision != nil && decision.AllowsDisclosure()
}

func validateVisibility(visibility Visibility) error {
	if visibility != VisibilityHidden && visibility != VisibilityVisible {
		return ErrInvalidVisibility
	}
	return nil
}
