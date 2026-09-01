package socialcontext

import "testing"

func TestEvaluateSignificancePreservesSuppressedAntecedentBarrierForBareAbandonment(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-distributed", Content: "最近開始深入研究分散式系統設計"},
		{ActivityID: "activity-french", Content: "最近開始學習法文"},
		{ActivityID: "activity-abandon", Content: "後來放棄了"},
	})

	assertDuplicateDecisionByID(t, decisions, "activity-distributed", SignificanceEligible, SuppressionNone)
	assertDuplicateDecisionByID(t, decisions, "activity-french", SignificanceSuppressed, SuppressionLowSignal)
	assertDuplicateDecisionByID(t, decisions, "activity-abandon", SignificanceSuppressed, SuppressionLowSignal)
}

func TestEvaluateSignificanceRecognizesCompoundMarathonReversalBeforeTrailingAction(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-marathon-stop", Content: "放棄馬拉松但開始工作"},
	})

	assertDuplicateDecisionByID(t, decisions, "activity-marathon-stop", SignificanceEligible, SuppressionNone)
}

func TestEvaluateSignificanceKeepsRoleDutyAntecedentAsBarrier(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-role-duty", Content: "開始研究分散式系統志工"},
		{ActivityID: "activity-abandon", Content: "放棄了"},
	})

	assertDuplicateDecisionByID(t, decisions, "activity-role-duty", SignificanceSuppressed, SuppressionLowSignal)
	assertDuplicateDecisionByID(t, decisions, "activity-abandon", SignificanceSuppressed, SuppressionLowSignal)
}

func TestEvaluateSignificanceTrimsTemporalBoundaryForShortMarathonReversal(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-marathon-stop", Content: "放棄馬拉松後開始工作"},
	})

	assertDuplicateDecisionByID(t, decisions, "activity-marathon-stop", SignificanceEligible, SuppressionNone)
}
