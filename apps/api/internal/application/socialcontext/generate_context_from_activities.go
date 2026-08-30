package socialcontext

import (
	"context"
	"errors"
	"strings"
	"time"

	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

var (
	// ErrActivityOwnerMismatch prevents one owner's candidate from incorporating another owner's private Activity.
	ErrActivityOwnerMismatch = errors.New("context candidate source activity owner mismatch")
	// ErrInvalidActivitySignal indicates that an Activity read port returned an invalid normalized signal.
	ErrInvalidActivitySignal = errors.New("invalid normalized activity signal")
	// ErrPureActivityReplay prevents a single raw Activity from being re-published as derived context.
	ErrPureActivityReplay = errors.New("generated context must not replay raw activity")
)

// ActivitySignal is the provider-neutral normalized Activity shape consumed by Social Context orchestration.
type ActivitySignal struct {
	ID            string
	OwnerID       domainidentity.ID
	Content       string
	Provenance    string
	OccurredAt    time.Time
	ContributedAt time.Time
}

// ActivitySignalReader retrieves normalized private Activity signals without exposing write aggregates.
type ActivitySignalReader interface {
	FindByIDs(ctx context.Context, ids []string) ([]ActivitySignal, error)
}

// ContextGenerationInput is the provider-neutral input to a derived-context generator.
type ContextGenerationInput struct {
	OwnerID domainidentity.ID
	Signals []ActivitySignal
}

// GeneratedContext is normalized generator output before Kin domain validation creates a candidate.
type GeneratedContext struct {
	Meaning string
}

// ContextGenerator describes the business capability to derive higher-level meaning from normalized Activity signals.
type ContextGenerator interface {
	Generate(ctx context.Context, input ContextGenerationInput) (GeneratedContext, error)
}

// CandidateRepository persists owner-private, unvalidated Context Candidates.
type CandidateRepository interface {
	Save(ctx context.Context, candidate domainsocialcontext.ContextCandidate) error
}

// CandidateIDGenerator provides stable candidate identifiers without coupling the use case to randomness infrastructure.
type CandidateIDGenerator interface {
	NewCandidateID(ctx context.Context) (string, error)
}

// Clock provides deterministic generation timestamps.
type Clock interface {
	Now() time.Time
}

// SignificanceReason explains the deterministic MVP 2 significance gate outcome.
type SignificanceReason string

const (
	// SignificanceNone means there are no usable signals.
	SignificanceNone SignificanceReason = "no-signals"
	// SignificanceDuplicateOnly means multiple signals carry no more meaning than the same repeated Activity.
	SignificanceDuplicateOnly SignificanceReason = "duplicate-only"
	// SignificancePotential means the signals may support higher-level derived meaning and may reach the generator.
	SignificancePotential SignificanceReason = "potentially-significant"
)

// SignificanceInput is the provider-neutral input to the deterministic MVP significance gate.
type SignificanceInput struct {
	Signals []ActivitySignal
}

// SignificanceResult records whether generation is authorized by the deterministic significance gate.
type SignificanceResult struct {
	Significant bool
	Reason      SignificanceReason
}

// EvaluateSignificance applies the minimal MVP 2 gate without invoking an AI or external provider.
func EvaluateSignificance(input SignificanceInput) SignificanceResult {
	if len(input.Signals) == 0 {
		return SignificanceResult{Reason: SignificanceNone}
	}
	if len(input.Signals) > 1 {
		first := normalizeComparable(input.Signals[0].Content)
		allSame := true
		for _, item := range input.Signals[1:] {
			if normalizeComparable(item.Content) != first {
				allSame = false
				break
			}
		}
		if allSame {
			return SignificanceResult{Reason: SignificanceDuplicateOnly}
		}
	}
	return SignificanceResult{Significant: true, Reason: SignificancePotential}
}

// GenerateContextFromActivitiesCommand identifies the owner and private Activity sources to evaluate.
type GenerateContextFromActivitiesCommand struct {
	OwnerID           string
	SourceActivityIDs []string
}

// GenerateContextResult distinguishes intentional suppression from a generated candidate.
type GenerateContextResult struct {
	Candidate  *domainsocialcontext.ContextCandidate
	Suppressed bool
}

// GenerateContextFromActivities evaluates significance, invokes a provider-neutral generator only when warranted, and persists a candidate.
type GenerateContextFromActivities struct {
	reader      ActivitySignalReader
	generator   ContextGenerator
	repository  CandidateRepository
	idGenerator CandidateIDGenerator
	clock       Clock
}

// NewGenerateContextFromActivities constructs the candidate-generation use case from explicit ports.
func NewGenerateContextFromActivities(reader ActivitySignalReader, generator ContextGenerator, repository CandidateRepository, idGenerator CandidateIDGenerator, clock Clock) GenerateContextFromActivities {
	return GenerateContextFromActivities{
		reader: reader, generator: generator, repository: repository, idGenerator: idGenerator, clock: clock,
	}
}

// Execute generates an owner-private Context Candidate only when the source set passes the MVP significance gate.
func (uc GenerateContextFromActivities) Execute(ctx context.Context, command GenerateContextFromActivitiesCommand) (GenerateContextResult, error) {
	ownerID, err := domainidentity.NewID(command.OwnerID)
	if err != nil {
		return GenerateContextResult{}, err
	}

	requestedIDs := normalizeUniqueIDs(command.SourceActivityIDs)
	if len(requestedIDs) == 0 {
		return GenerateContextResult{Suppressed: true}, nil
	}

	signals, err := uc.reader.FindByIDs(ctx, requestedIDs)
	if err != nil {
		return GenerateContextResult{}, err
	}
	normalizedSignals, err := normalizeSignals(ownerID, signals)
	if err != nil {
		return GenerateContextResult{}, err
	}

	decision := EvaluateSignificance(SignificanceInput{Signals: normalizedSignals})
	if !decision.Significant {
		return GenerateContextResult{Suppressed: true}, nil
	}

	generated, err := uc.generator.Generate(ctx, ContextGenerationInput{
		OwnerID: ownerID,
		Signals: append([]ActivitySignal(nil), normalizedSignals...),
	})
	if err != nil {
		return GenerateContextResult{}, err
	}
	meaning, err := domainsocialcontext.NewMeaning(generated.Meaning)
	if err != nil {
		return GenerateContextResult{}, err
	}
	if len(normalizedSignals) == 1 && meaning.IsPureReplayOf(normalizedSignals[0].Content) {
		return GenerateContextResult{}, ErrPureActivityReplay
	}

	candidateID, err := uc.idGenerator.NewCandidateID(ctx)
	if err != nil {
		return GenerateContextResult{}, err
	}
	sourceIDs := make([]string, 0, len(normalizedSignals))
	for _, item := range normalizedSignals {
		sourceIDs = append(sourceIDs, item.ID)
	}
	candidate, err := domainsocialcontext.NewContextCandidate(candidateID, ownerID, meaning, sourceIDs, uc.clock.Now())
	if err != nil {
		return GenerateContextResult{}, err
	}
	if err := uc.repository.Save(ctx, candidate); err != nil {
		return GenerateContextResult{}, err
	}

	return GenerateContextResult{Candidate: &candidate}, nil
}

func normalizeUniqueIDs(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func normalizeSignals(ownerID domainidentity.ID, signals []ActivitySignal) ([]ActivitySignal, error) {
	result := make([]ActivitySignal, 0, len(signals))
	seen := make(map[string]struct{}, len(signals))
	for _, item := range signals {
		item.ID = strings.TrimSpace(item.ID)
		item.Content = strings.TrimSpace(item.Content)
		item.Provenance = strings.TrimSpace(item.Provenance)
		if item.OwnerID != ownerID {
			return nil, ErrActivityOwnerMismatch
		}
		if item.ID == "" || item.Content == "" || item.OccurredAt.IsZero() || item.ContributedAt.IsZero() {
			return nil, ErrInvalidActivitySignal
		}
		if _, exists := seen[item.ID]; exists {
			continue
		}
		seen[item.ID] = struct{}{}
		result = append(result, item)
	}
	return result, nil
}

func normalizeComparable(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(value), " "))
}
