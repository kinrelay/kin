package activity

import (
	"context"
	"errors"
	"sync"

	domainactivity "github.com/kinrelay/kin/apps/api/internal/domain/activity"
)

// ErrActivityIDConflict indicates that an Activity ID is already bound to different content/state.
var ErrActivityIDConflict = errors.New("activity id already exists with different activity")

// MemoryRepository stores normalized Activities in memory for application tests and MVP composition.
type MemoryRepository struct {
	mu         sync.RWMutex
	activities map[domainactivity.ID]domainactivity.Activity
}

// NewMemoryRepository constructs an empty in-memory Activity repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{activities: make(map[domainactivity.ID]domainactivity.Activity)}
}

// Save creates an Activity by stable ID and treats an identical repeat as an idempotent retry.
func (r *MemoryRepository) Save(_ context.Context, value domainactivity.Activity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, found := r.activities[value.ID()]; found {
		if existing == value {
			return nil
		}
		return ErrActivityIDConflict
	}

	r.activities[value.ID()] = value
	return nil
}
