package socialcontext

import (
	"context"
	"errors"
	"testing"
	"time"

	domainsocialcontext "github.com/kinrelay/kin/apps/api/internal/domain/socialcontext"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

type activitySignalReaderFake struct {
	items    []ActivitySignal
	err      error
	calls    int
	askedIDs []string
}

func (f *activitySignalReaderFake) FindByIDs(_ context.Context, ids []string) ([]ActivitySignal, error) {
	f.calls++
	f.askedIDs = append([]string(nil), ids...)
	if f.err != nil {
		return nil, f.err
	}
	return append([]ActivitySignal(nil), f.items...), nil
}

type contextGeneratorFake struct {
	output GeneratedContext
	err    error
	calls  int
	input  ContextGenerationInput
}

func (f *contextGeneratorFake) Generate(_ context.Context, input ContextGenerationInput) (GeneratedContext, error) {
	f.calls++
	f.input = input
	if f.err != nil {
		return GeneratedContext{}, f.err
	}
	return f.output, nil
}

type candidateRepositoryFake struct {
	saved []domainsocialcontext.ContextCandidate
	err   error
}

func (f *candidateRepositoryFake) Save(_ context.Context, candidate domainsocialcontext.ContextCandidate) error {
	if f.err != nil {
		return f.err
	}
	f.saved = append(f.saved, candidate)
	return nil
}

type candidateIDGeneratorFake struct {
	id    string
	err   error
	calls int
}

func (f *candidateIDGeneratorFake) NewCandidateID(_ context.Context) (string, error) {
	f.calls++
	return f.id, f.err
}

type clockFake struct {
	now   time.Time
	calls int
}

func (f *clockFake) Now() time.Time {
	f.calls++
	return f.now
}

func appOwner(t *testing.T, value string) domainidentity.ID {
	t.Helper()
	id, err := domainidentity.NewID(value)
	if err != nil {
		t.Fatalf("NewID(%q): %v", value, err)
	}
	return id
}

func signal(t *testing.T, id, owner, content string, at time.Time) ActivitySignal {
	t.Helper()
	return ActivitySignal{
		ID:            id,
		OwnerID:       appOwner(t, owner),
		Content:       content,
		Provenance:    "manual",
		OccurredAt:    at,
		ContributedAt: at,
	}
}

func newUseCase(reader ActivitySignalReader, generator ContextGenerator, repository CandidateRepository, idGenerator CandidateIDGenerator, clock Clock) GenerateContextFromActivities {
	return NewGenerateContextFromActivities(reader, generator, repository, idGenerator, clock)
}

func TestGenerateContextSuppressesEmptySourceSetBeforeReadingOrGenerating(t *testing.T) {
	reader := &activitySignalReaderFake{}
	generator := &contextGeneratorFake{}
	repository := &candidateRepositoryFake{}
	ids := &candidateIDGeneratorFake{id: "candidate-1"}
	clock := &clockFake{now: time.Now()}
	useCase := newUseCase(reader, generator, repository, ids, clock)

	result, err := useCase.Execute(context.Background(), GenerateContextFromActivitiesCommand{OwnerID: "alice"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !result.Suppressed || result.Candidate != nil {
		t.Fatalf("result = %#v, want suppressed without candidate", result)
	}
	if reader.calls != 0 || generator.calls != 0 || len(repository.saved) != 0 || ids.calls != 0 || clock.calls != 0 {
		t.Fatalf("empty source caused side effects: reader=%d generator=%d saved=%d ids=%d clock=%d", reader.calls, generator.calls, len(repository.saved), ids.calls, clock.calls)
	}
}

func TestGenerateContextRejectsCrossOwnerSourcesBeforeGenerator(t *testing.T) {
	at := time.Date(2026, time.August, 30, 2, 0, 0, 0, time.UTC)
	reader := &activitySignalReaderFake{items: []ActivitySignal{
		signal(t, "activity-a", "alice", "reading CRDT papers", at),
		signal(t, "activity-b", "bob", "private bob activity", at),
	}}
	generator := &contextGeneratorFake{}
	repository := &candidateRepositoryFake{}
	useCase := newUseCase(reader, generator, repository, &candidateIDGeneratorFake{id: "candidate-1"}, &clockFake{now: at})

	result, err := useCase.Execute(context.Background(), GenerateContextFromActivitiesCommand{OwnerID: "alice", SourceActivityIDs: []string{"activity-a", "activity-b"}})
	if !errors.Is(err, ErrActivityOwnerMismatch) {
		t.Fatalf("Execute() error = %v, want %v", err, ErrActivityOwnerMismatch)
	}
	if result.Candidate != nil || generator.calls != 0 || len(repository.saved) != 0 {
		t.Fatalf("cross-owner source reached generation/persistence: result=%#v generator=%d saved=%d", result, generator.calls, len(repository.saved))
	}
}

func TestGenerateContextSuppressesDuplicateOnlySignalsBeforeGenerator(t *testing.T) {
	at := time.Date(2026, time.August, 30, 2, 0, 0, 0, time.UTC)
	reader := &activitySignalReaderFake{items: []ActivitySignal{
		signal(t, "activity-1", "alice", " Reading   CRDT papers ", at),
		signal(t, "activity-2", "alice", "reading crdt PAPERS", at.Add(time.Minute)),
	}}
	generator := &contextGeneratorFake{}
	repository := &candidateRepositoryFake{}
	useCase := newUseCase(reader, generator, repository, &candidateIDGeneratorFake{id: "candidate-1"}, &clockFake{now: at})

	result, err := useCase.Execute(context.Background(), GenerateContextFromActivitiesCommand{OwnerID: "alice", SourceActivityIDs: []string{"activity-1", "activity-2"}})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !result.Suppressed || result.Candidate != nil {
		t.Fatalf("result = %#v, want duplicate-only suppression", result)
	}
	if generator.calls != 0 || len(repository.saved) != 0 {
		t.Fatalf("duplicate-only signals reached generator/save: generator=%d saved=%d", generator.calls, len(repository.saved))
	}
}

func TestGenerateContextPassesProviderNeutralNormalizedSignalsAndPersistsCandidate(t *testing.T) {
	at := time.Date(2026, time.August, 30, 2, 0, 0, 0, time.UTC)
	generatedAt := at.Add(10 * time.Minute)
	alice := appOwner(t, "alice")
	reader := &activitySignalReaderFake{items: []ActivitySignal{
		signal(t, "activity-1", "alice", "reading CRDT papers", at),
		signal(t, "activity-2", "alice", "comparing consensus protocols", at.Add(time.Minute)),
	}}
	generator := &contextGeneratorFake{output: GeneratedContext{Meaning: "最近持續研究 distributed systems 的 consistency trade-offs"}}
	repository := &candidateRepositoryFake{}
	ids := &candidateIDGeneratorFake{id: "candidate-42"}
	clock := &clockFake{now: generatedAt}
	useCase := newUseCase(reader, generator, repository, ids, clock)

	result, err := useCase.Execute(context.Background(), GenerateContextFromActivitiesCommand{OwnerID: " alice ", SourceActivityIDs: []string{" activity-1 ", "activity-2", "activity-1"}})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.Suppressed || result.Candidate == nil {
		t.Fatalf("result = %#v, want generated candidate", result)
	}
	if generator.calls != 1 {
		t.Fatalf("generator calls = %d, want 1", generator.calls)
	}
	if generator.input.OwnerID != alice {
		t.Fatalf("generator owner = %q, want %q", generator.input.OwnerID, alice)
	}
	if len(generator.input.Signals) != 2 {
		t.Fatalf("generator signals = %#v, want 2 normalized signals", generator.input.Signals)
	}
	if generator.input.Signals[0].ID != "activity-1" || generator.input.Signals[0].Content != "reading CRDT papers" || generator.input.Signals[0].Provenance != "manual" {
		t.Fatalf("generator first signal = %#v, want provider-neutral normalized activity signal", generator.input.Signals[0])
	}
	if len(repository.saved) != 1 {
		t.Fatalf("saved count = %d, want 1", len(repository.saved))
	}
	candidate := repository.saved[0]
	if candidate.OwnerID() != alice || string(candidate.ID()) != "candidate-42" || candidate.Meaning().String() != generator.output.Meaning {
		t.Fatalf("saved candidate = %#v, wrong owner/id/meaning", candidate)
	}
	wantSources := []string{"activity-1", "activity-2"}
	gotSources := candidate.SourceActivityIDs()
	if len(gotSources) != len(wantSources) || gotSources[0] != wantSources[0] || gotSources[1] != wantSources[1] {
		t.Fatalf("candidate sources = %#v, want %#v", gotSources, wantSources)
	}
	if !candidate.GeneratedAt().Equal(generatedAt) || !candidate.IsPrivate() || candidate.IsValidatedSocialContext() {
		t.Fatalf("candidate lifecycle/timestamp invalid: generated=%v private=%v validated=%v", candidate.GeneratedAt(), candidate.IsPrivate(), candidate.IsValidatedSocialContext())
	}
}

func TestGenerateContextRejectsBlankGeneratedMeaningWithoutPersisting(t *testing.T) {
	at := time.Date(2026, time.August, 30, 2, 0, 0, 0, time.UTC)
	reader := &activitySignalReaderFake{items: []ActivitySignal{signal(t, "activity-1", "alice", "reading CRDT papers", at)}}
	generator := &contextGeneratorFake{output: GeneratedContext{Meaning: "  \n "}}
	repository := &candidateRepositoryFake{}
	ids := &candidateIDGeneratorFake{id: "candidate-1"}
	clock := &clockFake{now: at}
	useCase := newUseCase(reader, generator, repository, ids, clock)

	result, err := useCase.Execute(context.Background(), GenerateContextFromActivitiesCommand{OwnerID: "alice", SourceActivityIDs: []string{"activity-1"}})
	if !errors.Is(err, domainsocialcontext.ErrEmptyMeaning) {
		t.Fatalf("Execute() error = %v, want %v", err, domainsocialcontext.ErrEmptyMeaning)
	}
	if result.Candidate != nil || len(repository.saved) != 0 || ids.calls != 0 || clock.calls != 0 {
		t.Fatalf("blank generator output caused candidate side effects: result=%#v saved=%d ids=%d clock=%d", result, len(repository.saved), ids.calls, clock.calls)
	}
}

func TestGenerateContextRejectsSingleSourceRawReplayWithoutPersisting(t *testing.T) {
	at := time.Date(2026, time.August, 30, 2, 0, 0, 0, time.UTC)
	reader := &activitySignalReaderFake{items: []ActivitySignal{signal(t, "activity-1", "alice", "Reading CRDT Papers", at)}}
	generator := &contextGeneratorFake{output: GeneratedContext{Meaning: " reading   crdt papers "}}
	repository := &candidateRepositoryFake{}
	ids := &candidateIDGeneratorFake{id: "candidate-1"}
	clock := &clockFake{now: at}
	useCase := newUseCase(reader, generator, repository, ids, clock)

	result, err := useCase.Execute(context.Background(), GenerateContextFromActivitiesCommand{OwnerID: "alice", SourceActivityIDs: []string{"activity-1"}})
	if !errors.Is(err, ErrPureActivityReplay) {
		t.Fatalf("Execute() error = %v, want %v", err, ErrPureActivityReplay)
	}
	if result.Candidate != nil || len(repository.saved) != 0 || ids.calls != 0 || clock.calls != 0 {
		t.Fatalf("raw replay caused candidate side effects: result=%#v saved=%d ids=%d clock=%d", result, len(repository.saved), ids.calls, clock.calls)
	}
}
