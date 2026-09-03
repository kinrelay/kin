package privacy

import (
	"context"
	"errors"
	"strings"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

// ErrSocialContextNotFound indicates the requested Social Context does not exist.
var ErrSocialContextNotFound = errors.New("social context not found")

// SocialContextOwnerReader resolves canonical ownership evidence for a Social Context.
type SocialContextOwnerReader interface {
	OwnerOf(ctx context.Context, socialContextID string) (domainidentity.ID, bool, error)
}

// DisclosureDecisionRepository persists owner-controlled disclosure decisions.
type DisclosureDecisionRepository interface {
	Find(
		ctx context.Context,
		ownerID domainidentity.ID,
		socialContextID string,
		viewerID domainidentity.ID,
	) (domainprivacy.DisclosureDecision, bool, error)
	Save(ctx context.Context, decision domainprivacy.DisclosureDecision) error
}

// SetDisclosureDecisionCommand identifies the requester, context, viewer, and desired visibility.
type SetDisclosureDecisionCommand struct {
	RequesterID     string
	SocialContextID string
	ViewerID        string
	Visibility      domainprivacy.Visibility
}

// SetDisclosureDecision creates or changes one owner-controlled disclosure decision.
type SetDisclosureDecision struct {
	ownerReader SocialContextOwnerReader
	repository  DisclosureDecisionRepository
}

// NewSetDisclosureDecision constructs the write-side disclosure use case.
func NewSetDisclosureDecision(
	ownerReader SocialContextOwnerReader,
	repository DisclosureDecisionRepository,
) SetDisclosureDecision {
	return SetDisclosureDecision{ownerReader: ownerReader, repository: repository}
}

// Execute verifies canonical Social Context ownership before persisting the decision.
func (uc SetDisclosureDecision) Execute(
	ctx context.Context,
	command SetDisclosureDecisionCommand,
) (domainprivacy.DisclosureDecision, error) {
	requesterID, err := domainidentity.NewID(command.RequesterID)
	if err != nil {
		return domainprivacy.DisclosureDecision{}, err
	}
	viewerID, err := domainidentity.NewID(command.ViewerID)
	if err != nil {
		return domainprivacy.DisclosureDecision{}, err
	}
	socialContextID := strings.TrimSpace(command.SocialContextID)
	if socialContextID == "" {
		return domainprivacy.DisclosureDecision{}, domainprivacy.ErrBlankSocialContextID
	}

	ownerID, found, err := uc.ownerReader.OwnerOf(ctx, socialContextID)
	if err != nil {
		return domainprivacy.DisclosureDecision{}, err
	}
	if !found {
		return domainprivacy.DisclosureDecision{}, ErrSocialContextNotFound
	}
	if requesterID != ownerID {
		return domainprivacy.DisclosureDecision{}, domainprivacy.ErrOnlyContextOwnerCanChangeDisclosure
	}

	decision, exists, err := uc.repository.Find(ctx, ownerID, socialContextID, viewerID)
	if err != nil {
		return domainprivacy.DisclosureDecision{}, err
	}
	if exists {
		if err := decision.ChangeVisibility(requesterID, command.Visibility); err != nil {
			return domainprivacy.DisclosureDecision{}, err
		}
	} else {
		decision, err = domainprivacy.NewDisclosureDecision(
			requesterID,
			ownerID,
			socialContextID,
			viewerID,
			command.Visibility,
		)
		if err != nil {
			return domainprivacy.DisclosureDecision{}, err
		}
	}

	if err := uc.repository.Save(ctx, decision); err != nil {
		return domainprivacy.DisclosureDecision{}, err
	}
	return decision, nil
}
