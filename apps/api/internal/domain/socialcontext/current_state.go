package socialcontext

import (
	"errors"
	"strings"
	"time"
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

// CurrentStateMutation describes one semantic component mutation inside an
// owner-scoped atomic reconciliation batch. A nil Replacement retires the
// component at OccurredAt; a non-nil Replacement establishes its new state.
type CurrentStateMutation struct {
	SemanticID  SemanticIdentity
	OccurredAt  time.Time
	Replacement *SocialContext
}
