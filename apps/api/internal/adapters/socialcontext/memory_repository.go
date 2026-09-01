package socialcontext

import (
	"context"
	"fmt"
	"sync"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

type storedSocialContext struct {
	id         string
	ownerID    domainidentity.ID
	context    domainsocialcontext.SocialContext
	promotedAt time.Time
}

type currentStateKey struct {
	ownerID    domainidentity.ID
	semanticID domainsocialcontext.SemanticIdentity
}

type storedCurrentState struct {
	occurredAt time.Time
	context    *domainsocialcontext.SocialContext
}

// MemoryRepository is the MVP in-memory SocialContext persistence adapter used for application composition/tests.
type MemoryRepository struct {
	mu            sync.RWMutex
	entries       []storedSocialContext
	currentStates map[currentStateKey]storedCurrentState
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{currentStates: make(map[currentStateKey]storedCurrentState)}
}

func (r *MemoryRepository) SaveIfAbsent(_ context.Context, ownerID domainidentity.ID, socialContext domainsocialcontext.SocialContext) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, entry := range r.entries {
		if entry.ownerID == ownerID && entry.context.Meaning() == socialContext.Meaning() {
			return false, nil
		}
	}
	r.entries = append(r.entries, storedSocialContext{
		id:         fmt.Sprintf("social-context-%d", len(r.entries)+1),
		ownerID:    ownerID,
		context:    socialContext,
		promotedAt: time.Now().UTC(),
	})
	return true, nil
}

func (r *MemoryRepository) ReconcileOwnerCurrentState(
	_ context.Context,
	ownerID domainidentity.ID,
	mutations []domainsocialcontext.CurrentStateMutation,
) (int, error) {
	normalized := make([]domainsocialcontext.CurrentStateMutation, 0, len(mutations))
	for _, mutation := range mutations {
		semanticID, err := domainsocialcontext.NewSemanticIdentity(mutation.SemanticID.String())
		if err != nil {
			return 0, err
		}
		mutation.SemanticID = semanticID
		normalized = append(normalized, mutation)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	next := make(map[currentStateKey]storedCurrentState, len(r.currentStates)+len(normalized))
	for key, state := range r.currentStates {
		next[key] = state
	}

	changed := 0
	for _, mutation := range normalized {
		key := currentStateKey{ownerID: ownerID, semanticID: mutation.SemanticID}
		if current, exists := next[key]; exists && !mutation.OccurredAt.After(current.occurredAt) {
			continue
		}

		var storedReplacement *domainsocialcontext.SocialContext
		if mutation.Replacement != nil {
			copy := *mutation.Replacement
			storedReplacement = &copy
		}
		next[key] = storedCurrentState{occurredAt: mutation.OccurredAt, context: storedReplacement}
		changed++
	}

	r.currentStates = next
	return changed, nil
}

func (r *MemoryRepository) ReconcileCurrentState(
	ctx context.Context,
	ownerID domainidentity.ID,
	semanticID domainsocialcontext.SemanticIdentity,
	occurredAt time.Time,
	replacement *domainsocialcontext.SocialContext,
) (bool, error) {
	changed, err := r.ReconcileOwnerCurrentState(ctx, ownerID, []domainsocialcontext.CurrentStateMutation{{
		SemanticID:  semanticID,
		OccurredAt:  occurredAt,
		Replacement: replacement,
	}})
	if err != nil {
		return false, err
	}
	return changed > 0, nil
}

func (r *MemoryRepository) All() []domainsocialcontext.SocialContext {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domainsocialcontext.SocialContext, 0, len(r.entries)+len(r.currentStates))
	for _, entry := range r.entries {
		result = append(result, entry.context)
	}
	for _, state := range r.currentStates {
		if state.context != nil {
			result = append(result, *state.context)
		}
	}
	return result
}
