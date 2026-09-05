package friendpulse

import (
	"context"
	"errors"
	"sort"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainprivacy "github.com/kinrelay/kin/apps/api/internal/domain/privacy"
)

const maxPulseItems = 3

var ErrFriendPulseUnauthorized = errors.New("requester is not an active friend")

// Query identifies the authenticated viewer and the active friend whose pulse is requested.
type Query struct {
	AuthenticatedViewerID string
	FriendID              string
}

// Candidate is ranking metadata for a Social Context that may contribute a friend-visible projection.
// It deliberately carries no raw activity or Social Context meaning.
type Candidate struct {
	SocialContextID string
	SignalScore     int
	ObservedAt      time.Time
}

// Item contains only the privacy-approved projection exposed to the viewer.
type Item struct {
	Projection domainprivacy.ContextProjection
}

// Pulse is intentionally capped rather than representing a chronological activity history.
type Pulse struct {
	Items []Item
}

// ActiveFriendshipReader provides read-only relationship authorization evidence.
type ActiveFriendshipReader interface {
	IsActiveBetween(ctx context.Context, first, second domainidentity.ID) (bool, error)
}

// CandidateReader supplies ranking metadata without exposing raw private content.
type CandidateReader interface {
	ListForOwner(ctx context.Context, ownerID domainidentity.ID) ([]Candidate, error)
}

// ContextProjector resolves the viewer-specific Context Projection using the canonical privacy policy.
type ContextProjector interface {
	Project(
		ctx context.Context,
		viewerID domainidentity.ID,
		ownerID domainidentity.ID,
		socialContextID string,
	) (domainprivacy.ContextProjection, bool, error)
}

// GetFriendPulse returns a small deterministic set of privacy-approved context projections.
type GetFriendPulse struct {
	friendships ActiveFriendshipReader
	candidates  CandidateReader
	projector   ContextProjector
}

func NewGetFriendPulse(
	friendships ActiveFriendshipReader,
	candidates CandidateReader,
	projector ContextProjector,
) GetFriendPulse {
	return GetFriendPulse{friendships: friendships, candidates: candidates, projector: projector}
}

func (uc GetFriendPulse) Execute(ctx context.Context, query Query) (Pulse, error) {
	viewerID, err := domainidentity.NewID(query.AuthenticatedViewerID)
	if err != nil {
		return Pulse{}, err
	}
	friendID, err := domainidentity.NewID(query.FriendID)
	if err != nil {
		return Pulse{}, err
	}

	active, err := uc.friendships.IsActiveBetween(ctx, viewerID, friendID)
	if err != nil {
		return Pulse{}, err
	}
	if !active {
		return Pulse{}, ErrFriendPulseUnauthorized
	}

	candidates, err := uc.candidates.ListForOwner(ctx, friendID)
	if err != nil {
		return Pulse{}, err
	}
	candidates = append([]Candidate(nil), candidates...)
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].SignalScore != candidates[j].SignalScore {
			return candidates[i].SignalScore > candidates[j].SignalScore
		}
		if !candidates[i].ObservedAt.Equal(candidates[j].ObservedAt) {
			return candidates[i].ObservedAt.After(candidates[j].ObservedAt)
		}
		return candidates[i].SocialContextID < candidates[j].SocialContextID
	})

	pulse := Pulse{Items: make([]Item, 0, maxPulseItems)}
	for _, candidate := range candidates {
		projection, visible, err := uc.projector.Project(ctx, viewerID, friendID, candidate.SocialContextID)
		if err != nil {
			return Pulse{}, err
		}
		if !visible {
			continue
		}
		pulse.Items = append(pulse.Items, Item{Projection: projection})
		if len(pulse.Items) == maxPulseItems {
			break
		}
	}

	return pulse, nil
}
