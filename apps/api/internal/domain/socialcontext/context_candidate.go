package socialcontext

import (
	"errors"
	"strings"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

var (
	// ErrInvalidCandidateID indicates that a Context Candidate ID is empty after normalization.
	ErrInvalidCandidateID = errors.New("context candidate id must not be empty")
	// ErrEmptyMeaning indicates that generated candidate meaning is blank.
	ErrEmptyMeaning = errors.New("context candidate meaning must not be empty")
	// ErrMissingSourceActivity indicates that candidate provenance has no valid Activity source.
	ErrMissingSourceActivity = errors.New("context candidate must reference at least one source activity")
	// ErrInvalidCandidateTimestamp indicates that candidate generation time is missing.
	ErrInvalidCandidateTimestamp = errors.New("context candidate generated timestamp must not be zero")
)

// CandidateID identifies one unvalidated Context Candidate.
type CandidateID string

// CandidateState describes the lifecycle stage of a generated candidate.
type CandidateState string

// StateCandidate means the generated meaning has not yet been validated/promoted to Social Context.
const StateCandidate CandidateState = "candidate"

// Meaning is normalized derived social meaning proposed for a Context Candidate.
type Meaning struct {
	value string
}

// NewMeaning validates and normalizes generated candidate meaning.
func NewMeaning(value string) (Meaning, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return Meaning{}, ErrEmptyMeaning
	}
	return Meaning{value: value}, nil
}

// String returns normalized derived meaning.
func (m Meaning) String() string {
	return m.value
}

// IsPureReplayOf reports whether meaning differs from raw Activity content only by case or whitespace.
func (m Meaning) IsPureReplayOf(raw string) bool {
	return normalizeReplayComparable(m.value) == normalizeReplayComparable(raw)
}

func normalizeReplayComparable(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(value), " "))
}

// ContextCandidate is owner-private generated meaning awaiting explicit Social Context validation.
type ContextCandidate struct {
	id                CandidateID
	ownerID           domainidentity.ID
	meaning           Meaning
	sourceActivityIDs []string
	generatedAt       time.Time
}

// NewContextCandidate creates an owner-private, unvalidated candidate with abstract Activity provenance.
func NewContextCandidate(idValue string, ownerID domainidentity.ID, meaning Meaning, sourceActivityIDs []string, generatedAt time.Time) (ContextCandidate, error) {
	idValue = strings.TrimSpace(idValue)
	if idValue == "" {
		return ContextCandidate{}, ErrInvalidCandidateID
	}
	validatedOwnerID, err := domainidentity.NewID(string(ownerID))
	if err != nil {
		return ContextCandidate{}, err
	}
	validatedMeaning, err := NewMeaning(meaning.String())
	if err != nil {
		return ContextCandidate{}, err
	}
	if generatedAt.IsZero() {
		return ContextCandidate{}, ErrInvalidCandidateTimestamp
	}

	normalizedSources := make([]string, 0, len(sourceActivityIDs))
	seen := make(map[string]struct{}, len(sourceActivityIDs))
	for _, sourceID := range sourceActivityIDs {
		sourceID = strings.TrimSpace(sourceID)
		if sourceID == "" {
			return ContextCandidate{}, ErrMissingSourceActivity
		}
		if _, exists := seen[sourceID]; exists {
			continue
		}
		seen[sourceID] = struct{}{}
		normalizedSources = append(normalizedSources, sourceID)
	}
	if len(normalizedSources) == 0 {
		return ContextCandidate{}, ErrMissingSourceActivity
	}

	return ContextCandidate{
		id:                CandidateID(idValue),
		ownerID:           validatedOwnerID,
		meaning:           validatedMeaning,
		sourceActivityIDs: normalizedSources,
		generatedAt:       generatedAt,
	}, nil
}

// ID returns the stable candidate identifier.
func (c ContextCandidate) ID() CandidateID {
	return c.id
}

// OwnerID returns the Identity whose Activities produced this candidate.
func (c ContextCandidate) OwnerID() domainidentity.ID {
	return c.ownerID
}

// Meaning returns the generated derived meaning.
func (c ContextCandidate) Meaning() Meaning {
	return c.meaning
}

// SourceActivityIDs returns a defensive copy of abstract Activity provenance.
func (c ContextCandidate) SourceActivityIDs() []string {
	return append([]string(nil), c.sourceActivityIDs...)
}

// GeneratedAt returns when the candidate was generated.
func (c ContextCandidate) GeneratedAt() time.Time {
	return c.generatedAt
}

// State reports that this value is still an unvalidated candidate.
func (c ContextCandidate) State() CandidateState {
	return StateCandidate
}

// IsPrivate reports that MVP 2 candidates are owner-private inputs, never friend-visible output.
func (c ContextCandidate) IsPrivate() bool {
	return true
}

// IsValidatedSocialContext is false until a later explicit validation/promotion use case succeeds.
func (c ContextCandidate) IsValidatedSocialContext() bool {
	return false
}
