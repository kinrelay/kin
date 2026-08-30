package socialcontext

import (
	"context"
	"sort"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

// SocialContextReadModel is the purpose-built owner-private projection returned by Social Context queries.
type SocialContextReadModel struct {
	ID         string
	OwnerID    domainidentity.ID
	Meaning    string
	Provenance []string
	PromotedAt time.Time
}

// SocialContextReader provides owner-scoped projections without exposing write-side SocialContext aggregates.
type SocialContextReader interface {
	ListByOwner(ctx context.Context, ownerID domainidentity.ID) ([]SocialContextReadModel, error)
}

// ListMySocialContextsQuery identifies the owner whose private derived contexts are being listed.
type ListMySocialContextsQuery struct {
	RequesterID string
}

// ListMySocialContexts lists only the requester's validated private Social Context projections.
type ListMySocialContexts struct {
	reader SocialContextReader
}

// NewListMySocialContexts constructs the owner-only Social Context query.
func NewListMySocialContexts(reader SocialContextReader) ListMySocialContexts {
	return ListMySocialContexts{reader: reader}
}

// Execute returns an owner-filtered deterministic newest-first projection.
func (q ListMySocialContexts) Execute(ctx context.Context, query ListMySocialContextsQuery) ([]SocialContextReadModel, error) {
	requesterID, err := domainidentity.NewID(query.RequesterID)
	if err != nil {
		return nil, err
	}

	items, err := q.reader.ListByOwner(ctx, requesterID)
	if err != nil {
		return nil, err
	}

	result := make([]SocialContextReadModel, 0, len(items))
	for _, item := range items {
		if item.OwnerID != requesterID {
			continue
		}
		copyItem := item
		copyItem.Provenance = append([]string(nil), item.Provenance...)
		result = append(result, copyItem)
	}

	sort.Slice(result, func(i, j int) bool {
		if !result[i].PromotedAt.Equal(result[j].PromotedAt) {
			return result[i].PromotedAt.After(result[j].PromotedAt)
		}
		return result[i].ID < result[j].ID
	})

	return result, nil
}
