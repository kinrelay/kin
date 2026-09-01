package socialcontext

import "testing"

func TestEvaluateSignificanceKeepsNewestDuplicateRepresentative(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-a", Content: "最近開始研究分散式系統與可靠性工程取捨"},
		{ActivityID: "activity-stop", Content: "不再研究分散式系統"},
		{ActivityID: "activity-z", Content: "最近開始研究分散式系統與可靠性工程取捨"},
	})

	assertDuplicateDecisionByID(t, decisions, "activity-a", SignificanceSuppressed, SuppressionDuplicate)
	assertDuplicateDecisionByID(t, decisions, "activity-z", SignificanceEligible, SuppressionNone)
}

func TestEvaluateSignificanceKeepsShortSupportedMarathonReversalEligible(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{{ActivityID: "activity-stop", Content: "放棄馬拉松"}})
	if len(decisions) != 1 {
		t.Fatalf("decision count = %d, want 1", len(decisions))
	}
	if decisions[0].Status != SignificanceEligible || decisions[0].Reason != SuppressionNone {
		t.Fatalf("decision = %#v, want supported concise marathon reversal eligible", decisions[0])
	}
}
