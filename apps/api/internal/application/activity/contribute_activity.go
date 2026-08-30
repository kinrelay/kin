package activity

import (
	"context"
	"errors"
	"time"

	domainactivity "github.com/kinrelay/kin/apps/api/internal/domain/activity"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

// ErrIdentityNotFound indicates that the contributor is not a known Kin identity.
var ErrIdentityNotFound = errors.New("identity not found")

// IdentityDirectory answers whether an Identity may own an Activity.
type IdentityDirectory interface {
	Exists(ctx context.Context, id domainidentity.ID) (bool, error)
}

// Repository persists normalized Activities.
type Repository interface {
	Save(ctx context.Context, value domainactivity.Activity) error
}

// IDGenerator provides stable Activity IDs without coupling the use case to randomness or storage.
type IDGenerator interface {
	NewActivityID(ctx context.Context) (string, error)
}

// Clock provides deterministic contribution time to the use case.
type Clock interface {
	Now() time.Time
}

// ContributeActivityCommand is raw explicit user input for one manual contribution.
// It is intentionally separate from the normalized Activity domain model.
type ContributeActivityCommand struct {
	ContributorID string
	Content       string
	OccurredAt    time.Time
}

// ContributeActivity normalizes and persists one explicit private Activity contribution.
type ContributeActivity struct {
	identities IdentityDirectory
	repository Repository
	ids        IDGenerator
	clock      Clock
}

// NewContributeActivity constructs the contribution use case from inner-facing ports.
func NewContributeActivity(identities IdentityDirectory, repository Repository, ids IDGenerator, clock Clock) ContributeActivity {
	return ContributeActivity{
		identities: identities,
		repository: repository,
		ids:        ids,
		clock:      clock,
	}
}

// Execute validates raw contribution input, verifies ownership, normalizes the Activity, and persists it.
func (uc ContributeActivity) Execute(ctx context.Context, command ContributeActivityCommand) (domainactivity.Activity, error) {
	contributorID, err := domainidentity.NewID(command.ContributorID)
	if err != nil {
		return domainactivity.Activity{}, err
	}
	content, err := domainactivity.NewContent(command.Content)
	if err != nil {
		return domainactivity.Activity{}, err
	}

	exists, err := uc.identities.Exists(ctx, contributorID)
	if err != nil {
		return domainactivity.Activity{}, err
	}
	if !exists {
		return domainactivity.Activity{}, ErrIdentityNotFound
	}

	id, err := uc.ids.NewActivityID(ctx)
	if err != nil {
		return domainactivity.Activity{}, err
	}
	created, err := domainactivity.NewManual(id, contributorID, content, command.OccurredAt, uc.clock.Now())
	if err != nil {
		return domainactivity.Activity{}, err
	}
	if err := uc.repository.Save(ctx, created); err != nil {
		return domainactivity.Activity{}, err
	}

	return created, nil
}
