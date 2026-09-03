package privacy

import (
	"context"
	"errors"
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

type fakeActiveFriendshipReader struct {
	active map[string]bool
}

func (f fakeActiveFriendshipReader) IsActiveBetween(_ context.Context, first, second domainidentity.ID) (bool, error) {
	return f.active[string(first)+"|"+string(second)] || f.active[string(second)+"|"+string(first)], nil
}

type fakeFriendSocialContextReader struct {
	items map[string]FriendSocialContext
}

func (f fakeFriendSocialContextReader) FindByID(_ context.Context, id string) (FriendSocialContext, bool, error) {
	item, ok := f.items[id]
	return item, ok, nil
}

type fakeDisclosureDecisionReader struct {
	decisions map[string]domainprivacy.DisclosureDecision
}

func (f fakeDisclosureDecisionReader) Find(_ context.Context, ownerID domainidentity.ID, contextID string, viewerID domainidentity.ID) (domainprivacy.DisclosureDecision, bool, error) {
	decision, ok := f.decisions[string(ownerID)+"|"+contextID+"|"+string(viewerID)]
	return decision, ok, nil
}

func TestGetFriendContextProjectionUsesAuthenticatedCallerAndActiveRelationship(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	viewerID, _ := domainidentity.NewID("friend-1")
	decision, _ := domainprivacy.NewDisclosureDecision(ownerID, ownerID, "context-1", viewerID, domainprivacy.VisibilityVisible)

	uc := NewGetFriendContextProjection(
		fakeActiveFriendshipReader{active: map[string]bool{"owner-1|friend-1": true}},
		fakeFriendSocialContextReader{items: map[string]FriendSocialContext{"context-1": {ID: "context-1", OwnerID: ownerID, Meaning: "最近開始深入研究分散式系統"}}},
		fakeDisclosureDecisionReader{decisions: map[string]domainprivacy.DisclosureDecision{"owner-1|context-1|friend-1": decision}},
	)

	projection, visible, err := uc.Execute(context.Background(), GetFriendContextProjectionQuery{
		AuthenticatedViewerID: "friend-1",
		OwnerID:               "owner-1",
		SocialContextID:       "context-1",
	})
	if err != nil {
		t.Fatalf("execute projection query: %v", err)
	}
	if !visible || projection.Meaning != "最近開始深入研究分散式系統" {
		t.Fatalf("expected semantic projection for active authorized friend, got visible=%v projection=%+v", visible, projection)
	}
}

func TestGetFriendContextProjectionRejectsNonParticipant(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	uc := NewGetFriendContextProjection(
		fakeActiveFriendshipReader{active: map[string]bool{}},
		fakeFriendSocialContextReader{items: map[string]FriendSocialContext{"context-1": {ID: "context-1", OwnerID: ownerID, Meaning: "meaning"}}},
		fakeDisclosureDecisionReader{},
	)

	_, _, err := uc.Execute(context.Background(), GetFriendContextProjectionQuery{
		AuthenticatedViewerID: "stranger-1",
		OwnerID:               "owner-1",
		SocialContextID:       "context-1",
	})
	if !errors.Is(err, ErrFriendContextUnauthorized) {
		t.Fatalf("expected non-participant rejection, got %v", err)
	}
}

func TestGetFriendContextProjectionVariesByAuthenticatedViewer(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	visibleViewer, _ := domainidentity.NewID("friend-visible")
	hiddenViewer, _ := domainidentity.NewID("friend-hidden")
	visibleDecision, _ := domainprivacy.NewDisclosureDecision(ownerID, ownerID, "context-1", visibleViewer, domainprivacy.VisibilityVisible)
	hiddenDecision, _ := domainprivacy.NewDisclosureDecision(ownerID, ownerID, "context-1", hiddenViewer, domainprivacy.VisibilityHidden)

	uc := NewGetFriendContextProjection(
		fakeActiveFriendshipReader{active: map[string]bool{"owner-1|friend-visible": true, "owner-1|friend-hidden": true}},
		fakeFriendSocialContextReader{items: map[string]FriendSocialContext{"context-1": {ID: "context-1", OwnerID: ownerID, Meaning: "meaning"}}},
		fakeDisclosureDecisionReader{decisions: map[string]domainprivacy.DisclosureDecision{
			"owner-1|context-1|friend-visible": visibleDecision,
			"owner-1|context-1|friend-hidden":  hiddenDecision,
		}},
	)

	_, visible, err := uc.Execute(context.Background(), GetFriendContextProjectionQuery{AuthenticatedViewerID: "friend-visible", OwnerID: "owner-1", SocialContextID: "context-1"})
	if err != nil || !visible {
		t.Fatalf("expected visible viewer to receive projection, visible=%v err=%v", visible, err)
	}
	_, visible, err = uc.Execute(context.Background(), GetFriendContextProjectionQuery{AuthenticatedViewerID: "friend-hidden", OwnerID: "owner-1", SocialContextID: "context-1"})
	if err != nil || visible {
		t.Fatalf("expected hidden viewer not to receive projection, visible=%v err=%v", visible, err)
	}
}
