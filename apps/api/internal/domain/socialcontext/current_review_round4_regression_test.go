package socialcontext

import "testing"

func TestEvaluateSignificanceDoesNotTreatNoLongerAbandonmentAsReversal(t *testing.T) {
	content := "我不再放棄馬拉松比賽"
	decisions := EvaluateSignificance([]SignificanceSignal{{ActivityID: "activity-1", Content: content}})
	if len(decisions) != 1 || decisions[0].Status != SignificanceSuppressed || decisions[0].Reason != SuppressionLowSignal {
		t.Fatalf("EvaluateSignificance(%q) = %#v, want negated abandonment to remain low-signal", content, decisions)
	}
}
