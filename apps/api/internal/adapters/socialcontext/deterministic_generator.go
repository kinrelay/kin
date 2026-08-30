package socialcontext

import (
	"context"
	"fmt"

	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

// DeterministicGenerator is the provider-free MVP adapter used to make the derivation path executable and testable.
// It deliberately derives a higher-level statement from significance-approved signals instead of replaying raw Activity content.
type DeterministicGenerator struct{}

func NewDeterministicGenerator() DeterministicGenerator {
	return DeterministicGenerator{}
}

func (DeterministicGenerator) Generate(_ context.Context, input applicationsocialcontext.ContextGenerationInput) (applicationsocialcontext.GeneratedContext, error) {
	provenance := make([]string, 0, len(input.Activities))
	for _, activity := range input.Activities {
		provenance = append(provenance, activity.ID)
	}

	meaning := "近期有一個持續投入、值得之後再聊聊的主題"
	if len(input.Activities) > 1 {
		meaning = fmt.Sprintf("近期有 %d 個彼此相關、持續投入且值得之後再聊聊的主題訊號", len(input.Activities))
	}

	return applicationsocialcontext.GeneratedContext{
		Meaning:    meaning,
		Provenance: provenance,
	}, nil
}
