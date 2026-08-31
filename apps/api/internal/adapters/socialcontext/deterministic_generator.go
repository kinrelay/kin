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
	marathonMarkers = []string{"馬拉松"}
	marathonReversals = []string{
		"沒有準備", "不再準備", "停止準備", "沒有訓練", "不再訓練", "停止訓練",
		"取消參賽", "取消參加", "放棄了", "放棄", "不參加", "不參賽",
	}
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
		for _, topic := range compoundObjectlessAbandonmentTopics(activity.Content, recognized) {
			recognized, _ = removeRecognizedTopic(recognized, topic)
		}

		reversed := reversedTopics(activity.Content)
		if isStandaloneObjectlessAbandonment(activity.Content) && len(reversed) > 1 {
			// A bare "放棄了" is context-dependent. Bind it only to the nearest
			// recognized antecedent rather than retracting every topic whose grammar
			// happens to support objectless abandonment.
			if len(recognized) > 0 {
				nearestTopic := recognized[len(recognized)-1].topic
				if _, applies := reversed[nearestTopic]; applies {
					recognized, _ = removeRecognizedTopic(recognized, nearestTopic)
				}
			}
		} else {
			for topic := range reversed {
				recognized, _ = removeRecognizedTopic(recognized, topic)
			}
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
	for _, clause := range clauses {
		if isObjectlessAbandonmentClause(clause) {
			if len(topics) > 0 {
				nearestTopic := topics[len(topics)-1]
				topics = removeTopic(topics, nearestTopic)
			}
			continue
		}

		// Reconcile each clause against the topic state that exists at that point
		// in chronology. Future reversals must not erase an antecedent before an
		// intervening objectless abandonment has had a chance to bind to it.
		distributedReversed := hasTopicReversal([]string{clause}, distributedSystemsMarkers, distributedSystemsReversals)
		marathonReversed := hasTopicReversal([]string{clause}, marathonMarkers, marathonReversals)
		if distributedReversed {
			topics = removeTopic(topics, distributedSystemsTopic)
		}
		if marathonReversed {
			topics = removeTopic(topics, marathonTopic)
		}

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
			if !distributedReversed {
				topics = append(topics, distributedSystemsTopic)
			}
		case isAffirmativeMarathonParticipationClause(clause):
			if !marathonReversed {
				topics = append(topics, marathonTopic)
			}
		}
	}
	return topics
}

func compoundObjectlessAbandonmentTopics(content string, recognized []recognizedSignal) []string {
	clauses := splitSignalClauses(content)
	if len(clauses) < 2 {
		return nil
	}

	batchState := append([]recognizedSignal(nil), recognized...)
	localTopics := make([]string, 0, 2)
	boundTopics := make([]string, 0, 2)
	for _, clause := range clauses {
		if isObjectlessAbandonmentClause(clause) {
			var topic string
			if len(localTopics) > 0 {
				topic = localTopics[len(localTopics)-1]
				localTopics = removeTopic(localTopics, topic)
			} else if len(batchState) > 0 {
				topic = batchState[len(batchState)-1].topic
			}
			if topic != "" {
				boundTopics = append(boundTopics, topic)
				batchState, _ = removeRecognizedTopic(batchState, topic)
			}
			continue
		}

		// Keep the simulated clause-time state aligned with Generate/summarizeSignals.
		// Explicit reversals retire their topic before a later bare abandonment
		// chooses the nearest remaining antecedent; the bare abandonment itself is
		// handled above because its target is context-dependent.
		if hasTopicReversal([]string{clause}, distributedSystemsMarkers, distributedSystemsReversals) {
			localTopics = removeTopic(localTopics, distributedSystemsTopic)
			batchState, _ = removeRecognizedTopic(batchState, distributedSystemsTopic)
		}
		if hasTopicReversal([]string{clause}, marathonMarkers, marathonReversals) {
			localTopics = removeTopic(localTopics, marathonTopic)
			batchState, _ = removeRecognizedTopic(batchState, marathonTopic)
		}

		for _, topic := range summarizeSignals(clause) {
			localTopics = append(localTopics, topic)
		}
	}
	return boundTopics
}

