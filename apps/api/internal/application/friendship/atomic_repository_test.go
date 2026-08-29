package friendship

import (
	"context"

	domainfriendship "github.com/kinrelay/kin/apps/api/internal/domain/friendship"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func (f *friendshipRepositoryFake) CreateIfAbsent(_ context.Context, friendship domainfriendship.Friendship) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	key := friendshipKey(friendship.InviterID(), friendship.InviteeID())
	if _, exists := f.friendships[key]; exists {
		return false, nil
	}
	f.friendships[key] = friendship
	return true, nil
}

func (f *friendshipRepositoryFake) UpdateBetween(_ context.Context, first, second domainidentity.ID, update func(*domainfriendship.Friendship) error) (domainfriendship.Friendship, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	key := friendshipKey(first, second)
	friendship, exists := f.friendships[key]
	if !exists {
		return domainfriendship.Friendship{}, false, nil
	}
	if err := update(&friendship); err != nil {
		return domainfriendship.Friendship{}, true, err
	}
	f.friendships[key] = friendship
	return friendship, true, nil
}
