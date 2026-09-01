package socialcontext

import "testing"

func TestSplitSignificanceCompoundActionsUsesEarliestValidBoundary(t *testing.T) {
	got := splitSignificanceCompoundActions("開始準備馬拉松比賽且最近開始學習法文但不再研究分散式系統")
	want := []string{"開始準備馬拉松比賽", "最近開始學習法文", "不再研究分散式系統"}
	if len(got) != len(want) {
		t.Fatalf("split clauses = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("split clause[%d] = %q, want %q; all = %#v", i, got[i], want[i], got)
		}
	}
}

func TestEvaluateSignificanceAllowsBareAbandonmentToBindSuppressedSameTopicAntecedent(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-distributed-long", Content: "最近開始深入研究分散式系統設計與可靠性工程取捨"},
		{ActivityID: "activity-distributed-short", Content: "開始研究分散式系統"},
		{ActivityID: "activity-abandon", Content: "後來放棄了"},
	})

	if len(decisions) != 3 {
		t.Fatalf("decisions = %#v, want 3 decisions", decisions)
	}
	if decisions[0].Status != SignificanceEligible {
		t.Fatalf("first decision = %#v, want eligible", decisions[0])
	}
	if decisions[1].Status != SignificanceSuppressed || decisions[1].Reason != SuppressionLowSignal {
		t.Fatalf("second decision = %#v, want low-signal suppression", decisions[1])
	}
	if decisions[2].Status != SignificanceEligible {
		t.Fatalf("bare abandonment decision = %#v, want eligible because nearest suppressed antecedent is the same active topic", decisions[2])
	}
}
