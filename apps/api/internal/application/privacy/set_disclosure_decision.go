package privacy

import (
	"context"
	"errors"
	"strings"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

var ErrSocialContextNotFound = errors.New("social context not found")

type SocialContextOwnerReader interface {
	OwnerOf(ctx context.Context, socialContextID string) (domainidentity.ID, bool, error)
}

type DisclosureDecisionRepository interface {
	Find(
		ctx context.Context,
		ownerID domainidentity.ID,
		socialContextID string,
		viewerID domainidentity.ID,
	) (domainprivacy.DisclosureDecision, bool, error)
	Save(ctx context.Context, decision domainprivacy.DisclosureDecision) error
}

type SetDisclosureDecisionCommand struct {
	RequesterID     string
	SocialContextID string
	ViewerID        string
	Visibility      domainprivacy.Visibility
}

type SetDisclosureDecision struct {
	ownerReader SocialContextOwnerReader
	repository  DisclosureDecisionRepository
}

func NewSetDisclosureDecision(
	ownerReader SocialContextOwnerReader,
	repository DisclosureDecisionRepository,
) SetDisclosureDecision {
	return SetDisclosureDecision{ownerReader: ownerReader, repository: repository}
}

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
