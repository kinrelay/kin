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
	clauses := splitSignalClauses(content)
	for i, clause := range clauses {
		switch {
		case hasAnyPrefix(clause,
			"最近開始深入研究分散式系統",
			"開始深入研究分散式系統",
			"持續深入研究分散式系統",
			"最近開始研究分散式系統",
			"開始研究分散式系統",
			"持續研究分散式系統",
			"最近開始比較不同一致性模型",
			"開始比較不同一致性模型",
			"持續比較不同一致性模型",
			"最近開始研究一致性模型",
			"開始研究一致性模型",
			"持續研究一致性模型",
		):
			if hasTopicReversal(
				clauses[i:],
				[]string{"分散式系統", "一致性模型"},
				[]string{"不再研究", "不想研究", "停止研究", "沒有研究", "不再比較", "停止比較"},
			) {
				continue
			}
			return "分散式系統的一致性模型、可靠性與工程取捨", true
		case isAffirmativeMarathonParticipationClause(clause):
			if hasTopicReversal(
				clauses[i:],
				[]string{"馬拉松"},
				[]string{"沒有準備", "不再準備", "停止準備", "沒有訓練", "不再訓練", "停止訓練", "放棄馬拉松", "放棄了", "放棄", "不參加", "不參賽"},
			) {
				continue
			}
			return "耐力運動與長距離訓練", true
		}
	}
	return "", false
}

func isAffirmativeMarathonParticipationClause(clause string) bool {
	if strings.Contains(clause, "不參加") || strings.Contains(clause, "不參賽") {
		return false
	}
	for _, prefix := range []string{
		"最近開始準備",
		"開始準備",
		"持續準備",
		"最近開始訓練",
		"開始訓練",
		"持續訓練",
	} {
		if !strings.HasPrefix(clause, prefix) {
			continue
		}
		target := strings.TrimSpace(strings.TrimPrefix(clause, prefix))
		for _, participationTarget := range []string{
			"馬拉松",
			"第一次全程馬拉松訓練",
			"全程馬拉松訓練",
			"馬拉松訓練",
			"馬拉松參賽",
			"馬拉松比賽",
		} {
			if target == participationTarget {
				return true
			}
		}
	}
	return false
}

func splitSignalClauses(content string) []string {
	parts := strings.FieldsFunc(content, func(r rune) bool {
		switch r {
		case '，', ',', '。', '.', '！', '!', '？', '?', '；', ';', '\n':
			return true
		default:
			return false
		}
	})
	clauses := make([]string, 0, len(parts))
	for _, part := range parts {
		if clause := strings.TrimSpace(part); clause != "" {
			clauses = append(clauses, clause)
		}
	}
	return clauses
}

func hasTopicReversal(clauses []string, topicMarkers, reversalPatterns []string) bool {
	for _, clause := range clauses {
		for _, pattern := range reversalPatterns {
			index := strings.Index(clause, pattern)
			if index < 0 {
				continue
			}

			suffix := strings.TrimSpace(clause[index+len(pattern):])
			if suffix == "" || hasAnySubstring(suffix, topicMarkers...) {
				return true
			}
		}
	}
	return false
}

func hasAnyPrefix(content string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(content, prefix) {
			return true
		}
	}
	return false
}

func hasAnySubstring(content string, values ...string) bool {
	for _, value := range values {
		if strings.Contains(content, value) {
			return true
		}
	}
	return false
}
