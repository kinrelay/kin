package socialcontext

import (
	"context"
	"strings"

	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

// DeterministicGenerator is the provider-free MVP adapter used to make the derivation path executable and testable.
// It derives a stable higher-level statement from significance-approved signals without replaying raw Activity content verbatim.
type DeterministicGenerator struct{}

func NewDeterministicGenerator() DeterministicGenerator {
	return DeterministicGenerator{}
}

func (DeterministicGenerator) Generate(_ context.Context, input applicationsocialcontext.ContextGenerationInput) (applicationsocialcontext.GeneratedContext, error) {
	provenance := make([]string, 0, len(input.Activities))
	topics := make([]string, 0, len(input.Activities))
	for _, activity := range input.Activities {
		provenance = append(provenance, activity.ID)
		topic, ok := summarizeSignal(activity.Content)
		if !ok {
			return applicationsocialcontext.GeneratedContext{Provenance: provenance}, nil
		}
		topics = append(topics, topic)
	}

	return applicationsocialcontext.GeneratedContext{
		Meaning:    "近期持續關注與投入：" + strings.Join(topics, "；"),
		Provenance: provenance,
	}, nil
}

func summarizeSignal(content string) (string, bool) {
	summary := strings.TrimSpace(content)
	for _, prefix := range []string{
		"最近開始",
		"最近",
		"持續",
		"開始",
		"深入研究",
		"研究",
		"準備",
	} {
		if strings.HasPrefix(summary, prefix) {
			topic := strings.TrimSpace(strings.TrimPrefix(summary, prefix))
			if topic == "" || topic == summary {
				return "", false
			}
			return topic, true
		}
	}
	return "", false
}
