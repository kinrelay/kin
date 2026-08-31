package socialcontext

import "testing"

func TestEvaluateSignificanceAdmitsDistributedSystemsTechnicalWorkReversal(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{{
		ActivityID: "activity-1",
		Content:    "放棄分散式系統工作負載",
	}})
	if len(decisions) != 1 || decisions[0].Status != SignificanceEligible {
		t.Fatalf("EvaluateSignificance() = %#v, want technical 工作 object to remain an eligible distributed-systems reversal", decisions)
	}
}
