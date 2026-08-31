package socialcontext

import "testing"

func TestEvaluateSignificancePreservesDuplicateAntecedentBarrierForBareAbandonment(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-distributed", Content: "最近開始深入研究分散式系統設計"},
		{ActivityID: "activity-french-older", Content: "最近開始學習法文並且每天練習口說"},
		{ActivityID: "activity-abandon", Content: "後來放棄了"},
		{ActivityID: "activity-french-newer", Content: "最近開始學習法文並且每天練習口說"},
	})

	assertDuplicateDecisionByID(t, decisions, "activity-distributed", SignificanceEligible, SuppressionNone)
	assertDuplicateDecisionByID(t, decisions, "activity-french-older", SignificanceSuppressed, SuppressionDuplicate)
	assertDuplicateDecisionByID(t, decisions, "activity-abandon", SignificanceSuppressed, SuppressionLowSignal)
	assertDuplicateDecisionByID(t, decisions, "activity-french-newer", SignificanceEligible, SuppressionNone)
}

func TestEvaluateSignificanceRejectsMarathonLogisticsAsSupportedAntecedent(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-distributed", Content: "最近開始深入研究分散式系統設計"},
		{ActivityID: "activity-ticket", Content: "開始準備馬拉松門票"},
		{ActivityID: "activity-abandon", Content: "後來放棄了"},
	})

	assertDuplicateDecisionByID(t, decisions, "activity-distributed", SignificanceEligible, SuppressionNone)
	assertDuplicateDecisionByID(t, decisions, "activity-ticket", SignificanceSuppressed, SuppressionLowSignal)
	assertDuplicateDecisionByID(t, decisions, "activity-abandon", SignificanceSuppressed, SuppressionLowSignal)
}

func TestEvaluateSignificanceSplitsNegatedActionAfterShortMarathonReversal(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-marathon-stop", Content: "放棄馬拉松但不會放棄"},
	})

	assertDuplicateDecisionByID(t, decisions, "activity-marathon-stop", SignificanceEligible, SuppressionNone)
}
