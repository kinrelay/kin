package socialcontext

import (
	"errors"
	"strings"
)

var ErrBlankSemanticIdentity = errors.New("social context semantic identity is blank")

// SemanticIdentity identifies one owner-private social meaning lifecycle independently
// from the individual Activity provenance that currently supports it.
type SemanticIdentity string

func NewSemanticIdentity(value string) (SemanticIdentity, error) {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return "", ErrBlankSemanticIdentity
	}
	return SemanticIdentity(normalized), nil
}

func (id SemanticIdentity) String() string { return string(id) }
