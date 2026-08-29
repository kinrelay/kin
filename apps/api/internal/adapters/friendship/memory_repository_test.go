package friendship

import (
	"context"
	"errors"
	"sync"
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

	created, err := repository.CreateIfAbsent(ctx, invitation)
	if err != nil || !created {
		t.Fatalf("CreateIfAbsent(pending) = (%v, %v), want created", created, err)
	}
	found, ok, err := repository.FindBetween(ctx, bob, alice)
	if err != nil || !ok {
		t.Fatalf("FindBetween(reverse) = (%v, %v, %v), want friendship", found, ok, err)
	}
	if found.IsActive() {
		t.Fatal("pending friendship must remain pending")
	}

	updated, ok, err := repository.UpdateBetween(ctx, alice, bob, func(friendship *domainfriendship.Friendship) error {
		return friendship.Accept(bob)
	})
	if err != nil || !ok || !updated.IsActive() {
		t.Fatalf("UpdateBetween(active) = (%v, %v, %v), want one active aggregate", updated, ok, err)
	}
}

func TestMemoryRepositoryCreateIfAbsentIsAtomicForParticipantPair(t *testing.T) {
	ctx := context.Background()
	repository := NewMemoryRepository()
	alice := adapterID(t, "alice")
	bob := adapterID(t, "bob")
	forward, _ := domainfriendship.Invite(alice, bob)
	reverse, _ := domainfriendship.Invite(bob, alice)

	start := make(chan struct{})
	results := make(chan bool, 2)
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for _, invitation := range []domainfriendship.Friendship{forward, reverse} {
		wg.Add(1)
		go func(invitation domainfriendship.Friendship) {
			defer wg.Done()
			<-start
			created, err := repository.CreateIfAbsent(ctx, invitation)
			results <- created
			errs <- err
		}(invitation)
	}
	close(start)
	wg.Wait()
	close(results)
	close(errs)

	createdCount := 0
	for created := range results {
		if created {
			createdCount++
		}
	}
	for err := range errs {
		if err != nil {
			t.Fatalf("CreateIfAbsent() error = %v", err)
		}
	}
	if createdCount != 1 {
		t.Fatalf("created count = %d, want 1", createdCount)
	}
}

func TestMemoryRepositoryUpdateBetweenSerializesConcurrentAcceptance(t *testing.T) {
	ctx := context.Background()
	repository := NewMemoryRepository()
	alice := adapterID(t, "alice")
	bob := adapterID(t, "bob")
	invitation, _ := domainfriendship.Invite(alice, bob)
	created, err := repository.CreateIfAbsent(ctx, invitation)
	if err != nil || !created {
		t.Fatalf("CreateIfAbsent() = (%v, %v), want created", created, err)
	}

	start := make(chan struct{})
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, _, err := repository.UpdateBetween(ctx, alice, bob, func(friendship *domainfriendship.Friendship) error {
				return friendship.Accept(bob)
			})
			errs <- err
		}()
	}
	close(start)
	wg.Wait()
	close(errs)

	successes := 0
	alreadyActive := 0
	for err := range errs {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, domainfriendship.ErrAlreadyActive):
			alreadyActive++
		default:
			t.Fatalf("UpdateBetween() error = %v", err)
		}
	}
	if successes != 1 || alreadyActive != 1 {
		t.Fatalf("concurrent accepts = %d success / %d already-active, want 1 / 1", successes, alreadyActive)
	}
}
