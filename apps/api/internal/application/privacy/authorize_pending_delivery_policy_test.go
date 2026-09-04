package privacy

import (
	"context"
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

func TestAuthorizePendingDeliveryReprojectsAfterPrivacyDetailDowngrade(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	viewerID, _ := domainidentity.NewID("friend-1")
	visibleDecision, _ := domainprivacy.NewDisclosureDecision(ownerID, ownerID, "context-1", viewerID, domainprivacy.VisibilityVisible)
	intent, _ := domainprivacy.NewPendingDeliveryIntent(
		ownerID,
		viewerID,
		"context-1",
		domainprivacy.ContextProjection{Meaning: "full private detail including exact location"},
		domainprivacy.PrivacyPolicyRevision(1),
		domainprivacy.RelationshipRevision(1),
	)

	state := &fakeDeliveryAuthorizationState{snapshot: DeliveryAuthorizationSnapshot{
		Meaning:              "full private detail including exact location",
		Decision:             &visibleDecision,
		PrivacyRevision:      domainprivacy.PrivacyPolicyRevision(1),
		RelationshipActive:   true,
		RelationshipRevision: domainprivacy.RelationshipRevision(1),
	}}
	guard := &fakeAtomicRevisionGuard{state: state}
	guard.beforeCheck = func() {
		// Simulate a policy revision whose current privacy projection exposes less detail.
		state.snapshot.Meaning = "coarse location only"
		state.snapshot.PrivacyRevision = domainprivacy.PrivacyPolicyRevision(2)
	}

	uc := NewAuthorizePendingDelivery(fakeDeliveryAuthorizationReader{state: state}, guard)
	result, err := uc.Execute(context.Background(), intent)
	if err != nil {
		t.Fatalf("authorize pending delivery: %v", err)
	}
	if result.State() != domainprivacy.DeliveryIntentDispatchable {
		t.Fatalf("expected reauthorized intent to become dispatchable, got %q", result.State())
	}
	if result.Projection().Meaning != "coarse location only" {
		t.Fatalf("expected downgraded projection only, got %+v", result.Projection())
	}
	if result.PrivacyRevision() != domainprivacy.PrivacyPolicyRevision(2) {
		t.Fatalf("expected latest privacy revision 2, got %d", result.PrivacyRevision())
	}
	if guard.commits != 1 {
		t.Fatalf("expected only the refreshed projection to commit, got %d commits", guard.commits)
	}
}

func TestAuthorizePendingDeliveryDefaultsToDenyWithoutExplicitPermission(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	viewerID, _ := domainidentity.NewID("friend-1")
	intent, _ := domainprivacy.NewPendingDeliveryIntent(
		ownerID,
		viewerID,
		"context-1",
		domainprivacy.ContextProjection{Meaning: "stale projection"},
		domainprivacy.PrivacyPolicyRevision(1),
		domainprivacy.RelationshipRevision(1),
	)

	state := &fakeDeliveryAuthorizationState{snapshot: DeliveryAuthorizationSnapshot{
		Meaning:              "current meaning",
		Decision:             nil,
		PrivacyRevision:      domainprivacy.PrivacyPolicyRevision(2),
		RelationshipActive:   true,
		RelationshipRevision: domainprivacy.RelationshipRevision(1),
	}}
	guard := &fakeAtomicRevisionGuard{state: state}

	uc := NewAuthorizePendingDelivery(fakeDeliveryAuthorizationReader{state: state}, guard)
	result, err := uc.Execute(context.Background(), intent)
	if err != nil {
		t.Fatalf("authorize pending delivery: %v", err)
	}
	if result.State() != domainprivacy.DeliveryIntentCancelled {
		t.Fatalf("expected default-deny cancellation, got %q", result.State())
	}
	if guard.commits != 0 {
		t.Fatalf("missing permission must never commit, got %d commits", guard.commits)
	}
}
