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
	// SuppressionDuplicate means an earlier Activity in the same batch has equivalent normalized content.
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

// EvaluateSignificance applies the deterministic MVP 2 significance/suppression policy in input order.
func EvaluateSignificance(signals []SignificanceSignal) []SignificanceDecision {
	decisions := make([]SignificanceDecision, 0, len(signals))
	seenContent := make(map[string]struct{}, len(signals))

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
		if _, exists := seenContent[duplicateKey]; exists {
			decision.Status = SignificanceSuppressed
			decision.Reason = SuppressionDuplicate
			decisions = append(decisions, decision)
			continue
		}
		seenContent[duplicateKey] = struct{}{}

		if utf8.RuneCountInString(content) < MinimumMeaningfulRunes {
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

func normalizeSignificanceContent(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
