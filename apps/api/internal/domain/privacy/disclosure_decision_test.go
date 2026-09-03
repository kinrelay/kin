package privacy

import (
	"errors"
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func privacyTestID(t *testing.T, value string) domainidentity.ID {
	t.Helper()
	id, err := domainidentity.NewID(value)
	if err != nil {
		t.Fatalf("NewID(%q) error = %v", value, err)
	}
	return id
}

func TestDisclosureDecisionOwnerControlsVisibility(t *testing.T) {
	ownerID := privacyTestID(t, "owner-1")
	viewerID := privacyTestID(t, "friend-1")
	intruderID := privacyTestID(t, "intruder-1")

	decision, err := NewDisclosureDecision(ownerID, ownerID, "social-context-1", viewerID, VisibilityVisible)
	if err != nil {
		t.Fatalf("NewDisclosureDecision() error = %v", err)
	}
	if !decision.AllowsDisclosure() {
		t.Fatal("visible decision should allow disclosure")
	}

	if err := decision.ChangeVisibility(viewerID, VisibilityHidden); !errors.Is(err, ErrOnlyContextOwnerCanChangeDisclosure) {
		t.Fatalf("viewer ChangeVisibility() error = %v, want %v", err, ErrOnlyContextOwnerCanChangeDisclosure)
	}
	if !decision.AllowsDisclosure() {
		t.Fatal("failed viewer mutation must not change the owner's decision")
	}

	if err := decision.ChangeVisibility(intruderID, VisibilityHidden); !errors.Is(err, ErrOnlyContextOwnerCanChangeDisclosure) {
		t.Fatalf("third-party ChangeVisibility() error = %v, want %v", err, ErrOnlyContextOwnerCanChangeDisclosure)
	}

	if err := decision.ChangeVisibility(ownerID, VisibilityHidden); err != nil {
		t.Fatalf("owner ChangeVisibility(hidden) error = %v", err)
	}
	if decision.AllowsDisclosure() {
		t.Fatal("hidden decision should revoke disclosure")
	}
}

func TestDisclosurePermissionDefaultsToDeny(t *testing.T) {
	if AllowsDisclosure(nil) {
		t.Fatal("absence of a disclosure decision must deny disclosure")
	}

	ownerID := privacyTestID(t, "owner-1")
	viewerID := privacyTestID(t, "friend-1")
	decision, err := NewDisclosureDecision(ownerID, ownerID, "social-context-1", viewerID, VisibilityHidden)
	if err != nil {
		t.Fatalf("NewDisclosureDecision(hidden) error = %v", err)
	}
	if AllowsDisclosure(&decision) {
		t.Fatal("explicit hidden decision must deny disclosure")
	}
}

func TestDisclosureDecisionRejectsInvalidCreation(t *testing.T) {
	ownerID := privacyTestID(t, "owner-1")
	viewerID := privacyTestID(t, "friend-1")

	if _, err := NewDisclosureDecision(viewerID, ownerID, "social-context-1", viewerID, VisibilityVisible); !errors.Is(err, ErrOnlyContextOwnerCanChangeDisclosure) {
		t.Fatalf("non-owner NewDisclosureDecision() error = %v, want %v", err, ErrOnlyContextOwnerCanChangeDisclosure)
	}
	if _, err := NewDisclosureDecision(ownerID, ownerID, " ", viewerID, VisibilityVisible); !errors.Is(err, ErrBlankSocialContextID) {
		t.Fatalf("blank context NewDisclosureDecision() error = %v, want %v", err, ErrBlankSocialContextID)
	}
	if _, err := NewDisclosureDecision(ownerID, ownerID, "social-context-1", viewerID, Visibility("public")); !errors.Is(err, ErrInvalidVisibility) {
		t.Fatalf("invalid visibility NewDisclosureDecision() error = %v, want %v", err, ErrInvalidVisibility)
	}
}