func reversedTopics(content string) map[string]struct{} {
	clauses := splitSignalClauses(content)
	// Within compound content, a bare abandonment is resolved against the
	// nearest antecedent by summarizeSignals/Generate. Do not treat it as an
	// explicit reversal of every topic. Standalone abandonment remains
	// cross-Activity context-dependent and is handled by Generate.
	if len(clauses) > 1 {
		clauses = withoutObjectlessAbandonmentClauses(clauses)
	}
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

func removeTopic(topics []string, topic string) []string {
	kept := topics[:0]
	for _, value := range topics {
		if value != topic {
			kept = append(kept, value)
		}
	}
	return kept
}

func isAffirmativeMarathonParticipationClause(clause string) bool {
	if hasAffirmativeReversalOccurrence(clause, "不參加", "不參賽", "取消參賽", "取消參加") {
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
		target = trimMarathonTrailingDescription(target)
		if reversalObjectTargetsTopic(target, marathonMarkers) {
			return true
		}
	}
	return false
}

func trimMarathonTrailingDescription(target string) string {
	for _, boundary := range []string{"但目前", "但現在"} {
		if index := strings.Index(target, boundary); index >= 0 {
			remainder := strings.TrimSpace(target[index+len(boundary):])
			if strings.HasPrefix(remainder, "沒有受傷") {
				return strings.TrimSpace(target[:index])
			}
		}
	}
	return target
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
	rightSegments := splitCompoundActions(right)
	result := make([]string, 0, 1+len(rightSegments))
	if left != "" {
		result = append(result, left)
	}
	return append(result, rightSegments...)
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
		"後來不再", "後來停止", "後來沒有", "後來不想", "後來放棄", "後來取消",
		"不再", "停止", "沒有", "不想", "放棄", "取消參賽", "取消參加", "不參加", "不參賽",
		"不會停止", "不想停止", "不願停止", "沒有停止", "從未停止",
		"不會放棄", "不想放棄", "不願放棄", "沒有放棄", "從未放棄",
		"不會取消參賽", "不會取消參加", "不願取消參賽", "不願取消參加",
		"沒有取消參賽", "沒有取消參加", "從未取消參賽", "從未取消參加",
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
				if isNegatedReversalOccurrence(clause[:index]) {
					searchFrom = index + len(pattern)
					continue
				}
				suffix := strings.TrimSpace(clause[index+len(pattern):])
				object := reversalObject(suffix)
				if isObjectlessReversal(object) {
					preposedObject := reversalPreposedObject(clause[:index])
					if preposedObject == "" || reversalObjectTargetsTopic(preposedObject, topicMarkers) {
						return true
					}
				} else if reversalObjectTargetsTopic(object, topicMarkers) {
					return true
				}
				searchFrom = index + len(pattern)
			}
		}
	}
	return false
}

func hasAffirmativeReversalOccurrence(content string, patterns ...string) bool {
	for _, pattern := range patterns {
		searchFrom := 0
		for searchFrom < len(content) {
			relativeIndex := strings.Index(content[searchFrom:], pattern)
			if relativeIndex < 0 {
				break
			}
			index := searchFrom + relativeIndex
			if !isNegatedReversalOccurrence(content[:index]) {
				return true
			}
			searchFrom = index + len(pattern)
		}
	}
	return false
}

func isNegatedReversalOccurrence(prefix string) bool {
	prefix = strings.TrimSpace(prefix)
	return strings.HasSuffix(prefix, "不會") || strings.HasSuffix(prefix, "不想") || strings.HasSuffix(prefix, "不願") || strings.HasSuffix(prefix, "沒有") || strings.HasSuffix(prefix, "從未")
}

func reversalObject(suffix string) string {
	if boundaryIndex, _ := nextActionBoundary(suffix); boundaryIndex >= 0 {
		suffix = suffix[:boundaryIndex]
	}
	return strings.TrimSpace(suffix)
}

func reversalPreposedObject(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	for _, framing := range []string{"但後來", "後來", "最近", "目前", "現在", "已經"} {
		prefix = strings.TrimSpace(strings.TrimPrefix(prefix, framing))
	}
	prefix = strings.TrimSpace(strings.TrimSuffix(prefix, "我"))
	return prefix
}

func reversalObjectTargetsTopic(object string, topicMarkers []string) bool {
	object = strings.TrimSpace(object)
	if isMarathonTopicMarkerSet(topicMarkers) {
		return isMarathonParticipationObject(object)
	}
	return hasAnySubstring(object, topicMarkers...)
}

func isMarathonParticipationObject(object string) bool {
	object = strings.TrimSpace(object)
	if object == "馬拉松" {
		return true
	}
	return strings.HasSuffix(object, "馬拉松訓練") ||
		strings.HasSuffix(object, "馬拉松參賽") ||
		strings.HasSuffix(object, "馬拉松比賽")
}

func isMarathonTopicMarkerSet(topicMarkers []string) bool {
	return len(topicMarkers) == 1 && topicMarkers[0] == "馬拉松"
}

func withoutObjectlessAbandonmentClauses(clauses []string) []string {
	filtered := make([]string, 0, len(clauses))
	for _, clause := range clauses {
		if !isObjectlessAbandonmentClause(clause) {
			filtered = append(filtered, clause)
		}
	}
	return filtered
}

func isObjectlessAbandonmentClause(value string) bool {
	value = strings.TrimSpace(value)
	for _, framing := range []string{"但後來", "後來", "最近", "目前", "現在", "已經"} {
		value = strings.TrimSpace(strings.TrimPrefix(value, framing))
	}
	return value == "放棄" || value == "放棄了"
}

func isStandaloneObjectlessAbandonment(content string) bool {
	clauses := splitSignalClauses(content)
	return len(clauses) == 1 && isObjectlessAbandonmentClause(clauses[0])
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