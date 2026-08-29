package friendship

import (
	"context"
	"errors"
	"sort"
	"testing"

	domainfriendship "github.com/kinrelay/kin/apps/api/internal/domain/friendship"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

type identityDirectoryFake struct {
	identities map[domainidentity.ID]bool
}

func (f identityDirectoryFake) Exists(_ context.Context, id domainidentity.ID) (bool, error) {
	return f.identities[id], nil
}

type friendshipRepositoryFake struct {
	friendships map[string]domainfriendship.Friendship
}

func newFriendshipRepositoryFake() *friendshipRepositoryFake {
	return &friendshipRepositoryFake{friendships: make(map[string]domainfriendship.Friendship)}
}

func friendshipKey(first, second domainidentity.ID) string {
	parts := []string{string(first), string(second)}
	sort.Strings(parts)
	return parts[0] + "\x00" + parts[1]
}

func (f *friendshipRepositoryFake) FindBetween(_ context.Context, first, second domainidentity.ID) (domainfriendship.Friendship, bool, error) {
	found, ok := f.friendships[friendshipKey(first, second)]
	return found, ok, nil
}

func (f *friendshipRepositoryFake) Save(_ context.Context, friendship domainfriendship.Friendship) error {
	f.friendships[friendshipKey(friendship.InviterID(), friendship.InviteeID())] = friendship
	return nil
}

func appID(t *testing.T, value string) domainidentity.ID {
	t.Helper()
	id, err := domainidentity.NewID(value)
	if err != nil {
		t.Fatalf("NewID(%q): %v", value, err)
	}
	return id
}

func testDirectory(t *testing.T, values ...string) identityDirectoryFake {
	t.Helper()
	identities := make(map[domainidentity.ID]bool, len(values))
	for _, value := range values {
		identities[appID(t, value)] = true
	}
	return identityDirectoryFake{identities: identities}
}

func TestInviteFriendValidatesIdentityParticipantsAndPersistsPendingInvitation(t *testing.T) {
	ctx := context.Background()
	repository := newFriendshipRepositoryFake()
	useCase := NewInviteFriend(testDirectory(t, "alice", "bob"), repository)

	created, err := useCase.Execute(ctx, InviteFriendCommand{InviterID: "alice", InviteeID: "bob"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if created.IsActive() {
		t.Fatal("invitation must remain pending until acceptance")
	}
	persisted, ok, err := repository.FindBetween(ctx, appID(t, "alice"), appID(t, "bob"))
	if err != nil || !ok {
		t.Fatalf("FindBetween() = (%v, %v, %v), want persisted friendship", persisted, ok, err)
	}
}

func TestInviteFriendRejectsMissingIdentity(t *testing.T) {
	useCase := NewInviteFriend(testDirectory(t, "alice"), newFriendshipRepositoryFake())

	_, err := useCase.Execute(context.Background(), InviteFriendCommand{InviterID: "alice", InviteeID: "missing"})
	if !errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("Execute() error = %v, want %v", err, ErrIdentityNotFound)
	}
}

func TestInviteFriendRejectsSelfInvite(t *testing.T) {
	useCase := NewInviteFriend(testDirectory(t, "alice"), newFriendshipRepositoryFake())

	_, err := useCase.Execute(context.Background(), InviteFriendCommand{InviterID: "alice", InviteeID: "alice"})
	if !errors.Is(err, domainfriendship.ErrSelfInvite) {
		t.Fatalf("Execute() error = %v, want %v", err, domainfriendship.ErrSelfInvite)
	}
}

func TestInviteFriendRejectsExistingParticipantPair(t *testing.T) {
	ctx := context.Background()
	repository := newFriendshipRepositoryFake()
	useCase := NewInviteFriend(testDirectory(t, "alice", "bob"), repository)
	if _, err := useCase.Execute(ctx, InviteFriendCommand{InviterID: "alice", InviteeID: "bob"}); err != nil {
		t.Fatalf("first Execute() error = %v", err)
	}

	_, err := useCase.Execute(ctx, InviteFriendCommand{InviterID: "bob", InviteeID: "alice"})
	if !errors.Is(err, ErrFriendshipAlreadyExists) {
		t.Fatalf("reverse duplicate Execute() error = %v, want %v", err, ErrFriendshipAlreadyExists)
	}
}

func TestAcceptFriendshipOnlyAllowsDesignatedInvitee(t *testing.T) {
	ctx := context.Background()
	repository := newFriendshipRepositoryFake()
	directory := testDirectory(t, "alice", "bob", "charlie")
	invite := NewInviteFriend(directory, repository)
	accept := NewAcceptFriendship(directory, repository)
	if _, err := invite.Execute(ctx, InviteFriendCommand{InviterID: "alice", InviteeID: "bob"}); err != nil {
		t.Fatalf("invite Execute() error = %v", err)
	}

	for _, actor := range []string{"alice", "charlie"} {
		_, err := accept.Execute(ctx, AcceptFriendshipCommand{InviterID: "alice", InviteeID: "bob", ActorID: actor})
		if !errors.Is(err, domainfriendship.ErrOnlyInviteeCanAccept) {
			t.Fatalf("actor %q error = %v, want %v", actor, err, domainfriendship.ErrOnlyInviteeCanAccept)
		}
	}
}

func TestAcceptFriendshipActivatesSingleAggregateAndRepeatedAcceptFails(t *testing.T) {
	ctx := context.Background()
	repository := newFriendshipRepositoryFake()
	directory := testDirectory(t, "alice", "bob")
	invite := NewInviteFriend(directory, repository)
	accept := NewAcceptFriendship(directory, repository)
	if _, err := invite.Execute(ctx, InviteFriendCommand{InviterID: "alice", InviteeID: "bob"}); err != nil {
		t.Fatalf("invite Execute() error = %v", err)
	}

	active, err := accept.Execute(ctx, AcceptFriendshipCommand{InviterID: "alice", InviteeID: "bob", ActorID: "bob"})
	if err != nil {
		t.Fatalf("accept Execute() error = %v", err)
	}
	if !active.IsActive() {
		t.Fatal("accepted friendship must be active")
	}

	_, err = accept.Execute(ctx, AcceptFriendshipCommand{InviterID: "alice", InviteeID: "bob", ActorID: "bob"})
	if !errors.Is(err, domainfriendship.ErrAlreadyActive) {
		t.Fatalf("repeated accept error = %v, want %v", err, domainfriendship.ErrAlreadyActive)
	}
}

func TestAcceptFriendshipRejectsMissingIdentityAndInvitation(t *testing.T) {
	ctx := context.Background()
	repository := newFriendshipRepositoryFake()
	directory := testDirectory(t, "alice", "bob")
	accept := NewAcceptFriendship(directory, repository)

	_, err := accept.Execute(ctx, AcceptFriendshipCommand{InviterID: "alice", InviteeID: "bob", ActorID: "missing"})
	if !errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("missing actor error = %v, want %v", err, ErrIdentityNotFound)
	}

	_, err = accept.Execute(ctx, AcceptFriendshipCommand{InviterID: "alice", InviteeID: "bob", ActorID: "bob"})
	if !errors.Is(err, ErrFriendshipNotFound) {
		t.Fatalf("missing invitation error = %v, want %v", err, ErrFriendshipNotFound)
	}
}
