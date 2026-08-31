package socialcontext

import (
	"strings"
	"unicode/utf8"
)

// MinimumMeaningfulRunes is the first MVP hypothesis for rejecting overly short Activity signals.
// It is deliberately explicit and deterministic so the product threshold can evolve with evidence.
const MinimumMeaningfulRunes = 12

// SignificanceStatus is the outcome of evaluating one normalized Activity signal.
type SignificanceStatus string

const (
	// SignificanceEligible means the Activity may proceed to a later derivation use case.
	SignificanceEligible SignificanceStatus = "eligible"
	// SignificanceSuppressed means the Activity must not proceed to derivation in this evaluation.
	SignificanceSuppressed SignificanceStatus = "suppressed"
)

// SuppressionReason explains a deterministic Kin-owned significance decision.
type SuppressionReason string

const (
	// SuppressionNone is used for eligible Activities.
	SuppressionNone SuppressionReason = ""
	// SuppressionInvalidSignal means a supposedly normalized Activity is missing required identity/content.
	SuppressionInvalidSignal SuppressionReason = "invalid-signal"
	// SuppressionLowSignal means the Activity is too short to satisfy the first MVP meaningful-signal threshold.
	SuppressionLowSignal SuppressionReason = "low-signal"
	// SuppressionDuplicate means another Activity in the same batch is the canonical representative for equivalent normalized content.
	SuppressionDuplicate SuppressionReason = "duplicate"
)

// SignificanceSignal is the minimal provider-neutral Activity input required by significance policy.
type SignificanceSignal struct {
	ActivityID string
	Content    string
}

// SignificanceDecision records the Activity-scoped result without creating Context Candidate or Social Context state.
type SignificanceDecision struct {
	ActivityID string
	Status     SignificanceStatus
	Reason     SuppressionReason
}

// EvaluateSignificance applies the deterministic MVP 2 significance/suppression policy while preserving input order.
func EvaluateSignificance(signals []SignificanceSignal) []SignificanceDecision {
	canonicalByContent := canonicalActivityIDs(signals)
	decisions := make([]SignificanceDecision, 0, len(signals))

	for _, signal := range signals {
		activityID := strings.TrimSpace(signal.ActivityID)
		content := normalizeSignificanceContent(signal.Content)
		decision := SignificanceDecision{ActivityID: activityID}

		if activityID == "" || content == "" {
			decision.Status = SignificanceSuppressed
			decision.Reason = SuppressionInvalidSignal
			decisions = append(decisions, decision)
			continue
		}

		duplicateKey := strings.ToLower(content)
		if canonicalByContent[duplicateKey] != activityID {
			decision.Status = SignificanceSuppressed
			decision.Reason = SuppressionDuplicate
			decisions = append(decisions, decision)
			continue
		}

		if utf8.RuneCountInString(content) < MinimumMeaningfulRunes && !isExplicitReversalSignal(content) {
			decision.Status = SignificanceSuppressed
			decision.Reason = SuppressionLowSignal
			decisions = append(decisions, decision)
			continue
		}

		decision.Status = SignificanceEligible
		decision.Reason = SuppressionNone
		decisions = append(decisions, decision)
	}

	return decisions
}

// canonicalActivityIDs treats the input order as occurrence chronology and keeps
// the newest representative for equivalent normalized content. This lets a later
// affirmative signal re-establish current state after an intervening reversal.
func canonicalActivityIDs(signals []SignificanceSignal) map[string]string {
	canonical := make(map[string]string, len(signals))
	for _, signal := range signals {
		activityID := strings.TrimSpace(signal.ActivityID)
		content := normalizeSignificanceContent(signal.Content)
		if activityID == "" || content == "" {
			continue
		}
		canonical[strings.ToLower(content)] = activityID
	}
	return canonical
}

