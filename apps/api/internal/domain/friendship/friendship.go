package friendship

import (
	"errors"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

var (
	// ErrSelfInvite indicates that an identity attempted to invite itself.
	ErrSelfInvite = errors.New("friendship participants must be distinct")
	// ErrOnlyInviteeCanAccept indicates that someone other than the designated invitee attempted acceptance.
	ErrOnlyInviteeCanAccept = errors.New("only the designated invitee can accept friendship")
	// ErrAlreadyActive indicates that an already active friendship was accepted again.
	ErrAlreadyActive = errors.New("friendship is already active")
)

// Friendship represents one close-friend relationship invitation and its acceptance state.
type Friendship struct {
	inviterID domainidentity.ID
	inviteeID domainidentity.ID
	active    bool
}

// Invite creates a pending friendship invitation between two distinct identities.
func Invite(inviterID, inviteeID domainidentity.ID) (Friendship, error) {
	inviterID, err := domainidentity.NewID(string(inviterID))
	if err != nil {
		return Friendship{}, err
	}
	inviteeID, err = domainidentity.NewID(string(inviteeID))
	if err != nil {
		return Friendship{}, err
	}
	if inviterID == inviteeID {
		return Friendship{}, ErrSelfInvite
	}

	return Friendship{inviterID: inviterID, inviteeID: inviteeID}, nil
}

// InviterID returns the identity that initiated the friendship invitation.
func (f Friendship) InviterID() domainidentity.ID {
	return f.inviterID
}

// InviteeID returns the identity designated to accept the friendship invitation.
func (f Friendship) InviteeID() domainidentity.ID {
	return f.inviteeID
}

// IsActive reports whether the invitation has been accepted.
func (f Friendship) IsActive() bool {
	return f.active
}

// Accept activates the friendship only when performed by the designated invitee.
func (f *Friendship) Accept(actorID domainidentity.ID) error {
	if f.active {
		return ErrAlreadyActive
	}

	actorID, err := domainidentity.NewID(string(actorID))
	if err != nil {
		return err
	}
	if actorID != f.inviteeID {
		return ErrOnlyInviteeCanAccept
	}

	f.active = true
	return nil
}
