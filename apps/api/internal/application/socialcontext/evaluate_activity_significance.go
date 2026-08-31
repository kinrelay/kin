package socialcontext

import (
	"context"
	"errors"
	"sort"
	"time"

	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

var (
	// ErrActivityOwnerMismatch protects owner-private Activity inputs if a read adapter violates its owner-scoped contract.
	ErrActivityOwnerMismatch = errors.New("significance activity owner mismatch")
)

// ActivityForSignificance is the minimal normalized, owner-private Activity projection required by this use case.
type ActivityForSignificance struct {
	ID            string
	OwnerID       domainidentity.ID
	Content       string
	OccurredAt    time.Time
	ContributedAt time.Time
}

// ActivitySignificanceReader exposes only the requested owner-private normalized Activity batch needed for significance evaluation.
type ActivitySignificanceReader interface {
	ListOwnerPrivateNormalized(ctx context.Context, ownerID domainidentity.ID, activityIDs []string) ([]ActivityForSignificance, error)
}

// EvaluateActivitySignificanceQuery evaluates an explicit batch of the requester's own private Activities.
type EvaluateActivitySignificanceQuery struct {
	RequesterID string
	ActivityIDs []string
}

// EvaluateActivitySignificance orchestrates owner-scoped Activity reading and deterministic domain policy.
type EvaluateActivitySignificance struct {
	reader ActivitySignificanceReader
}

// NewEvaluateActivitySignificance constructs the significance use case around an explicit Activity read port.
func NewEvaluateActivitySignificance(reader ActivitySignificanceReader) EvaluateActivitySignificance {
	return EvaluateActivitySignificance{reader: reader}
}

// Execute returns Activity-scoped significance decisions without creating Context Candidate or Social Context state.
func (uc EvaluateActivitySignificance) Execute(ctx context.Context, query EvaluateActivitySignificanceQuery) ([]domainsocialcontext.SignificanceDecision, error) {
	requesterID, err := domainidentity.NewID(query.RequesterID)
	if err != nil {
		return nil, err
	}
	if len(query.ActivityIDs) == 0 {
		return []domainsocialcontext.SignificanceDecision{}, nil
	}

	activityIDs := append([]string(nil), query.ActivityIDs...)
	activities, err := uc.reader.ListOwnerPrivateNormalized(ctx, requesterID, activityIDs)
	if err != nil {
		return nil, err
	}
	if len(activities) == 0 {
		return []domainsocialcontext.SignificanceDecision{}, nil
	}

	for _, activity := range activities {
		if activity.OwnerID != requesterID {
			return nil, ErrActivityOwnerMismatch
		}
	}
	// Duplicate significance is a chronology-sensitive policy: when equivalent
	// signals repeat, the newest occurrence is the representative. Do not trust
	// adapter/request order to encode that meaning. ContributedAt is the explicit
	// secondary chronology when two Activities share the same OccurredAt.
	sort.SliceStable(activities, func(i, j int) bool {
		left, right := activities[i], activities[j]
		if !left.OccurredAt.IsZero() && !right.OccurredAt.IsZero() && !left.OccurredAt.Equal(right.OccurredAt) {
			return left.OccurredAt.Before(right.OccurredAt)
		}
		if !left.ContributedAt.IsZero() && !right.ContributedAt.IsZero() && !left.ContributedAt.Equal(right.ContributedAt) {
			return left.ContributedAt.Before(right.ContributedAt)
		}
		return false
	})

	signals := make([]domainsocialcontext.SignificanceSignal, 0, len(activities))
	for _, activity := range activities {
		signals = append(signals, domainsocialcontext.SignificanceSignal{
			ActivityID: activity.ID,
			Content:    activity.Content,
		})
	}

	return domainsocialcontext.EvaluateSignificance(signals), nil
}
