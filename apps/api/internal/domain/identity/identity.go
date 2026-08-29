package identity

import (
	"errors"
	"strings"
)

// ErrInvalidID indicates that an identity ID is empty after normalization.
var ErrInvalidID = errors.New("identity id must not be empty")

// ID identifies a Kin identity.
type ID string

// NewID validates and normalizes an identity ID.
func NewID(value string) (ID, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidID
	}

	return ID(value), nil
}

// Identity represents a Kin identity with a stable validated ID.
type Identity struct {
	id ID
}

// New creates an Identity while enforcing the ID invariant.
func New(id ID) (Identity, error) {
	normalizedID, err := NewID(string(id))
	if err != nil {
		return Identity{}, err
	}

	return Identity{id: normalizedID}, nil
}

// ID returns the identity's stable ID.
func (i Identity) ID() ID {
	return i.id
}
