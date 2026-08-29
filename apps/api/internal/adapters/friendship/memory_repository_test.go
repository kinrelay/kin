package friendship

import (
	"context"
	"testing"

	domainfriendship "github.com/kinrelay/kin/apps/api/internal/domain/friendship"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func adapterID(t *testing.T, value string) domainidentity.ID {
	t.Helper()
	id, err := domainidentity.NewID(value)
	if err != nil {
		t.Fatalf("NewID(%q): %v", value, err)
	}
	return id
}

func TestMemoryRepositoryFindsParticipantPairRegardlessOfOrderAndUpdatesSingleAggregate(t *testing.T) {
	ctx := context.Background()
	repository := NewMemoryRepository()
	alice := adapterID(t, "alice")
	bob := adapterID(t, "bob")
	invitation, err := domainfriendship.Invite(alice, bob)
	if err != nil {
		t.Fatalf("Invite() error = %v", err)
	}

	if err := repository.Save(ctx, invitation); err != nil {
		t.Fatalf("Save(pending) error = %v", err)
	}
	found, ok, err := repository.FindBetween(ctx, bob, alice)
	if err != nil || !ok {
		t.Fatalf("FindBetween(reverse) = (%v, %v, %v), want friendship", found, ok, err)
	}
	if found.IsActive() {
		t.Fatal("pending friendship must remain pending")
	}

	if err := invitation.Accept(bob); err != nil {
		t.Fatalf("Accept() error = %v", err)
	}
	if err := repository.Save(ctx, invitation); err != nil {
		t.Fatalf("Save(active) error = %v", err)
	}
	found, ok, err = repository.FindBetween(ctx, alice, bob)
	if err != nil || !ok || !found.IsActive() {
		t.Fatalf("FindBetween(active) = (%v, %v, %v), want one active aggregate", found, ok, err)
	}
}
