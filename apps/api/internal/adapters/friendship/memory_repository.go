package friendship

import (
	"context"
	"sort"
	"sync"

	domainfriendship "github.com/kinrelay/kin/apps/api/internal/domain/friendship"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

// MemoryRepository stores one Friendship aggregate per unordered participant pair.
type MemoryRepository struct {
	mu          sync.RWMutex
	friendships map[string]domainfriendship.Friendship
}

// NewMemoryRepository constructs an empty in-memory Friendship repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{friendships: make(map[string]domainfriendship.Friendship)}
}

// FindBetween returns the Friendship aggregate for an unordered participant pair.
func (r *MemoryRepository) FindBetween(_ context.Context, first, second domainidentity.ID) (domainfriendship.Friendship, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	found, ok := r.friendships[memoryKey(first, second)]
	return found, ok, nil
}

// CreateIfAbsent atomically creates the participant-pair aggregate when none exists.
func (r *MemoryRepository) CreateIfAbsent(_ context.Context, friendship domainfriendship.Friendship) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := memoryKey(friendship.InviterID(), friendship.InviteeID())
	if _, exists := r.friendships[key]; exists {
		return false, nil
	}
	r.friendships[key] = friendship
	return true, nil
}

// UpdateBetween atomically loads, mutates through domain behavior, and persists one aggregate.
func (r *MemoryRepository) UpdateBetween(_ context.Context, first, second domainidentity.ID, update func(*domainfriendship.Friendship) error) (domainfriendship.Friendship, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := memoryKey(first, second)
	friendship, exists := r.friendships[key]
	if !exists {
		return domainfriendship.Friendship{}, false, nil
	}
	if err := update(&friendship); err != nil {
		return domainfriendship.Friendship{}, true, err
	}
	r.friendships[key] = friendship
	return friendship, true, nil
}

func memoryKey(first, second domainidentity.ID) string {
	parts := []string{string(first), string(second)}
	sort.Strings(parts)
	return parts[0] + "\x00" + parts[1]
}
