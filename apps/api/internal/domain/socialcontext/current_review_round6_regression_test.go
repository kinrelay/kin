package socialcontext

import "testing"

func TestEvaluateSignificanceDoesNotAdmitDistributedSystemsRoleDutyReversal(t *testing.T) {
	for _, content := range []string{
		"放棄分散式系統助教工作",
		"放棄分散式系統志工工作",
	} {
		t.Run(content, func(t *testing.T) {
			decisions := EvaluateSignificance([]SignificanceSignal{{ActivityID: "activity-1", Content: content}})
			if len(decisions) != 1 || decisions[0].Status != SignificanceSuppressed || decisions[0].Reason != SuppressionLowSignal {
				t.Fatalf("EvaluateSignificance(%q) = %#v, want role-duty reversal to remain low-signal", content, decisions)
			}
		})
	}
}
