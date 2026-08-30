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
		topic, ok := summarizeSignal(activity.Content)
		if !ok {
			continue
		}
		provenance = append(provenance, activity.ID)
		if _, seen := seenTopics[topic]; seen {
			continue
		}
		seenTopics[topic] = struct{}{}
		topics = append(topics, topic)
	}
	if len(topics) == 0 {
		return applicationsocialcontext.GeneratedContext{}, nil
	}
	sort.Strings(topics)

	return applicationsocialcontext.GeneratedContext{
		Meaning:    "近期關注" + strings.Join(topics, "；"),
		Provenance: provenance,
	}, nil
}

func summarizeSignal(content string) (string, bool) {
	normalized := strings.TrimSpace(content)
	if containsAny(normalized,
		"不研究", "不再研究", "停止研究", "沒有研究",
		"不準備", "沒有準備", "不訓練", "沒有訓練",
	) {
		return "", false
	}

	switch {
	case containsAny(normalized, "深入研究分散式系統", "研究分散式系統", "比較不同一致性模型", "研究一致性模型"):
		return "分散式系統的一致性模型、可靠性與工程取捨", true
	case strings.Contains(normalized, "馬拉松") && containsAny(normalized, "準備", "訓練"):
		return "耐力運動與長距離訓練", true
	default:
		return "", false
	}
}

func containsAny(content string, patterns ...string) bool {
	for _, pattern := range patterns {
		if strings.Contains(content, pattern) {
			return true
		}
	}
	return false
}
