package socialcontext

import (
	"context"
	"errors"
	"testing"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

var errReconciliationPersistence = errors.New("reconciliation persistence failed")

type atomicReconciliationRepositoryFake struct {
	attempts int
	applied  []domainsocialcontext.CurrentStateMutation
	fail     error
}

func (f *atomicReconciliationRepositoryFake) ReconcileOwnerCurrentState(_ context.Context, _ domainidentity.ID, mutations []domainsocialcontext.CurrentStateMutation) (int, error) {
	f.attempts++
	if f.fail != nil {
		return 0, f.fail
	}
	f.applied = append([]domainsocialcontext.CurrentStateMutation(nil), mutations...)
	return len(mutations), nil
}

func TestReconcileCurrentSocialContextDelegatesOneAtomicOwnerBatch(t *testing.T) {
	repository := &atomicReconciliationRepositoryFake{}
	useCase := NewReconcileCurrentSocialContext(repository)
	occurredAt := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)

	changed, err := useCase.Execute(context.Background(), ReconcileCurrentSocialContextCommand{
		OwnerID: "owner-1",
		Mutations: []domainsocialcontext.CurrentStateMutation{
			{SemanticID: "marathon", OccurredAt: occurredAt},
			{SemanticID: "distributed-systems", OccurredAt: occurredAt},
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if changed != 2 {
		t.Fatalf("Execute() changed = %d, want 2", changed)
	}
	if repository.attempts != 1 {
		t.Fatalf("repository attempts = %d, want exactly one atomic call", repository.attempts)
	}
	if len(repository.applied) != 2 {
		t.Fatalf("repository applied = %#v, want complete owner batch", repository.applied)
	}
}

func TestReconcileCurrentSocialContextDoesNotSplitRetireAndReplacementWhenPersistenceFails(t *testing.T) {
	repository := &atomicReconciliationRepositoryFake{fail: errReconciliationPersistence}
	useCase := NewReconcileCurrentSocialContext(repository)
	occurredAt := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)

	_, err := useCase.Execute(context.Background(), ReconcileCurrentSocialContextCommand{
		OwnerID: "owner-1",
		Mutations: []domainsocialcontext.CurrentStateMutation{
			{SemanticID: "marathon", OccurredAt: occurredAt, Replacement: nil},
			{SemanticID: "distributed-systems", OccurredAt: occurredAt},
		},
	})
	if !errors.Is(err, errReconciliationPersistence) {
		t.Fatalf("Execute() error = %v, want persistence error", err)
	}
	if repository.attempts != 1 {
		t.Fatalf("repository attempts = %d, want exactly one atomic call", repository.attempts)
	}
	if len(repository.applied) != 0 {
		t.Fatalf("repository applied = %#v, want no partial state after failed atomic persistence", repository.applied)
	}
}

func TestReconcileCurrentSocialContextRejectsInvalidOwnerBeforePersistence(t *testing.T) {
	repository := &atomicReconciliationRepositoryFake{}
	useCase := NewReconcileCurrentSocialContext(repository)

	if _, err := useCase.Execute(context.Background(), ReconcileCurrentSocialContextCommand{OwnerID: ""}); err == nil {
		t.Fatal("Execute() error = nil, want invalid owner")
	}
	if repository.attempts != 0 {
		t.Fatalf("repository attempts = %d, want no persistence for invalid owner", repository.attempts)
	}
}
