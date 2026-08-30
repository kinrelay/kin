package socialcontext

import (
	"context"
	"testing"
	"time"

	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func repositoryCandidate(t *testing.T, id string) domainsocialcontext.ContextCandidate {
	t.Helper()
	owner, err := domainidentity.NewID("alice")
	if err != nil {
		t.Fatalf("NewID(): %v", err)
	}
	meaning, err := domainsocialcontext.NewMeaning("最近持續研究 distributed systems")
	if err != nil {
		t.Fatalf("NewMeaning(): %v", err)
	}
	candidate, err := domainsocialcontext.NewContextCandidate(
		id,
		owner,
		meaning,
		[]string{"activity-1", "activity-2"},
		time.Date(2026, time.August, 30, 2, 30, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("NewContextCandidate(): %v", err)
	}
	return candidate
}

func TestMemoryCandidateRepositorySavesOwnerPrivateUnvalidatedCandidate(t *testing.T) {
	repository := NewMemoryCandidateRepository()
	candidate := repositoryCandidate(t, "candidate-1")

	if err := repository.Save(context.Background(), candidate); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	repository.mu.RLock()
	stored, exists := repository.candidates[string(candidate.ID())]
	repository.mu.RUnlock()
	if !exists {
		t.Fatal("candidate was not persisted")
	}
	if stored.ID() != candidate.ID() || stored.OwnerID() != candidate.OwnerID() || stored.Meaning() != candidate.Meaning() {
		t.Fatalf("stored candidate = %#v, want %#v", stored, candidate)
	}
	if !stored.IsPrivate() || stored.IsValidatedSocialContext() {
		t.Fatalf("stored lifecycle invalid: private=%v validated=%v", stored.IsPrivate(), stored.IsValidatedSocialContext())
	}
}

func TestMemoryCandidateRepositorySeparatesCandidatesByStableID(t *testing.T) {
	repository := NewMemoryCandidateRepository()
	for _, candidate := range []domainsocialcontext.ContextCandidate{
		repositoryCandidate(t, "candidate-1"),
		repositoryCandidate(t, "candidate-2"),
	} {
		if err := repository.Save(context.Background(), candidate); err != nil {
			t.Fatalf("Save(%q) error = %v", candidate.ID(), err)
		}
	}

	repository.mu.RLock()
	count := len(repository.candidates)
	repository.mu.RUnlock()
	if count != 2 {
		t.Fatalf("candidate count = %d, want 2", count)
	}
}
