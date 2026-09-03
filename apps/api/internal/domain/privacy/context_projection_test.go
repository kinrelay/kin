package privacy

import (
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func TestProjectContextDefaultsToDenyWithoutPermission(t *testing.T) {
	projection, visible := ProjectContext("最近開始深入研究分散式系統", nil)
	if visible {
		t.Fatal("expected context to remain hidden without an explicit disclosure decision")
	}
	if projection != (ContextProjection{}) {
		t.Fatalf("expected no projection when permission is absent, got %+v", projection)
	}
}

func TestProjectContextReturnsOnlySemanticMeaningForVisibleDecision(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	viewerID, _ := domainidentity.NewID("friend-1")
	decision, err := NewDisclosureDecision(ownerID, ownerID, "context-1", viewerID, VisibilityVisible)
	if err != nil {
		t.Fatalf("create visible decision: %v", err)
	}

	projection, visible := ProjectContext("最近開始深入研究分散式系統", &decision)
	if !visible {
		t.Fatal("expected explicit visible decision to produce a projection")
	}
	if projection.Meaning != "最近開始深入研究分散式系統" {
		t.Fatalf("expected semantic meaning only, got %+v", projection)
	}
}

func TestProjectContextHidesAfterDecisionIsRevoked(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	viewerID, _ := domainidentity.NewID("friend-1")
	decision, err := NewDisclosureDecision(ownerID, ownerID, "context-1", viewerID, VisibilityVisible)
	if err != nil {
		t.Fatalf("create visible decision: %v", err)
	}
	if err := decision.ChangeVisibility(ownerID, VisibilityHidden); err != nil {
		t.Fatalf("hide decision: %v", err)
	}

	projection, visible := ProjectContext("最近開始深入研究分散式系統", &decision)
	if visible {
		t.Fatal("expected hidden/revoked decision to suppress projection")
	}
	if projection != (ContextProjection{}) {
		t.Fatalf("expected no projection after revoke, got %+v", projection)
	}
}
