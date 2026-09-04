package privacy

import (
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func TestPendingDeliveryIntentBindsProjectionRevisions(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	viewerID, _ := domainidentity.NewID("friend-1")
	projection := ContextProjection{Meaning: "最近開始深入研究分散式系統"}

	intent, err := NewPendingDeliveryIntent(
		ownerID,
		viewerID,
		"context-1",
		projection,
		PrivacyPolicyRevision(7),
		RelationshipRevision(11),
	)
	if err != nil {
		t.Fatalf("create pending delivery intent: %v", err)
	}

	if intent.State() != DeliveryIntentPending {
		t.Fatalf("expected pending state, got %q", intent.State())
	}
	if intent.OwnerID() != ownerID || intent.ViewerID() != viewerID {
		t.Fatalf("expected owner/viewer binding, got owner=%q viewer=%q", intent.OwnerID(), intent.ViewerID())
	}
	if intent.SocialContextID() != "context-1" {
		t.Fatalf("expected context binding, got %q", intent.SocialContextID())
	}
	if intent.Projection() != projection {
		t.Fatalf("expected projection snapshot, got %+v", intent.Projection())
	}
	if intent.PrivacyRevision() != PrivacyPolicyRevision(7) {
		t.Fatalf("expected privacy revision 7, got %d", intent.PrivacyRevision())
	}
	if intent.RelationshipRevision() != RelationshipRevision(11) {
		t.Fatalf("expected relationship revision 11, got %d", intent.RelationshipRevision())
	}
}