func isExplicitReversalSignal(content string) bool {
	content = strings.TrimSpace(content)
	for _, clause := range splitSignificanceClauses(content) {
		if isObjectlessAbandonmentSignal(clause) {
			return true
		}
	}

	if hasTopicBoundSignificanceReversal(content, []string{"分散式系統", "一致性模型"},
		"不再深入研究", "停止深入研究", "沒有深入研究", "不想深入研究", "放棄深入研究",
		"不再研究", "停止研究", "不想研究", "沒有研究", "放棄研究",
		"不再比較", "停止比較", "放棄了", "放棄",
	) {
		return true
	}
	if hasTopicBoundSignificanceReversal(content, []string{"馬拉松"},
		"不再準備", "停止準備", "沒有準備",
		"不再訓練", "停止訓練", "沒有訓練",
		"取消參賽", "取消參加", "不參加", "不參賽", "放棄",
	) {
		return true
	}
	return false
}

func hasTopicBoundSignificanceReversal(content string, markers []string, patterns ...string) bool {
	for _, clause := range splitSignificanceClauses(content) {
		for _, pattern := range patterns {
			searchFrom := 0
			for searchFrom < len(clause) {
				relativeIndex := strings.Index(clause[searchFrom:], pattern)
				if relativeIndex < 0 {
					break
				}
				index := searchFrom + relativeIndex
				if isNegatedSignificanceReversal(clause[:index]) {
					searchFrom = index + len(pattern)
					continue
				}

				object := strings.TrimSpace(clause[index+len(pattern):])
				if object == "" || object == "了" {
					preposedObject := significanceReversalPreposedObject(clause[:index])
					if significanceReversalObjectTargetsTopic(preposedObject, markers) {
						return true
					}
				} else if significanceReversalObjectTargetsTopic(object, markers) {
					return true
				}
				searchFrom = index + len(pattern)
			}
		}
	}
	return false
}

func significanceReversalPreposedObject(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	for _, framing := range []string{"但後來", "後來", "最近", "目前", "現在", "已經"} {
		prefix = strings.TrimSpace(strings.TrimPrefix(prefix, framing))
	}
	return strings.TrimSpace(strings.TrimSuffix(prefix, "我"))
}

func significanceReversalObjectTargetsTopic(object string, markers []string) bool {
	object = strings.TrimSpace(object)
	if len(markers) == 1 && markers[0] == "馬拉松" {
		return isMarathonSignificanceParticipationObject(object)
	}
	return hasAnySignificanceMarker(object, markers...)
}

func isMarathonSignificanceParticipationObject(object string) bool {
	object = strings.TrimSpace(object)
	if object == "馬拉松" {
		return true
	}
	for _, suffix := range []string{"馬拉松訓練", "馬拉松參賽", "馬拉松比賽"} {
		if !strings.HasSuffix(object, suffix) {
			continue
		}
		modifier := strings.TrimSpace(strings.TrimSuffix(object, suffix))
		return modifier == "" || modifier == "第一次全程"
	}
	return false
}

func splitSignificanceClauses(content string) []string {
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

func isNegatedSignificanceReversal(prefix string) bool {
	prefix = strings.TrimSpace(prefix)
	return strings.HasSuffix(prefix, "不會") || strings.HasSuffix(prefix, "不想") || strings.HasSuffix(prefix, "不願") || strings.HasSuffix(prefix, "沒有") || strings.HasSuffix(prefix, "從未")
}

func isObjectlessAbandonmentSignal(content string) bool {
	value := strings.TrimSpace(content)
	for _, framing := range []string{"但後來", "後來", "最近", "目前", "現在", "已經"} {
		value = strings.TrimSpace(strings.TrimPrefix(value, framing))
	}
	return value == "放棄" || value == "放棄了"
}

func hasAnySignificanceMarker(content string, markers ...string) bool {
	for _, marker := range markers {
		if strings.Contains(content, marker) {
			return true
		}
	}
	return false
}

func normalizeSignificanceContent(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
