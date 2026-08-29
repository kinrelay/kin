package identity

import (
	"errors"
	"strings"
)

var ErrInvalidID = errors.New("identity id must not be empty")

type ID string

func NewID(value string) (ID, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ErrInvalidID
	}

	return ID(value), nil
}

type Identity struct {
	id ID
}

func New(id ID) (Identity, error) {
	normalizedID, err := NewID(string(id))
	if err != nil {
		return Identity{}, err
	}

	return Identity{id: normalizedID}, nil
}

func (i Identity) ID() ID {
	return i.id
}
