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

func (r *MemoryRepository) ReconcileCurrentState(
	_ context.Context,
	ownerID domainidentity.ID,
	semanticID domainsocialcontext.SemanticIdentity,
	occurredAt time.Time,
	replacement *domainsocialcontext.SocialContext,
) (bool, error) {
	normalizedSemanticID, err := domainsocialcontext.NewSemanticIdentity(semanticID.String())
	if err != nil {
		return false, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	key := currentStateKey{ownerID: ownerID, semanticID: normalizedSemanticID}
	if current, exists := r.currentStates[key]; exists && !occurredAt.After(current.occurredAt) {
		return false, nil
	}

	var storedReplacement *domainsocialcontext.SocialContext
	if replacement != nil {
		copy := *replacement
		storedReplacement = &copy
	}
	r.currentStates[key] = storedCurrentState{occurredAt: occurredAt, context: storedReplacement}
	return true, nil
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
