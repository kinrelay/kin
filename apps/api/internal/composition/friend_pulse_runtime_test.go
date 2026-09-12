package composition

import (
	"context"
	"errors"
	"testing"

	applicationfriendpulse "github.com/kinrelay/kin/apps/api/internal/application/friendpulse"
)

func TestNewInMemoryFriendPulseRuntimeProvidesRunnableDelivery(t *testing.T) {
	runtime := NewInMemoryFriendPulseRuntime()

	_, err := runtime.FriendPulse.Get(context.Background(), "viewer-1", "friend-1")
	if !errors.Is(err, applicationfriendpulse.ErrFriendPulseUnauthorized) {
		t.Fatalf("expected wired runtime to enforce active friendship, got %v", err)
	}
}
