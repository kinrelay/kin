package socialcontext

import (
	"context"
	"errors"

	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

var (
	// ErrActivityOwnerMismatch protects owner-private Activity inputs if a read adapter violates its owner-scoped contract.
	ErrActivityOwnerMismatch = errors.New("significance activity owner mismatch")
)

// ActivityForSignificance is the minimal normalized, owner-private Activity projection required by this use case.
type ActivityForSignificance struct {
	ID      string
	OwnerID domainidentity.ID
	Content string
}

// ActivitySignificanceReader exposes only owner-private normalized Activities needed for significance evaluation.
type ActivitySignificanceReader interface {
	ListOwnerPrivateNormalized(ctx context.Context, ownerID domainidentity.ID) ([]ActivityForSignificance, error)
}

// EvaluateActivitySignificanceQuery evaluates the requester's own private Activities.
type EvaluateActivitySignificanceQuery struct {
	RequesterID string
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

	activities, err := uc.reader.ListOwnerPrivateNormalized(ctx, requesterID)
	if err != nil {
		return nil, err
	}
	if len(activities) == 0 {
		return []domainsocialcontext.SignificanceDecision{}, nil
	}

	signals := make([]domainsocialcontext.SignificanceSignal, 0, len(activities))
	for _, activity := range activities {
		if activity.OwnerID != requesterID {
			return nil, ErrActivityOwnerMismatch
		}
		signals = append(signals, domainsocialcontext.SignificanceSignal{
			ActivityID: activity.ID,
			Content:    activity.Content,
		})
	}

	return domainsocialcontext.EvaluateSignificance(signals), nil
}
