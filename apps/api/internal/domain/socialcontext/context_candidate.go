package socialcontext

import (
	"errors"
	"strings"
)

var (
	ErrBlankContextMeaning      = errors.New("context meaning is blank")
	ErrMissingContextProvenance = errors.New("context provenance is missing")
	ErrSourceReplay             = errors.New("context meaning only replays source activity")
)

// ContextCandidate is derived meaning awaiting Kin-owned validation/promotion.
type ContextCandidate struct {
	meaning    string
	provenance []string
}

// SourceActivity is the minimal authorized source material needed to validate provenance and replay.
type SourceActivity struct {
	ID      string
	Content string
}

// SocialContext is owner-private derived social meaning promoted from a valid candidate.
type SocialContext struct {
	meaning    string
	provenance []string
}

func NewContextCandidate(meaning string, provenance []string) (ContextCandidate, error) {
	normalizedMeaning := normalizeContextText(meaning)
	if normalizedMeaning == "" {
		return ContextCandidate{}, ErrBlankContextMeaning
	}

	normalizedProvenance := normalizeProvenance(provenance)
	if len(normalizedProvenance) == 0 {
		return ContextCandidate{}, ErrMissingContextProvenance
	}

	return ContextCandidate{meaning: normalizedMeaning, provenance: normalizedProvenance}, nil
}

func PromoteContextCandidate(candidate ContextCandidate, sources []SourceActivity) (SocialContext, error) {
	authorized := make(map[string]string, len(sources))
	for _, source := range sources {
		id := strings.TrimSpace(source.ID)
		if id != "" {
			authorized[id] = normalizeContextText(source.Content)
		}
	}

	for _, sourceID := range candidate.provenance {
		content, ok := authorized[sourceID]
		if !ok {
			return SocialContext{}, ErrMissingContextProvenance
		}
		if strings.EqualFold(candidate.meaning, content) {
			return SocialContext{}, ErrSourceReplay
		}
	}

	return SocialContext{
		meaning:    candidate.meaning,
		provenance: append([]string(nil), candidate.provenance...),
	}, nil
}

func (c SocialContext) Meaning() string { return c.meaning }

func (c SocialContext) Provenance() []string {
	return append([]string(nil), c.provenance...)
}

func (c SocialContext) IsPrivate() bool { return true }

func normalizeProvenance(ids []string) []string {
	seen := make(map[string]struct{}, len(ids))
	normalized := make([]string, 0, len(ids))
	for _, raw := range ids {
		id := strings.TrimSpace(raw)
		if id == "" {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		normalized = append(normalized, id)
	}
	return normalized
}

func normalizeContextText(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
