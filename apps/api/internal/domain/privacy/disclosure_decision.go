package privacy

import (
	"errors"
	"strings"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

var (
	// ErrOnlyContextOwnerCanChangeDisclosure indicates a non-owner attempted to mutate disclosure.
	ErrOnlyContextOwnerCanChangeDisclosure = errors.New("only the context owner can change disclosure")
	// ErrBlankSocialContextID indicates the decision did not identify a Social Context.
	ErrBlankSocialContextID = errors.New("social context id is blank")
	// ErrInvalidVisibility indicates an unsupported disclosure visibility.
	ErrInvalidVisibility = errors.New("disclosure visibility is invalid")
)

// Visibility is the owner-selected disclosure state for one specific viewer.
type Visibility string

const (
	// VisibilityHidden denies disclosure to the specific viewer.
	VisibilityHidden Visibility = "hidden"
	// VisibilityVisible allows disclosure to the specific viewer.
	VisibilityVisible Visibility = "visible"
)

// DisclosureDecision records an owner's disclosure choice for one Social Context and viewer.
type DisclosureDecision struct {
	ownerID         domainidentity.ID
	socialContextID string
	viewerID        domainidentity.ID
	visibility      Visibility
}

// NewDisclosureDecision creates an owner-authorized relationship-specific disclosure decision.
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

// ChangeVisibility applies a visibility change only when initiated by the Context Owner.
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

// OwnerID returns the Social Context owner who controls this decision.
func (d DisclosureDecision) OwnerID() domainidentity.ID { return d.ownerID }

// SocialContextID returns the Social Context governed by this decision.
func (d DisclosureDecision) SocialContextID() string { return d.socialContextID }

// ViewerID returns the specific viewer governed by this decision.
func (d DisclosureDecision) ViewerID() domainidentity.ID { return d.viewerID }

// Visibility returns the current owner-selected visibility.
func (d DisclosureDecision) Visibility() Visibility { return d.visibility }

// AllowsDisclosure reports whether this decision currently permits disclosure.
func (d DisclosureDecision) AllowsDisclosure() bool { return d.visibility == VisibilityVisible }

// AllowsDisclosure applies default-deny semantics when no decision exists.
func AllowsDisclosure(decision *DisclosureDecision) bool {
	return decision != nil && decision.AllowsDisclosure()
}

func validateVisibility(visibility Visibility) error {
	if visibility != VisibilityHidden && visibility != VisibilityVisible {
		return ErrInvalidVisibility
	}
	return nil
}
