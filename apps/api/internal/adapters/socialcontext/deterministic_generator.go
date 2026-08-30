package socialcontext

import (
	"context"
	"sort"
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
	seenTopics := make(map[string]struct{}, len(input.Activities))
	for _, activity := range input.Activities {
		provenance = append(provenance, activity.ID)
		topic, ok := summarizeSignal(activity.Content)
		if !ok {
			continue
		}
		if _, seen := seenTopics[topic]; seen {
			continue
		}
		seenTopics[topic] = struct{}{}
		topics = append(topics, topic)
	}
	if len(topics) == 0 {
		return applicationsocialcontext.GeneratedContext{Provenance: provenance}, nil
	}
	sort.Strings(topics)

	return applicationsocialcontext.GeneratedContext{
		Meaning:    "近期關注" + strings.Join(topics, "；"),
		Provenance: provenance,
	}, nil
}

func summarizeSignal(content string) (string, bool) {
	normalized := strings.TrimSpace(content)
	switch {
	case strings.Contains(normalized, "分散式系統"), strings.Contains(normalized, "一致性模型"):
		return "分散式系統的一致性模型、可靠性與工程取捨", true
	case strings.Contains(normalized, "馬拉松") && (strings.Contains(normalized, "準備") || strings.Contains(normalized, "訓練")):
		return "耐力運動與長距離訓練", true
	default:
		return "", false
	}
}
