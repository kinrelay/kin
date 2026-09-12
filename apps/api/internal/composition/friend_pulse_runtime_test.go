package composition

import "testing"

func TestNewInMemoryFriendPulseRuntimeProvidesRunnableDelivery(t *testing.T) {
	runtime := NewInMemoryFriendPulseRuntime()
	if runtime.FriendPulse == (FriendPulseFlow{}) {
		t.Fatal("runtime must expose a wired Friend Pulse delivery flow")
	}
}
