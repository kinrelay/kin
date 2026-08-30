package socialcontext

import (
	"context"
	"errors"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
)

var (
	ErrContextActivityOwnerMismatch = errors.New("context activity owner mismatch")
	ErrContextActivityNotRequested  = errors.New("context activity was not requested")
)

type ActivityForContext struct {
	ID      string
	OwnerID domainidentity.ID
	Content string
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
	Save(ctx context.Context, socialContext domainsocialcontext.SocialContext) error
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

	signals := make([]domainsocialcontext.SignificanceSignal, 0, len(activities))
	byID := make(map[string]ActivityForContext, len(activities))
	for _, activity := range activities {
		if activity.OwnerID != requesterID {
			return DerivationOutcome{}, ErrContextActivityOwnerMismatch
		}
		if _, requested := requestedIDs[activity.ID]; !requested {
			return DerivationOutcome{}, ErrContextActivityNotRequested
		}
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
	candidate, err := domainsocialcontext.NewContextCandidate(generated.Meaning, generated.Provenance)
	if err != nil {
		return DerivationOutcome{Status: DerivationRejected, Reason: err}, nil
	}
	socialContext, err := domainsocialcontext.PromoteContextCandidate(candidate, eligibleSources)
	if err != nil {
		return DerivationOutcome{Status: DerivationRejected, Reason: err}, nil
	}
	if err := uc.repository.Save(ctx, socialContext); err != nil {
		return DerivationOutcome{}, err
	}

	return DerivationOutcome{Status: DerivationPromoted}, nil
}
