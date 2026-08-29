package friendship

import (
	"errors"
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func mustIdentityID(t *testing.T, value string) domainidentity.ID {
	t.Helper()
	id, err := domainidentity.NewID(value)
	if err != nil {
		t.Fatalf("NewID(%q): %v", value, err)
	}
	return id
}

func TestInviteRequiresDistinctParticipants(t *testing.T) {
	alice := mustIdentityID(t, "alice")

	_, err := Invite(alice, alice)
	if !errors.Is(err, ErrSelfInvite) {
		t.Fatalf("Invite(self) error = %v, want %v", err, ErrSelfInvite)
	}
}

func TestInvitationRecordsParticipantsAndStartsPending(t *testing.T) {
	alice := mustIdentityID(t, "alice")
	bob := mustIdentityID(t, "bob")

	friendship, err := Invite(alice, bob)
	if err != nil {
		t.Fatalf("Invite() error = %v", err)
	}
	if friendship.InviterID() != alice {
		t.Fatalf("InviterID() = %q, want %q", friendship.InviterID(), alice)
	}
	if friendship.InviteeID() != bob {
		t.Fatalf("InviteeID() = %q, want %q", friendship.InviteeID(), bob)
	}
	if friendship.IsActive() {
		t.Fatal("new invitation must not be active")
	}
}

func TestOnlyDesignatedInviteeCanAccept(t *testing.T) {
	alice := mustIdentityID(t, "alice")
	bob := mustIdentityID(t, "bob")
	charlie := mustIdentityID(t, "charlie")

	tests := []struct {
		name  string
		actor domainidentity.ID
	}{
		{name: "inviter", actor: alice},
		{name: "third party", actor: charlie},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			friendship, err := Invite(alice, bob)
			if err != nil {
				t.Fatalf("Invite() error = %v", err)
			}

			err = friendship.Accept(tt.actor)
			if !errors.Is(err, ErrOnlyInviteeCanAccept) {
				t.Fatalf("Accept(%q) error = %v, want %v", tt.actor, err, ErrOnlyInviteeCanAccept)
			}
			if friendship.IsActive() {
				t.Fatal("unauthorized acceptance must not activate friendship")
			}
		})
	}
}

func TestInviteeAcceptanceActivatesFriendship(t *testing.T) {
	alice := mustIdentityID(t, "alice")
	bob := mustIdentityID(t, "bob")
	friendship, err := Invite(alice, bob)
	if err != nil {
		t.Fatalf("Invite() error = %v", err)
	}

	if err := friendship.Accept(bob); err != nil {
		t.Fatalf("Accept(invitee) error = %v", err)
	}
	if !friendship.IsActive() {
		t.Fatal("accepted friendship must be active")
	}
}

func TestRepeatedAcceptanceFailsExplicitly(t *testing.T) {
	alice := mustIdentityID(t, "alice")
	bob := mustIdentityID(t, "bob")
	friendship, err := Invite(alice, bob)
	if err != nil {
		t.Fatalf("Invite() error = %v", err)
	}
	if err := friendship.Accept(bob); err != nil {
		t.Fatalf("first Accept() error = %v", err)
	}

	err = friendship.Accept(bob)
	if !errors.Is(err, ErrAlreadyActive) {
		t.Fatalf("second Accept() error = %v, want %v", err, ErrAlreadyActive)
	}
}
