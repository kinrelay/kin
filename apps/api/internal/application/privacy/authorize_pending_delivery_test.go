package privacy

import (
	"context"
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

type fakeDeliveryAuthorizationState struct {
	snapshot DeliveryAuthorizationSnapshot
}

type fakeDeliveryAuthorizationReader struct {
	state *fakeDeliveryAuthorizationState
}

func (f fakeDeliveryAuthorizationReader) Read(
	_ context.Context,
	_, _ domainidentity.ID,
	_ string,
) (DeliveryAuthorizationSnapshot, error) {
	return f.state.snapshot, nil
}

type fakeAtomicRevisionGuard struct {
	state       *fakeDeliveryAuthorizationState
	beforeCheck func()
	commits     int
}

func (f *fakeAtomicRevisionGuard) WithMatchingRevisions(
	_ context.Context,
	_ domainidentity.ID,
	_ string,
	_ domainidentity.ID,
	privacyRevision domainprivacy.PrivacyPolicyRevision,
	relationshipRevision domainprivacy.RelationshipRevision,
	commit func() error,
) (bool, error) {
	if f.beforeCheck != nil {
		hook := f.beforeCheck
		f.beforeCheck = nil
		hook()
	}
	if f.state.snapshot.PrivacyRevision != privacyRevision || f.state.snapshot.RelationshipRevision != relationshipRevision {
		return false, nil
	}
	if err := commit(); err != nil {
		return false, err
	}
	f.commits++
	return true, nil
}

func TestAuthorizePendingDeliveryCancelsWhenPrivacyChangesDuringCommitRace(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	viewerID, _ := domainidentity.NewID("friend-1")
	visibleDecision, _ := domainprivacy.NewDisclosureDecision(ownerID, ownerID, "context-1", viewerID, domainprivacy.VisibilityVisible)
	hiddenDecision, _ := domainprivacy.NewDisclosureDecision(ownerID, ownerID, "context-1", viewerID, domainprivacy.VisibilityHidden)
	intent, _ := domainprivacy.NewPendingDeliveryIntent(
		ownerID,
		viewerID,
		"context-1",
		domainprivacy.ContextProjection{Meaning: "stale private detail"},
		domainprivacy.PrivacyPolicyRevision(1),
		domainprivacy.RelationshipRevision(1),
	)

	state := &fakeDeliveryAuthorizationState{snapshot: DeliveryAuthorizationSnapshot{
		Meaning:              "current private detail",
		Decision:             &visibleDecision,
		PrivacyRevision:      domainprivacy.PrivacyPolicyRevision(1),
		RelationshipActive:   true,
		RelationshipRevision: domainprivacy.RelationshipRevision(1),
	}}
	guard := &fakeAtomicRevisionGuard{state: state}
	guard.beforeCheck = func() {
		state.snapshot.Decision = &hiddenDecision
		state.snapshot.PrivacyRevision = domainprivacy.PrivacyPolicyRevision(2)
	}

	uc := NewAuthorizePendingDelivery(fakeDeliveryAuthorizationReader{state: state}, guard)
	result, err := uc.Execute(context.Background(), intent)
	if err != nil {
		t.Fatalf("authorize pending delivery: %v", err)
	}
	if result.State() != domainprivacy.DeliveryIntentCancelled {
		t.Fatalf("expected cancelled intent after privacy revoke race, got %q", result.State())
	}
	if guard.commits != 0 {
		t.Fatalf("stale payload must never commit, got %d commits", guard.commits)
	}
}

func TestAuthorizePendingDeliveryCancelsWhenRelationshipChangesDuringCommitRace(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	viewerID, _ := domainidentity.NewID("friend-1")
	visibleDecision, _ := domainprivacy.NewDisclosureDecision(ownerID, ownerID, "context-1", viewerID, domainprivacy.VisibilityVisible)
	intent, _ := domainprivacy.NewPendingDeliveryIntent(
		ownerID,
		viewerID,
		"context-1",
		domainprivacy.ContextProjection{Meaning: "stale private detail"},
		domainprivacy.PrivacyPolicyRevision(1),
		domainprivacy.RelationshipRevision(3),
	)

	state := &fakeDeliveryAuthorizationState{snapshot: DeliveryAuthorizationSnapshot{
		Meaning:              "current private detail",
		Decision:             &visibleDecision,
		PrivacyRevision:      domainprivacy.PrivacyPolicyRevision(1),
		RelationshipActive:   true,
		RelationshipRevision: domainprivacy.RelationshipRevision(3),
	}}
	guard := &fakeAtomicRevisionGuard{state: state}
	guard.beforeCheck = func() {
		state.snapshot.RelationshipActive = false
		state.snapshot.RelationshipRevision = domainprivacy.RelationshipRevision(4)
	}

	uc := NewAuthorizePendingDelivery(fakeDeliveryAuthorizationReader{state: state}, guard)
	result, err := uc.Execute(context.Background(), intent)
	if err != nil {
		t.Fatalf("authorize pending delivery: %v", err)
	}
	if result.State() != domainprivacy.DeliveryIntentCancelled {
		t.Fatalf("expected cancelled intent after relationship race, got %q", result.State())
	}
	if guard.commits != 0 {
		t.Fatalf("stale relationship must never commit, got %d commits", guard.commits)
	}
}

func TestAuthorizePendingDeliveryReprojectsAndCommitsInsideRevisionGuard(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	viewerID, _ := domainidentity.NewID("friend-1")
	visibleDecision, _ := domainprivacy.NewDisclosureDecision(ownerID, ownerID, "context-1", viewerID, domainprivacy.VisibilityVisible)
	intent, _ := domainprivacy.NewPendingDeliveryIntent(
		ownerID,
		viewerID,
		"context-1",
		domainprivacy.ContextProjection{Meaning: "old projection"},
		domainprivacy.PrivacyPolicyRevision(1),
		domainprivacy.RelationshipRevision(1),
	)

	state := &fakeDeliveryAuthorizationState{snapshot: DeliveryAuthorizationSnapshot{
		Meaning:              "latest semantic projection",
		Decision:             &visibleDecision,
		PrivacyRevision:      domainprivacy.PrivacyPolicyRevision(2),
		RelationshipActive:   true,
		RelationshipRevision: domainprivacy.RelationshipRevision(5),
	}}
	guard := &fakeAtomicRevisionGuard{state: state}

	uc := NewAuthorizePendingDelivery(fakeDeliveryAuthorizationReader{state: state}, guard)
	result, err := uc.Execute(context.Background(), intent)
	if err != nil {
		t.Fatalf("authorize pending delivery: %v", err)
	}
	if result.State() != domainprivacy.DeliveryIntentDispatchable {
		t.Fatalf("expected dispatchable intent, got %q", result.State())
	}
	if result.Projection().Meaning != "latest semantic projection" {
		t.Fatalf("expected latest projection, got %+v", result.Projection())
	}
	if result.PrivacyRevision() != domainprivacy.PrivacyPolicyRevision(2) || result.RelationshipRevision() != domainprivacy.RelationshipRevision(5) {
		t.Fatalf("expected latest revisions, got privacy=%d relationship=%d", result.PrivacyRevision(), result.RelationshipRevision())
	}
	if guard.commits != 1 {
		t.Fatalf("expected exactly one guarded commit, got %d", guard.commits)
	}
}
