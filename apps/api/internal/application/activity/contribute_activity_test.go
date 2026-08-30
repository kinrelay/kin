package activity

import (
	"context"
	"errors"
	"testing"
	"time"

	domainactivity "github.com/kinrelay/kin/apps/api/internal/domain/activity"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

type activityIdentityDirectoryFake struct {
	existing map[domainidentity.ID]bool
	err      error
}

func (f activityIdentityDirectoryFake) Exists(_ context.Context, id domainidentity.ID) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return f.existing[id], nil
}

type activityRepositoryFake struct {
	saved []domainactivity.Activity
	err   error
}

func (f *activityRepositoryFake) Save(_ context.Context, value domainactivity.Activity) error {
	if f.err != nil {
		return f.err
	}
	f.saved = append(f.saved, value)
	return nil
}

type activityIDGeneratorFake struct {
	id  string
	err error
}

func (f activityIDGeneratorFake) NewActivityID(_ context.Context) (string, error) {
	return f.id, f.err
}

type activityClockFake struct {
	now time.Time
}

func (f activityClockFake) Now() time.Time {
	return f.now
}

func applicationIdentityID(t *testing.T, value string) domainidentity.ID {
	t.Helper()
	id, err := domainidentity.NewID(value)
	if err != nil {
		t.Fatalf("NewID(%q): %v", value, err)
	}
	return id
}

func TestContributeActivityNormalizesAndPersistsPrivateManualActivity(t *testing.T) {
	ctx := context.Background()
	owner := applicationIdentityID(t, "alice")
	repository := &activityRepositoryFake{}
	occurredAt := time.Date(2026, time.August, 29, 9, 30, 0, 0, time.UTC)
	contributedAt := time.Date(2026, time.August, 29, 10, 0, 0, 0, time.UTC)
	useCase := NewContributeActivity(
		activityIdentityDirectoryFake{existing: map[domainidentity.ID]bool{owner: true}},
		repository,
		activityIDGeneratorFake{id: "activity-1"},
		activityClockFake{now: contributedAt},
	)

	created, err := useCase.Execute(ctx, ContributeActivityCommand{
		ContributorID: " alice ",
		Content:       "  researching local-first apps  ",
		OccurredAt:    occurredAt,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(repository.saved) != 1 {
		t.Fatalf("saved count = %d, want 1", len(repository.saved))
	}
	if created.ID() != domainactivity.ID("activity-1") {
		t.Fatalf("ID() = %q, want activity-1", created.ID())
	}
	if created.OwnerID() != owner {
		t.Fatalf("OwnerID() = %q, want %q", created.OwnerID(), owner)
	}
	if created.Content().String() != "researching local-first apps" {
		t.Fatalf("Content() = %q, want normalized content", created.Content().String())
	}
	if created.Provenance() != domainactivity.ProvenanceManual {
		t.Fatalf("Provenance() = %q, want manual", created.Provenance())
	}
	if !created.IsPrivate() {
		t.Fatal("contributed Activity must remain private")
	}
	if !created.OccurredAt().Equal(occurredAt) || !created.ContributedAt().Equal(contributedAt) {
		t.Fatalf("timestamps = (%v, %v), want (%v, %v)", created.OccurredAt(), created.ContributedAt(), occurredAt, contributedAt)
	}
}

func TestContributeActivityRejectsMissingContributorWithoutSaving(t *testing.T) {
	repository := &activityRepositoryFake{}
	useCase := NewContributeActivity(
		activityIdentityDirectoryFake{existing: map[domainidentity.ID]bool{}},
		repository,
		activityIDGeneratorFake{id: "activity-1"},
		activityClockFake{now: time.Date(2026, time.August, 29, 10, 0, 0, 0, time.UTC)},
	)

	_, err := useCase.Execute(context.Background(), ContributeActivityCommand{
		ContributorID: "missing",
		Content:       "reading about CRDTs",
		OccurredAt:    time.Date(2026, time.August, 29, 9, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("Execute() error = %v, want %v", err, ErrIdentityNotFound)
	}
	if len(repository.saved) != 0 {
		t.Fatalf("saved count = %d, want 0", len(repository.saved))
	}
}

func TestContributeActivityRejectsBlankContentBeforePersistence(t *testing.T) {
	owner := applicationIdentityID(t, "alice")
	repository := &activityRepositoryFake{}
	useCase := NewContributeActivity(
		activityIdentityDirectoryFake{existing: map[domainidentity.ID]bool{owner: true}},
		repository,
		activityIDGeneratorFake{id: "activity-1"},
		activityClockFake{now: time.Date(2026, time.August, 29, 10, 0, 0, 0, time.UTC)},
	)

	_, err := useCase.Execute(context.Background(), ContributeActivityCommand{
		ContributorID: "alice",
		Content:       "   ",
		OccurredAt:    time.Date(2026, time.August, 29, 9, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, domainactivity.ErrEmptyContent) {
		t.Fatalf("Execute() error = %v, want %v", err, domainactivity.ErrEmptyContent)
	}
	if len(repository.saved) != 0 {
		t.Fatalf("saved count = %d, want 0", len(repository.saved))
	}
}

func TestContributeActivityPropagatesPortFailuresWithoutSavingInvalidState(t *testing.T) {
	owner := applicationIdentityID(t, "alice")
	now := time.Date(2026, time.August, 29, 10, 0, 0, 0, time.UTC)
	portErr := errors.New("port failed")

	tests := []struct {
		name       string
		directory  activityIdentityDirectoryFake
		generator  activityIDGeneratorFake
		repository *activityRepositoryFake
	}{
		{
			name:       "identity directory",
			directory:  activityIdentityDirectoryFake{err: portErr},
			generator:  activityIDGeneratorFake{id: "activity-1"},
			repository: &activityRepositoryFake{},
		},
		{
			name:       "id generator",
			directory:  activityIdentityDirectoryFake{existing: map[domainidentity.ID]bool{owner: true}},
			generator:  activityIDGeneratorFake{err: portErr},
			repository: &activityRepositoryFake{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := NewContributeActivity(tt.directory, tt.repository, tt.generator, activityClockFake{now: now})
			_, err := useCase.Execute(context.Background(), ContributeActivityCommand{
				ContributorID: "alice",
				Content:       "reading about CRDTs",
				OccurredAt:    now,
			})
			if !errors.Is(err, portErr) {
				t.Fatalf("Execute() error = %v, want %v", err, portErr)
			}
			if len(tt.repository.saved) != 0 {
				t.Fatalf("saved count = %d, want 0", len(tt.repository.saved))
			}
		})
	}
}
