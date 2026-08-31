package socialcontext

import (
	"context"
	"sort"
	"strings"

	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

const (
	distributedSystemsTopic = "分散式系統的一致性模型、可靠性與工程取捨"
	marathonTopic           = "耐力運動與長距離訓練"
)

var (
	distributedSystemsMarkers = []string{"分散式系統", "一致性模型"}
	distributedSystemsReversals = []string{
		"不再深入研究", "不想深入研究", "停止深入研究", "沒有深入研究", "放棄深入研究",
		"不再研究", "不想研究", "停止研究", "沒有研究", "放棄研究",
		"不再比較", "停止比較", "放棄了", "放棄",
	}
	marathonMarkers   = []string{"馬拉松"}
	marathonReversals = []string{"沒有準備", "不再準備", "停止準備", "沒有訓練", "不再訓練", "停止訓練", "放棄馬拉松", "放棄了", "放棄", "不參加", "不參賽"}
)

// DeterministicGenerator is the provider-free MVP adapter used to make the derivation path executable and testable.
// It derives a stable higher-level statement from significance-approved signals without replaying raw Activity content verbatim.
type DeterministicGenerator struct{}

type recognizedSignal struct {
	activityID string
	topic      string
}

func NewDeterministicGenerator() DeterministicGenerator {
	return DeterministicGenerator{}
}

func (DeterministicGenerator) Generate(_ context.Context, input applicationsocialcontext.ContextGenerationInput) (applicationsocialcontext.GeneratedContext, error) {
	recognized := make([]recognizedSignal, 0, len(input.Activities))
	for _, activity := range input.Activities {
		// Input order is chronological (oldest to newest). A reversal removes only
		// previously recognized state in this requested derivation batch, allowing
		// a genuinely later affirmative signal to establish the topic again.
		for topic := range reversedTopics(activity.Content) {
			recognized, _ = removeRecognizedTopic(recognized, topic)
		}

		seenInActivity := make(map[string]struct{}, 2)
		for _, topic := range summarizeSignals(activity.Content) {
			if _, seen := seenInActivity[topic]; seen {
				continue
			}
			seenInActivity[topic] = struct{}{}
			recognized = append(recognized, recognizedSignal{activityID: activity.ID, topic: topic})
		}
	}
	if len(recognized) == 0 {
		return applicationsocialcontext.GeneratedContext{}, nil
	}

	provenance := make([]string, 0, len(recognized))
	topics := make([]string, 0, len(recognized))
	seenActivityIDs := make(map[string]struct{}, len(recognized))
	seenTopics := make(map[string]struct{}, len(recognized))
	for _, signal := range recognized {
		if _, seen := seenActivityIDs[signal.activityID]; !seen {
			seenActivityIDs[signal.activityID] = struct{}{}
			provenance = append(provenance, signal.activityID)
		}
		if _, seen := seenTopics[signal.topic]; seen {
			continue
		}
		seenTopics[signal.topic] = struct{}{}
		topics = append(topics, signal.topic)
	}
	sort.Strings(topics)

	return applicationsocialcontext.GeneratedContext{
		Meaning:    "近期關注" + strings.Join(topics, "；"),
		Provenance: provenance,
	}, nil
}

func summarizeSignal(content string) (string, bool) {
	topics := summarizeSignals(content)
	if len(topics) == 0 {
		return "", false
	}
	return topics[0], true
}

func summarizeSignals(content string) []string {
	clauses := splitSignalClauses(content)
	topics := make([]string, 0, 2)
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
			if hasTopicReversal(clauses[i:], distributedSystemsMarkers, distributedSystemsReversals) {
				continue
			}
			topics = append(topics, distributedSystemsTopic)
		case isAffirmativeMarathonParticipationClause(clause):
			if hasTopicReversal(clauses[i:], marathonMarkers, marathonReversals) {
				continue
			}
			topics = append(topics, marathonTopic)
		}
	}
	return topics
}

func reversedTopics(content string) map[string]struct{} {
	clauses := splitSignalClauses(content)
	reversed := make(map[string]struct{}, 2)
	if hasTopicReversal(clauses, distributedSystemsMarkers, distributedSystemsReversals) {
		reversed[distributedSystemsTopic] = struct{}{}
	}
	if hasTopicReversal(clauses, marathonMarkers, marathonReversals) {
		reversed[marathonTopic] = struct{}{}
	}
	return reversed
}

func removeRecognizedTopic(recognized []recognizedSignal, topic string) ([]recognizedSignal, []string) {
	kept := make([]recognizedSignal, 0, len(recognized))
	retired := make([]string, 0)
	for _, signal := range recognized {
		if signal.topic == topic {
			retired = append(retired, signal.activityID)
			continue
		}
		kept = append(kept, signal)
	}
	return kept, retired
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
		for _, segment := range splitCompoundActions(strings.TrimSpace(part)) {
			if segment != "" {
				clauses = append(clauses, segment)
			}
		}
	}
	return clauses
}

func splitCompoundActions(clause string) []string {
	if clause == "" {
		return nil
	}
	boundaryIndex, boundaryLength := nextActionBoundary(clause)
	if boundaryIndex < 0 {
		return []string{clause}
	}
	left := strings.TrimSpace(clause[:boundaryIndex])
	right := strings.TrimSpace(clause[boundaryIndex+boundaryLength:])
	result := make([]string, 0, 1+len(splitCompoundActions(right)))
	if left != "" {
		result = append(result, left)
	}
	return append(result, splitCompoundActions(right)...)
}

func nextActionBoundary(content string) (int, int) {
	bestIndex := -1
	bestLength := 0
	for _, boundary := range []string{"以及", "並且", "同時", "並", "且", "而", "但"} {
		searchFrom := 0
		for searchFrom < len(content) {
			relative := strings.Index(content[searchFrom:], boundary)
			if relative < 0 {
				break
			}
			index := searchFrom + relative
			remainder := strings.TrimSpace(content[index+len(boundary):])
			if startsSupportedAction(remainder) && (bestIndex < 0 || index < bestIndex) {
				bestIndex = index
				bestLength = len(boundary)
			}
			searchFrom = index + len(boundary)
		}
	}
	return bestIndex, bestLength
}

func startsSupportedAction(content string) bool {
	return hasAnyPrefix(content,
		"最近開始", "開始", "持續",
		"後來不再", "後來停止", "後來沒有", "後來不想", "後來放棄",
		"不再", "停止", "沒有", "不想", "放棄", "不參加", "不參賽",
	)
}

func hasTopicReversal(clauses []string, topicMarkers, reversalPatterns []string) bool {
	for _, clause := range clauses {
		for _, pattern := range reversalPatterns {
			searchFrom := 0
			for searchFrom < len(clause) {
				relativeIndex := strings.Index(clause[searchFrom:], pattern)
				if relativeIndex < 0 {
					break
				}
				index := searchFrom + relativeIndex
				suffix := strings.TrimSpace(clause[index+len(pattern):])
				object := reversalObject(suffix)
				if isObjectlessReversal(object) || hasAnySubstring(object, topicMarkers...) {
					return true
				}
				searchFrom = index + len(pattern)
			}
		}
	}
	return false
}

func reversalObject(suffix string) string {
	if boundaryIndex, _ := nextActionBoundary(suffix); boundaryIndex >= 0 {
		suffix = suffix[:boundaryIndex]
	}
	return strings.TrimSpace(suffix)
}

func isObjectlessReversal(object string) bool {
	switch strings.TrimSpace(object) {
	case "", "了":
		return true
	default:
		return false
	}
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
