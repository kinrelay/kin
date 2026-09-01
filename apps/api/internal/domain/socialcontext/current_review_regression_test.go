package socialcontext

import "testing"

func TestEvaluateSignificanceDoesNotLetUnsupportedShortAbandonmentBypassLowSignal(t *testing.T) {
	for _, content := range []string{"放棄午餐", "不會放棄午餐"} {
		t.Run(content, func(t *testing.T) {
			decisions := EvaluateSignificance([]SignificanceSignal{{ActivityID: "activity-1", Content: content}})
			if len(decisions) != 1 || decisions[0].Status != SignificanceSuppressed || decisions[0].Reason != SuppressionLowSignal {
				t.Fatalf("EvaluateSignificance(%q) = %#v, want low-signal suppression", content, decisions)
			}
		})
	}
}
