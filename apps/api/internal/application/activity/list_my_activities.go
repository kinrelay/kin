package activity

import (
	"context"
	"sort"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

// ActivityReadModel is the purpose-built private projection returned by Activity queries.
type ActivityReadModel struct {
	ID            string
	OwnerID       domainidentity.ID
	Content       string
	Provenance    string
	OccurredAt    time.Time
	ContributedAt time.Time
}

// ActivityReader provides owner-scoped Activity projections without exposing write aggregates.
type ActivityReader interface {
	ListByOwner(ctx context.Context, ownerID domainidentity.ID) ([]ActivityReadModel, error)
}

// ListMyActivitiesQuery identifies the requester whose private Activities are being listed.
type ListMyActivitiesQuery struct {
	RequesterID string
}

// ListMyActivities lists only the requester's private Activity projections.
type ListMyActivities struct {
	reader ActivityReader
}

// NewListMyActivities constructs the owner-only Activity query.
func NewListMyActivities(reader ActivityReader) ListMyActivities {
	return ListMyActivities{reader: reader}
}

// Execute returns an owner-filtered, deterministic newest-first Activity projection.
func (q ListMyActivities) Execute(ctx context.Context, query ListMyActivitiesQuery) ([]ActivityReadModel, error) {
	requesterID, err := domainidentity.NewID(query.RequesterID)
	if err != nil {
		return nil, err
	}

	items, err := q.reader.ListByOwner(ctx, requesterID)
	if err != nil {
		return nil, err
	}

	result := make([]ActivityReadModel, 0, len(items))
	for _, item := range items {
		if item.OwnerID == requesterID {
			result = append(result, item)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if !result[i].ContributedAt.Equal(result[j].ContributedAt) {
			return result[i].ContributedAt.After(result[j].ContributedAt)
		}
		if !result[i].OccurredAt.Equal(result[j].OccurredAt) {
			return result[i].OccurredAt.After(result[j].OccurredAt)
		}
		return result[i].ID < result[j].ID
	})

	return result, nil
}
