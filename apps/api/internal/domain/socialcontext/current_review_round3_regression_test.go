package socialcontext

import "testing"

func TestEvaluateSignificanceDoesNotLetNegatedShortReversalsBypassLowSignal(t *testing.T) {
	for _, content := range []string{
		"沒有停止研究分散式系統",
		"從未放棄馬拉松",
	} {
		t.Run(content, func(t *testing.T) {
			decisions := EvaluateSignificance([]SignificanceSignal{{ActivityID: "activity-1", Content: content}})
			if len(decisions) != 1 || decisions[0].Status != SignificanceSuppressed || decisions[0].Reason != SuppressionLowSignal {
				t.Fatalf("EvaluateSignificance(%q) = %#v, want negated short reversal to remain low-signal", content, decisions)
			}
		})
	}
}

func TestEvaluateSignificanceDoesNotBindShortReversalToUnrelatedTopicMarker(t *testing.T) {
	content := "馬拉松，放棄午餐"
	decisions := EvaluateSignificance([]SignificanceSignal{{ActivityID: "activity-1", Content: content}})
	if len(decisions) != 1 || decisions[0].Status != SignificanceSuppressed || decisions[0].Reason != SuppressionLowSignal {
		t.Fatalf("EvaluateSignificance(%q) = %#v, want unrelated short reversal to remain low-signal", content, decisions)
	}
}

func TestEvaluateSignificanceDoesNotTreatPreposedMarathonWorkAsParticipationReversal(t *testing.T) {
	content := "馬拉松工作我放棄了"
	decisions := EvaluateSignificance([]SignificanceSignal{{ActivityID: "activity-1", Content: content}})
	if len(decisions) != 1 || decisions[0].Status != SignificanceSuppressed || decisions[0].Reason != SuppressionLowSignal {
		t.Fatalf("EvaluateSignificance(%q) = %#v, want marathon work abandonment to remain low-signal", content, decisions)
	}
}
