package socialcontext

import "testing"

func TestEvaluateSignificancePreservesRepeatedObjectlessAbandonments(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-distributed", Content: "最近開始研究分散式系統"},
		{ActivityID: "activity-stop-distributed", Content: "後來放棄了"},
		{ActivityID: "activity-marathon", Content: "最近開始準備馬拉松比賽"},
		{ActivityID: "activity-stop-marathon", Content: "後來放棄了"},
	})

	assertDuplicateDecisionByID(t, decisions, "activity-stop-distributed", SignificanceEligible, SuppressionNone)
	assertDuplicateDecisionByID(t, decisions, "activity-stop-marathon", SignificanceEligible, SuppressionNone)
}
