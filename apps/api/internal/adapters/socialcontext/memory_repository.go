package socialcontext

import (
	"context"
	"sync"

	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

// MemoryRepository is the MVP in-memory SocialContext persistence adapter used for application composition/tests.
type MemoryRepository struct {
	mu       sync.RWMutex
	contexts []domainsocialcontext.SocialContext
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (r *MemoryRepository) Save(_ context.Context, socialContext domainsocialcontext.SocialContext) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.contexts = append(r.contexts, socialContext)
	return nil
}

func (r *MemoryRepository) All() []domainsocialcontext.SocialContext {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]domainsocialcontext.SocialContext(nil), r.contexts...)
}
