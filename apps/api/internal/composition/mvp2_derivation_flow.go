package composition

import (
	activityadapter "github.com/kinrelay/kin/apps/api/internal/adapters/activity"
	socialcontextadapter "github.com/kinrelay/kin/apps/api/internal/adapters/socialcontext"
	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

// MVP2DerivationFlow wires the provider-free in-memory adapters needed to execute the MVP 2 core vertical path.
type MVP2DerivationFlow struct {
	Activities *activityadapter.MemoryRepository
	Derive     applicationsocialcontext.DeriveContextCandidate
	List       applicationssocialcontext.ListMySocialContexts
}

// NewMVP2DerivationFlow builds the smallest executable MVP 2 composition without delivery or external-provider infrastructure.
func NewMVP2DerivationFlow() MVP2DerivationFlow {
	activities := activityadapter.NewMemoryRepository()
	activityReader := activityadapter.NewMemoryReadRepository(activities)
	contexts := socialcontextadapter.NewMemoryRepository()
	contextReader := socialcontextadapter.NewMemoryReadRepository(contexts)
	generator := socialcontextadapter.NewDeterministicGenerator()

	return MVP2DerivationFlow{
		Activities: activities,
		Derive:     applicationsocialcontext.NewDeriveContextCandidate(activityReader, generator, contexts),
		List:       applicationssocialcontext.NewListMySocialContexts(contextReader),
	}
}
