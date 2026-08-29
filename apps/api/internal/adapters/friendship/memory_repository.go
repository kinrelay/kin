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

// Save creates or replaces the single aggregate for its participant pair.
func (r *MemoryRepository) Save(_ context.Context, friendship domainfriendship.Friendship) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.friendships[memoryKey(friendship.InviterID(), friendship.InviteeID())] = friendship
	return nil
}

func memoryKey(first, second domainidentity.ID) string {
	parts := []string{string(first), string(second)}
	sort.Strings(parts)
	return parts[0] + "\x00" + parts[1]
}
