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

// MemoryRepository is the MVP in-memory SocialContext persistence adapter used for application composition/tests.
type MemoryRepository struct {
	mu      sync.RWMutex
	entries []storedSocialContext
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
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

func (r *MemoryRepository) RetireByProvenance(_ context.Context, ownerID domainidentity.ID, activityIDs []string) (int, error) {
	if len(activityIDs) == 0 {
		return 0, nil
	}
	retiredIDs := make(map[string]struct{}, len(activityIDs))
	for _, activityID := range activityIDs {
		if activityID != "" {
			retiredIDs[activityID] = struct{}{}
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	kept := r.entries[:0]
	retired := 0
	for _, entry := range r.entries {
		if entry.ownerID != ownerID || !provenanceIntersects(entry.context.Provenance(), retiredIDs) {
			kept = append(kept, entry)
			continue
		}
		retired++
	}
	r.entries = kept
	return retired, nil
}

func provenanceIntersects(provenance []string, retiredIDs map[string]struct{}) bool {
	for _, activityID := range provenance {
		if _, retired := retiredIDs[activityID]; retired {
			return true
		}
	}
	return false
}

func (r *MemoryRepository) All() []domainsocialcontext.SocialContext {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domainsocialcontext.SocialContext, 0, len(r.entries))
	for _, entry := range r.entries {
		result = append(result, entry.context)
	}
	return result
}
