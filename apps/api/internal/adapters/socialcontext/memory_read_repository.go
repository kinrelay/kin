package socialcontext

import (
	"context"

	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

// MemoryReadRepository projects owner-scoped Social Context read models from in-memory write storage.
type MemoryReadRepository struct {
	source *MemoryRepository
}

// NewMemoryReadRepository constructs a dedicated Social Context read adapter over the in-memory source.
func NewMemoryReadRepository(source *MemoryRepository) *MemoryReadRepository {
	return &MemoryReadRepository{source: source}
}

// ListByOwner returns purpose-built projections for exactly one Social Context owner.
func (r *MemoryReadRepository) ListByOwner(_ context.Context, ownerID domainidentity.ID) ([]applicationsocialcontext.SocialContextReadModel, error) {
	r.source.mu.RLock()
	defer r.source.mu.RUnlock()

	result := make([]applicationsocialcontext.SocialContextReadModel, 0)
	for _, entry := range r.source.entries {
		if entry.ownerID != ownerID {
			continue
		}
		result = append(result, applicationsocialcontext.SocialContextReadModel{
			ID:         entry.id,
			OwnerID:    entry.ownerID,
			Meaning:    entry.context.Meaning(),
			Provenance: entry.context.Provenance(),
			PromotedAt: entry.promotedAt,
		})
	}
	return result, nil
}
