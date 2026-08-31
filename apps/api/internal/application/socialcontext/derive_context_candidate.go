package socialcontext

import (
	"context"
	"errors"
	"sort"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

var (
	ErrContextActivityOwnerMismatch = errors.New("context activity owner mismatch")
	ErrContextActivityNotRequested  = errors.New("context activity was not requested")
)

type ActivityForContext struct {
	ID            string
	OwnerID       domainidentity.ID
	Content       string
	OccurredAt    time.Time
	ContributedAt time.Time
}

type ContextActivityReader interface {
	ListOwnerPrivateNormalized(ctx context.Context, ownerID domainidentity.ID, activityIDs []string) ([]ActivityForContext, error)
}

type ContextGenerationActivity struct {
	ID      string
	Content string
}

type ContextGenerationInput struct {
	Activities []ContextGenerationActivity
}

type GeneratedContext struct {
	Meaning    string
	Provenance []string
}

type ContextGenerator interface {
	Generate(ctx context.Context, input ContextGenerationInput) (GeneratedContext, error)
}

type SocialContextRepository interface {
	SaveIfAbsent(ctx context.Context, ownerID domainidentity.ID, socialContext domainsocialcontext.SocialContext) (bool, error)
}

type DeriveContextCandidateCommand struct {
	RequesterID string
	ActivityIDs []string
}

type DerivationStatus string

const (
	DerivationPromoted   DerivationStatus = "promoted"
	DerivationRejected   DerivationStatus = "rejected"
	DerivationSuppressed DerivationStatus = "suppressed"
)

type DerivationOutcome struct {
	Status DerivationStatus
	Reason error
}

type DeriveContextCandidate struct {
	reader     ContextActivityReader
	generator  ContextGenerator
	repository SocialContextRepository
}

func NewDeriveContextCandidate(reader ContextActivityReader, generator ContextGenerator, repository SocialContextRepository) DeriveContextCandidate {
	return DeriveContextCandidate{reader: reader, generator: generator, repository: repository}
}

func (uc DeriveContextCandidate) Execute(ctx context.Context, command DeriveContextCandidateCommand) (DerivationOutcome, error) {
	requesterID, err := domainidentity.NewID(command.RequesterID)
	if err != nil {
		return DerivationOutcome{}, err
	}
	if len(command.ActivityIDs) == 0 {
		return DerivationOutcome{Status: DerivationSuppressed}, nil
	}

	requestedIDs := make(map[string]struct{}, len(command.ActivityIDs))
	for _, id := range command.ActivityIDs {
		requestedIDs[id] = struct{}{}
	}
	activities, err := uc.reader.ListOwnerPrivateNormalized(ctx, requesterID, append([]string(nil), command.ActivityIDs...))
	if err != nil {
		return DerivationOutcome{}, err
	}

	for _, activity := range activities {
		if activity.OwnerID != requesterID {
			return DerivationOutcome{}, ErrContextActivityOwnerMismatch
		}
		if _, requested := requestedIDs[activity.ID]; !requested {
			return DerivationOutcome{}, ErrContextActivityNotRequested
		}
	}
	// Reversal reconciliation represents current state within the requested
	// derivation batch, so derive from occurrence chronology rather than
	// caller-supplied Activity ID order. ContributedAt is the explicit secondary
	// chronology when two Activities share the same OccurredAt.
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
	byID := make(map[string]ActivityForContext, len(activities))
	for _, activity := range activities {
		byID[activity.ID] = activity
		signals = append(signals, domainsocialcontext.SignificanceSignal{ActivityID: activity.ID, Content: activity.Content})
	}

	decisions := domainsocialcontext.EvaluateSignificance(signals)
	eligible := make([]ContextGenerationActivity, 0, len(decisions))
	eligibleSources := make([]domainsocialcontext.SourceActivity, 0, len(decisions))
	for _, decision := range decisions {
		if decision.Status != domainsocialcontext.SignificanceEligible {
			continue
		}
		activity := byID[decision.ActivityID]
		eligible = append(eligible, ContextGenerationActivity{ID: activity.ID, Content: activity.Content})
		eligibleSources = append(eligibleSources, domainsocialcontext.SourceActivity{ID: activity.ID, Content: activity.Content})
	}
	if len(eligible) == 0 {
		return DerivationOutcome{Status: DerivationSuppressed}, nil
	}

	generated, err := uc.generator.Generate(ctx, ContextGenerationInput{Activities: eligible})
	if err != nil {
		return DerivationOutcome{}, err
	}
	if generated.Meaning == "" {
		return DerivationOutcome{Status: DerivationSuppressed}, nil
	}

	candidate, err := domainsocialcontext.NewContextCandidate(generated.Meaning, generated.Provenance)
	if err != nil {
		return DerivationOutcome{Status: DerivationRejected, Reason: err}, nil
	}
	socialContext, err := domainsocialcontext.PromoteContextCandidate(candidate, eligibleSources)
	if err != nil {
		return DerivationOutcome{Status: DerivationRejected, Reason: err}, nil
	}
	inserted, err := uc.repository.SaveIfAbsent(ctx, requesterID, socialContext)
	if err != nil {
		return DerivationOutcome{}, err
	}
	if !inserted {
		return DerivationOutcome{Status: DerivationSuppressed}, nil
	}

	return DerivationOutcome{Status: DerivationPromoted}, nil
}
