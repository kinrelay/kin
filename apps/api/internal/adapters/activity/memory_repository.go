package activity

import (
	"context"
	"sync"

	domainactivity "github.com/kinrelay/kin/apps/api/internal/domain/activity"
)

// MemoryRepository stores normalized Activities in memory for application tests and MVP composition.
type MemoryRepository struct {
	mu         sync.RWMutex
	activities map[domainactivity.ID]domainactivity.Activity
}

// NewMemoryRepository constructs an empty in-memory Activity repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{activities: make(map[domainactivity.ID]domainactivity.Activity)}
}

// Save stores the latest state for one stable Activity ID.
func (r *MemoryRepository) Save(_ context.Context, value domainactivity.Activity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.activities[value.ID()] = value
	return nil
}
