package socialcontext

import "testing"

func TestEvaluateSignificanceStripsBlockedAbandonmentFromIndependentCompoundSignal(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-distributed", Content: "最近開始深入研究分散式系統設計"},
		{ActivityID: "activity-french", Content: "最近開始學習法文"},
		{ActivityID: "activity-compound", Content: "後來放棄了，最近開始準備馬拉松比賽"},
	})

	assertDuplicateDecisionByID(t, decisions, "activity-distributed", SignificanceEligible, SuppressionNone)
	assertDuplicateDecisionByID(t, decisions, "activity-french", SignificanceSuppressed, SuppressionLowSignal)
	assertDuplicateDecisionByID(t, decisions, "activity-compound", SignificanceEligible, SuppressionNone)

	for _, decision := range decisions {
		if decision.ActivityID != "activity-compound" {
			continue
		}
		if decision.DerivationContent != "最近開始準備馬拉松比賽" {
			t.Fatalf("DerivationContent = %q, want independent supported clause only", decision.DerivationContent)
		}
		return
	}
	t.Fatal("missing activity-compound decision")
}
