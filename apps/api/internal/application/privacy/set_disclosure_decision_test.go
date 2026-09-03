package privacy

import (
	"context"
	"errors"
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

type fakeSocialContextOwnerReader struct {
	ownerID domainidentity.ID
	found   bool
	err     error
}

func (f fakeSocialContextOwnerReader) OwnerOf(context.Context, string) (domainidentity.ID, bool, error) {
	return f.ownerID, f.found, f.err
}

type fakeDisclosureDecisionRepository struct {
	decision domainprivacy.DisclosureDecision
	found    bool
	findErr  error
	saveErr  error
	saved    []domainprivacy.DisclosureDecision
}

func (f *fakeDisclosureDecisionRepository) Find(
	context.Context,
	domainidentity.ID,
	string,
	domainidentity.ID,
) (domainprivacy.DisclosureDecision, bool, error) {
	if f.findErr != nil {
		return domainprivacy.DisclosureDecision{}, false, f.findErr
	}
	return f.decision, f.found, nil
}

func (f *fakeDisclosureDecisionRepository) Save(_ context.Context, decision domainprivacy.DisclosureDecision) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.saved = append(f.saved, decision)
	f.decision = decision
	f.found = true
	return nil
}

func appPrivacyTestID(t *testing.T, value string) domainidentity.ID {
	t.Helper()
	id, err := domainidentity.NewID(value)
	if err != nil {
		t.Fatalf("NewID(%q) error = %v", value, err)
	}
	return id
}

func TestSetDisclosureDecisionOwnerCanCreateAndRevoke(t *testing.T) {
	ctx := context.Background()
	ownerID := appPrivacyTestID(t, "owner-1")
	viewerID := appPrivacyTestID(t, "friend-1")
	repository := &fakeDisclosureDecisionRepository{}
	useCase := NewSetDisclosureDecision(
		fakeSocialContextOwnerReader{ownerID: ownerID, found: true},
		repository,
	)

	created, err := useCase.Execute(ctx, SetDisclosureDecisionCommand{
		RequesterID:     string(ownerID),
		SocialContextID: "social-context-1",
		ViewerID:        string(viewerID),
		Visibility:      domainprivacy.VisibilityVisible,
	})
	if err != nil {
		t.Fatalf("Execute(visible) error = %v", err)
	}
	if !created.AllowsDisclosure() {
		t.Fatal("owner-created visible decision should allow disclosure")
	}
	if len(repository.saved) != 1 {
		t.Fatalf("saved decisions = %d, want 1", len(repository.saved))
	}

	revoked, err := useCase.Execute(ctx, SetDisclosureDecisionCommand{
		RequesterID:     string(ownerID),
		SocialContextID: "social-context-1",
		ViewerID:        string(viewerID),
		Visibility:      domainprivacy.VisibilityHidden,
	})
	if err != nil {
		t.Fatalf("Execute(hidden) error = %v", err)
	}
	if revoked.AllowsDisclosure() {
		t.Fatal("owner-hidden decision should revoke disclosure")
	}
	if len(repository.saved) != 2 {
		t.Fatalf("saved decisions after revoke = %d, want 2", len(repository.saved))
	}
}

func TestSetDisclosureDecisionRejectsViewerAndThirdPartyMutation(t *testing.T) {
	ctx := context.Background()
	ownerID := appPrivacyTestID(t, "owner-1")
	viewerID := appPrivacyTestID(t, "friend-1")
	thirdPartyID := appPrivacyTestID(t, "other-1")

	for _, requesterID := range []domainidentity.ID{viewerID, thirdPartyID} {
		repository := &fakeDisclosureDecisionRepository{}
		useCase := NewSetDisclosureDecision(
			fakeSocialContextOwnerReader{ownerID: ownerID, found: true},
			repository,
		)

		_, err := useCase.Execute(ctx, SetDisclosureDecisionCommand{
			RequesterID:     string(requesterID),
			SocialContextID: "social-context-1",
			ViewerID:        string(viewerID),
			Visibility:      domainprivacy.VisibilityVisible,
		})
		if !errors.Is(err, domainprivacy.ErrOnlyContextOwnerCanChangeDisclosure) {
			t.Fatalf("requester %q Execute() error = %v, want %v", requesterID, err, domainprivacy.ErrOnlyContextOwnerCanChangeDisclosure)
		}
		if len(repository.saved) != 0 {
			t.Fatalf("requester %q unexpectedly persisted a decision", requesterID)
		}
	}
}

func TestCheckDisclosurePermissionDefaultsToDeny(t *testing.T) {
	ctx := context.Background()
	ownerID := appPrivacyTestID(t, "owner-1")
	viewerID := appPrivacyTestID(t, "friend-1")
	repository := &fakeDisclosureDecisionRepository{}
	checker := NewCheckDisclosurePermission(repository)

	allowed, err := checker.Execute(ctx, CheckDisclosurePermissionQuery{
		OwnerID:         string(ownerID),
		SocialContextID: "social-context-1",
		ViewerID:        string(viewerID),
	})
	if err != nil {
		t.Fatalf("Execute(absent) error = %v", err)
	}
	if allowed {
		t.Fatal("absent decision must deny disclosure")
	}
}

func TestSetDisclosureDecisionRejectsMissingSocialContext(t *testing.T) {
	ownerID := appPrivacyTestID(t, "owner-1")
	viewerID := appPrivacyTestID(t, "friend-1")
	repository := &fakeDisclosureDecisionRepository{}
	useCase := NewSetDisclosureDecision(
		fakeSocialContextOwnerReader{found: false},
		repository,
	)

	_, err := useCase.Execute(context.Background(), SetDisclosureDecisionCommand{
		RequesterID:     string(ownerID),
		SocialContextID: "missing-context",
		ViewerID:        string(viewerID),
		Visibility:      domainprivacy.VisibilityVisible,
	})
	if !errors.Is(err, ErrSocialContextNotFound) {
		t.Fatalf("Execute(missing context) error = %v, want %v", err, ErrSocialContextNotFound)
	}
	if len(repository.saved) != 0 {
		t.Fatal("missing social context must not persist disclosure decision")
	}
}
