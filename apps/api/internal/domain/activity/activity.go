package activity

import (
	"errors"
	"strings"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

var (
	// ErrInvalidID indicates that an Activity ID is empty after normalization.
	ErrInvalidID = errors.New("activity id must not be empty")
	// ErrEmptyContent indicates that a contribution does not contain meaningful content.
	ErrEmptyContent = errors.New("activity content must not be empty")
	// ErrInvalidTimestamp indicates that required Activity time metadata is missing.
	ErrInvalidTimestamp = errors.New("activity timestamp must not be zero")
)

// ID identifies one normalized Kin Activity.
type ID string

// Content is the normalized human-meaningful content of an Activity.
type Content struct {
	value string
}

// NewContent validates and normalizes Activity content.
func NewContent(value string) (Content, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return Content{}, ErrEmptyContent
	}
	return Content{value: value}, nil
}

// String returns normalized Activity content.
func (c Content) String() string {
	return c.value
}

// Provenance describes how an Activity entered Kin.
type Provenance string

// ProvenanceManual means the owner explicitly contributed the Activity to Kin.
const ProvenanceManual Provenance = "manual"

// Activity is a normalized, owner-private signal authorized for Kin to retain.
type Activity struct {
	id            ID
	ownerID       domainidentity.ID
	content       Content
	provenance    Provenance
	occurredAt    time.Time
	contributedAt time.Time
}

// NewManual creates a normalized private Activity from an explicit manual contribution.
func NewManual(idValue string, ownerID domainidentity.ID, content Content, occurredAt, contributedAt time.Time) (Activity, error) {
	idValue = strings.TrimSpace(idValue)
	if idValue == "" {
		return Activity{}, ErrInvalidID
	}
	validatedOwnerID, err := domainidentity.NewID(string(ownerID))
	if err != nil {
		return Activity{}, err
	}
	if strings.TrimSpace(content.String()) == "" {
		return Activity{}, ErrEmptyContent
	}
	if occurredAt.IsZero() || contributedAt.IsZero() {
		return Activity{}, ErrInvalidTimestamp
	}

	return Activity{
		id:            ID(idValue),
		ownerID:       validatedOwnerID,
		content:       content,
		provenance:    ProvenanceManual,
		occurredAt:    occurredAt,
		contributedAt: contributedAt,
	}, nil
}

// ID returns the Activity's stable identifier.
func (a Activity) ID() ID {
	return a.id
}

// OwnerID returns the Identity that contributed and owns the Activity.
func (a Activity) OwnerID() domainidentity.ID {
	return a.ownerID
}

// Content returns normalized Activity content.
func (a Activity) Content() Content {
	return a.content
}

// Provenance returns how the Activity entered Kin.
func (a Activity) Provenance() Provenance {
	return a.provenance
}

// IsPrivate reports the MVP 1 visibility contract: Activities are private inputs.
func (a Activity) IsPrivate() bool {
	return true
}

// OccurredAt returns when the contributed activity occurred.
func (a Activity) OccurredAt() time.Time {
	return a.occurredAt
}

// ContributedAt returns when the owner explicitly contributed the Activity to Kin.
func (a Activity) ContributedAt() time.Time {
	return a.contributedAt
}
