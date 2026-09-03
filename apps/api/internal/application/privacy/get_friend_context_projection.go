package privacy

import (
	"context"
	"errors"
	"strings"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

var ErrFriendContextUnauthorized = errors.New("requester is not an active friendship participant")

// FriendSocialContext is the minimal private context data needed to evaluate friend visibility.
type FriendSocialContext struct {
	ID      string
	OwnerID domainidentity.ID
	Meaning string
}

// ActiveFriendshipReader provides read-only active relationship evidence.
type ActiveFriendshipReader interface {
	IsActiveBetween(ctx context.Context, first, second domainidentity.ID) (bool, error)
}

// FriendSocialContextReader loads the canonical Social Context owner and semantic meaning.
type FriendSocialContextReader interface {
	FindByID(ctx context.Context, id string) (FriendSocialContext, bool, error)
}

// DisclosureDecisionReader reads relationship-specific disclosure state without mutation.
type DisclosureDecisionReader interface {
	Find(
		ctx context.Context,
		ownerID domainidentity.ID,
		socialContextID string,
		viewerID domainidentity.ID,
	) (domainprivacy.DisclosureDecision, bool, error)
}

// GetFriendContextProjectionQuery identifies the authenticated viewer and target owner/context.
// There is deliberately no caller-supplied ViewerID field separate from authentication identity.
type GetFriendContextProjectionQuery struct {
	AuthenticatedViewerID string
	OwnerID               string
	SocialContextID       string
}

// GetFriendContextProjection evaluates privacy before any future relevance/ranking stage.
type GetFriendContextProjection struct {
	friendships ActiveFriendshipReader
	contexts    FriendSocialContextReader
	decisions   DisclosureDecisionReader
}

func NewGetFriendContextProjection(
	friendships ActiveFriendshipReader,
	contexts FriendSocialContextReader,
	decisions DisclosureDecisionReader,
) GetFriendContextProjection {
	return GetFriendContextProjection{friendships: friendships, contexts: contexts, decisions: decisions}
}

func (uc GetFriendContextProjection) Execute(
	ctx context.Context,
	query GetFriendContextProjectionQuery,
) (domainprivacy.ContextProjection, bool, error) {
	viewerID, err := domainidentity.NewID(query.AuthenticatedViewerID)
	if err != nil {
		return domainprivacy.ContextProjection{}, false, err
	}
	ownerID, err := domainidentity.NewID(query.OwnerID)
	if err != nil {
		return domainprivacy.ContextProjection{}, false, err
	}
	contextID := strings.TrimSpace(query.SocialContextID)
	if contextID == "" {
		return domainprivacy.ContextProjection{}, false, domainprivacy.ErrBlankSocialContextID
	}

	socialContext, found, err := uc.contexts.FindByID(ctx, contextID)
	if err != nil {
		return domainprivacy.ContextProjection{}, false, err
	}
	if !found || socialContext.OwnerID != ownerID {
		return domainprivacy.ContextProjection{}, false, ErrSocialContextNotFound
	}

	active, err := uc.friendships.IsActiveBetween(ctx, ownerID, viewerID)
	if err != nil {
		return domainprivacy.ContextProjection{}, false, err
	}
	if !active {
		return domainprivacy.ContextProjection{}, false, ErrFriendContextUnauthorized
	}

	decision, exists, err := uc.decisions.Find(ctx, ownerID, contextID, viewerID)
	if err != nil {
		return domainprivacy.ContextProjection{}, false, err
	}
	if !exists {
		return domainprivacy.ContextProjection{}, false, nil
	}

	projection, visible := domainprivacy.ProjectContext(socialContext.Meaning, &decision)
	return projection, visible, nil
}
