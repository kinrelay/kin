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

func canonicalActivityIDs(signals []SignificanceSignal) map[string]string {
	canonical := make(map[string]string, len(signals))
	for _, signal := range signals {
		activityID := strings.TrimSpace(signal.ActivityID)
		content := normalizeSignificanceContent(signal.Content)
		if activityID == "" || content == "" {
			continue
		}
		key := strings.ToLower(content)
		if current, exists := canonical[key]; !exists || activityID < current {
			canonical[key] = activityID
		}
	}
	return canonical
}

func isExplicitReversalSignal(content string) bool {
	if utf8.RuneCountInString(content) < 6 {
		return false
	}
	for _, marker := range []string{
		"不再研究", "停止研究", "不想研究", "沒有研究",
		"不再比較", "停止比較",
		"不再準備", "停止準備", "沒有準備",
		"不再訓練", "停止訓練", "沒有訓練",
		"不參加", "不參賽", "放棄",
	} {
		if strings.Contains(content, marker) {
			return true
		}
	}
	return false
}

func normalizeSignificanceContent(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
