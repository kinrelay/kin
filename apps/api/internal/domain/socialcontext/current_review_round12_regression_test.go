package socialcontext

import "testing"

func TestEvaluateSignificanceKeepsSuppressedSupportedAntecedentAsBarrier(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-distributed", Content: "最近開始深入研究分散式系統設計"},
		{ActivityID: "activity-marathon", Content: "開始準備馬拉松"},
		{ActivityID: "activity-abandon", Content: "後來放棄了"},
	})

	assertDuplicateDecisionByID(t, decisions, "activity-distributed", SignificanceEligible, SuppressionNone)
	assertDuplicateDecisionByID(t, decisions, "activity-marathon", SignificanceSuppressed, SuppressionLowSignal)
	assertDuplicateDecisionByID(t, decisions, "activity-abandon", SignificanceSuppressed, SuppressionLowSignal)
}

func TestEvaluateSignificancePreservesIndependentClauseBesideBlockedAbandonment(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-french", Content: "最近開始學習法文"},
		{ActivityID: "activity-compound", Content: "後來放棄了，最近開始深入研究分散式系統設計"},
	})

	assertDuplicateDecisionByID(t, decisions, "activity-french", SignificanceSuppressed, SuppressionLowSignal)
	assertDuplicateDecisionByID(t, decisions, "activity-compound", SignificanceEligible, SuppressionNone)
}
