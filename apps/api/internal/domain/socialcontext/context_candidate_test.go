package socialcontext

import (
	"errors"
	"testing"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func candidateOwner(t *testing.T, value string) domainidentity.ID {
	t.Helper()
	id, err := domainidentity.NewID(value)
	if err != nil {
		t.Fatalf("NewID(%q): %v", value, err)
	}
	return id
}

func TestNewMeaningNormalizesAndRejectsBlankMeaning(t *testing.T) {
	meaning, err := NewMeaning("  最近持續研究 agent workflow  ")
	if err != nil {
		t.Fatalf("NewMeaning() error = %v", err)
	}
	if got, want := meaning.String(), "最近持續研究 agent workflow"; got != want {
		t.Fatalf("meaning = %q, want %q", got, want)
	}

	_, err = NewMeaning(" \n\t ")
	if !errors.Is(err, ErrEmptyMeaning) {
		t.Fatalf("blank meaning error = %v, want %v", err, ErrEmptyMeaning)
	}
}

func TestMeaningDetectsWhitespaceAndCaseNormalizedRawReplay(t *testing.T) {
	meaning, err := NewMeaning("  Reading   CRDT Papers ")
	if err != nil {
		t.Fatalf("NewMeaning(): %v", err)
	}
	if !meaning.IsPureReplayOf("reading crdt papers") {
		t.Fatal("IsPureReplayOf() = false, want true for case/whitespace-only rewrite")
	}
	if meaning.IsPureReplayOf("recently exploring distributed systems concepts") {
		t.Fatal("IsPureReplayOf() = true for genuinely derived meaning")
	}
}

func TestNewContextCandidateCapturesPrivateUnvalidatedCandidateStateAndProvenance(t *testing.T) {
	owner := candidateOwner(t, "alice")
	meaning, err := NewMeaning("最近持續研究 distributed systems 的 consistency trade-offs")
	if err != nil {
		t.Fatalf("NewMeaning(): %v", err)
	}
	generatedAt := time.Date(2026, time.August, 30, 1, 30, 0, 0, time.UTC)

	candidate, err := NewContextCandidate(
		" candidate-1 ",
		owner,
		meaning,
		[]string{" activity-2 ", "activity-1", "activity-2"},
		generatedAt,
	)
	if err != nil {
		t.Fatalf("NewContextCandidate() error = %v", err)
	}
	if got, want := string(candidate.ID()), "candidate-1"; got != want {
		t.Fatalf("ID() = %q, want %q", got, want)
	}
	if candidate.OwnerID() != owner {
		t.Fatalf("OwnerID() = %q, want %q", candidate.OwnerID(), owner)
	}
	if candidate.Meaning() != meaning {
		t.Fatalf("Meaning() = %#v, want %#v", candidate.Meaning(), meaning)
	}
	if got, want := candidate.State(), StateCandidate; got != want {
		t.Fatalf("State() = %q, want %q", got, want)
	}
	if !candidate.IsPrivate() {
		t.Fatal("IsPrivate() = false, want owner-private candidate")
	}
	if candidate.IsValidatedSocialContext() {
		t.Fatal("IsValidatedSocialContext() = true, candidate must not be promoted implicitly")
	}
	if !candidate.GeneratedAt().Equal(generatedAt) {
		t.Fatalf("GeneratedAt() = %v, want %v", candidate.GeneratedAt(), generatedAt)
	}
	gotSources := candidate.SourceActivityIDs()
	wantSources := []string{"activity-2", "activity-1"}
	if len(gotSources) != len(wantSources) {
		t.Fatalf("SourceActivityIDs() = %#v, want %#v", gotSources, wantSources)
	}
	for i := range wantSources {
		if gotSources[i] != wantSources[i] {
			t.Fatalf("SourceActivityIDs()[%d] = %q, want %q", i, gotSources[i], wantSources[i])
		}
	}
	gotSources[0] = "mutated"
	if candidate.SourceActivityIDs()[0] != "activity-2" {
		t.Fatal("SourceActivityIDs() exposes mutable internal provenance")
	}
}

func TestNewContextCandidateRejectsInvalidRequiredFields(t *testing.T) {
	owner := candidateOwner(t, "alice")
	meaning, err := NewMeaning("derived meaning")
	if err != nil {
		t.Fatalf("NewMeaning(): %v", err)
	}
	now := time.Date(2026, time.August, 30, 1, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		id        string
		owner     domainidentity.ID
		meaning   Meaning
		sources   []string
		generated time.Time
		wantErr   error
	}{
		{name: "empty id", id: "  ", owner: owner, meaning: meaning, sources: []string{"activity-1"}, generated: now, wantErr: ErrInvalidCandidateID},
		{name: "invalid owner", id: "candidate-1", owner: " ", meaning: meaning, sources: []string{"activity-1"}, generated: now, wantErr: domainidentity.ErrInvalidID},
		{name: "empty source list", id: "candidate-1", owner: owner, meaning: meaning, generated: now, wantErr: ErrMissingSourceActivity},
		{name: "blank source id", id: "candidate-1", owner: owner, meaning: meaning, sources: []string{"  "}, generated: now, wantErr: ErrMissingSourceActivity},
		{name: "zero generated time", id: "candidate-1", owner: owner, meaning: meaning, sources: []string{"activity-1"}, wantErr: ErrInvalidCandidateTimestamp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewContextCandidate(tt.id, tt.owner, tt.meaning, tt.sources, tt.generated)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewContextCandidate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
