package friendship

import (
	"context"
	"testing"

	domainfriendship "github.com/kinrelay/kin/apps/api/internal/domain/friendship"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func TestMemoryRepositoryFindActiveBetweenProjectsOnlyAcceptedFriendship(t *testing.T) {
	ctx := context.Background()
	userA, _ := domainidentity.NewID("user-a")
	userB, _ := domainidentity.NewID("user-b")
	invitation, err := domainfriendship.Invite(userA, userB)
	if err != nil {
		t.Fatalf("Invite() error = %v", err)
	}
	repository := NewMemoryRepository()
	if created, err := repository.CreateIfAbsent(ctx, invitation); err != nil || !created {
		t.Fatalf("CreateIfAbsent() = (%v, %v), want (true, nil)", created, err)
	}

	if _, found, err := repository.FindActiveBetween(ctx, userA, userB); err != nil || found {
		t.Fatalf("pending FindActiveBetween() = (_, %v, %v), want (_, false, nil)", found, err)
	}

	if _, _, err := repository.UpdateBetween(ctx, userA, userB, func(friendship *domainfriendship.Friendship) error {
		return friendship.Accept(userB)
	}); err != nil {
		t.Fatalf("UpdateBetween() error = %v", err)
	}

	fromA, found, err := repository.FindActiveBetween(ctx, userA, userB)
	if err != nil {
		t.Fatalf("FindActiveBetween() error = %v", err)
	}
	if !found {
		t.Fatal("FindActiveBetween() found = false, want true")
	}
	fromB, found, err := repository.FindActiveBetween(ctx, userB, userA)
	if err != nil || !found {
		t.Fatalf("reverse FindActiveBetween() = (_, %v, %v), want (_, true, nil)", found, err)
	}
	if fromA != fromB {
		t.Fatalf("participant projections differ: fromA=%#v fromB=%#v", fromA, fromB)
	}
	if !fromA.Active || fromA.FirstParticipantID != "user-a" || fromA.SecondParticipantID != "user-b" {
		t.Fatalf("projection = %#v, want active user-a/user-b", fromA)
	}
}
