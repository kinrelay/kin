package activity

import (
	"context"
	"sort"

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

// ListOwnerPrivateNormalized exposes only explicitly requested, owner-private normalized Activities to context derivation in occurrence order.
func (r *MemoryReadRepository) ListOwnerPrivateNormalized(_ context.Context, ownerID domainidentity.ID, activityIDs []string) ([]applicationsocialcontext.ActivityForContext, error) {
	r.source.mu.RLock()
	defer r.source.mu.RUnlock()

	selected := make([]domainactivity.Activity, 0, len(activityIDs))
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
		selected = append(selected, value)
	}

	sort.SliceStable(selected, func(i, j int) bool {
		left, right := selected[i], selected[j]
		if left.OccurredAt().Equal(right.OccurredAt()) {
			return string(left.ID()) < string(right.ID())
		}
		return left.OccurredAt().Before(right.OccurredAt())
	})

	result := make([]applicationsocialcontext.ActivityForContext, 0, len(selected))
	for _, value := range selected {
		result = append(result, applicationsocialcontext.ActivityForContext{
			ID:         string(value.ID()),
			OwnerID:    value.OwnerID(),
			Content:    value.Content().String(),
			OccurredAt: value.OccurredAt(),
		})
	}
	return result, nil
}
