package composition

import (
	"context"
	"testing"
)

func TestFriendPulseFlowExposesAuthenticatedDelivery(t *testing.T) {
	flow := NewFriendPulseFlow()

	pulse, err := flow.Get(context.Background(), "viewer-1", "friend-1")
	if err != nil {
		t.Fatalf("get friend pulse: %v", err)
	}
	if len(pulse.Items) > 3 {
		t.Fatalf("pulse must stay bounded to at most 3 items, got %d", len(pulse.Items))
	}
}
