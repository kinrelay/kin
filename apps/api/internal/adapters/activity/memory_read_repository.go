package activity

import (
	"context"

	applicationactivity "github.com/kinrelay/kin/apps/api/internal/application/activity"
	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
	domainactivity "github.com/kinrelay/kin/apps/api/internal/domain/activity"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

// MemoryReadRepository projects owner-scoped Activity read models from in-memory write storage.
type MemoryReadRepository struct {
	source *MemoryRepository
}

// NewMemoryReadRepository constructs a dedicated Activity read adapter over the in-memory source.
func NewMemoryReadRepository(source *MemoryRepository) *MemoryReadRepository {
	return &MemoryReadRepository{source: source}
}

// ListByOwner returns purpose-built projections for exactly one Activity owner.
func (r *MemoryReadRepository) ListByOwner(_ context.Context, ownerID domainidentity.ID) ([]applicationactivity.ActivityReadModel, error) {
	r.source.mu.RLock()
	defer r.source.mu.RUnlock()

	result := make([]applicationactivity.ActivityReadModel, 0)
	for _, value := range r.source.activities {
		if value.OwnerID() != ownerID {
			continue
		}
		result = append(result, applicationactivity.ActivityReadModel{
			ID:            string(value.ID()),
			OwnerID:       value.OwnerID(),
			Content:       value.Content().String(),
			Provenance:    string(value.Provenance()),
			OccurredAt:    value.OccurredAt(),
			ContributedAt: value.ContributedAt(),
		})
	}

	return result, nil
}

// ListOwnerPrivateNormalized exposes only explicitly requested, owner-private normalized Activities to context derivation.
func (r *MemoryReadRepository) ListOwnerPrivateNormalized(_ context.Context, ownerID domainidentity.ID, activityIDs []string) ([]applicationsocialcontext.ActivityForContext, error) {
	r.source.mu.RLock()
	defer r.source.mu.RUnlock()

	result := make([]applicationsocialcontext.ActivityForContext, 0, len(activityIDs))
	seen := make(map[string]struct{}, len(activityIDs))
	for _, id := range activityIDs {
		if _, duplicate := seen[id]; duplicate {
			continue
		}
		seen[id] = struct{}{}

		value, ok := r.source.activities[domainactivity.ID(id)]
		if !ok || value.OwnerID() != ownerID || !value.IsPrivate() {
			continue
		}
		result = append(result, applicationsocialcontext.ActivityForContext{
			ID:      id,
			OwnerID: value.OwnerID(),
			Content: value.Content().String(),
		})
	}

	return result, nil
}
