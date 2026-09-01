package socialcontext

import "testing"

func TestEvaluateSignificanceCarriesSuppressedAntecedentAcrossOtherEligibleTopics(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-distributed-long", Content: "最近開始深入研究分散式系統設計與可靠性工程取捨"},
		{ActivityID: "activity-marathon", Content: "最近開始準備第一次全程馬拉松比賽"},
		{ActivityID: "activity-distributed-short", Content: "開始研究分散式系統"},
		{ActivityID: "activity-abandon", Content: "後來放棄了"},
	})

	if len(decisions) != 4 {
		t.Fatalf("decisions = %#v, want 4 decisions", decisions)
	}
	if decisions[2].Status != SignificanceSuppressed || decisions[2].Reason != SuppressionLowSignal {
		t.Fatalf("suppressed antecedent = %#v, want low-signal suppression", decisions[2])
	}
	if decisions[3].Status != SignificanceEligible {
		t.Fatalf("bare abandonment = %#v, want eligible because its nearest suppressed antecedent is the still-active distributed-systems topic", decisions[3])
	}
}

func TestEvaluateSignificanceStripsBlockedAbandonmentBeforeTemporalSupportedAction(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-distributed", Content: "最近開始深入研究分散式系統設計與可靠性工程取捨"},
		{ActivityID: "activity-unsupported", Content: "最近開始學習法文"},
		{ActivityID: "activity-compound", Content: "後來放棄了然後開始準備馬拉松比賽"},
	})

	if len(decisions) != 3 {
		t.Fatalf("decisions = %#v, want 3 decisions", decisions)
	}
	if decisions[1].Status != SignificanceSuppressed {
		t.Fatalf("unsupported antecedent = %#v, want suppressed", decisions[1])
	}
	if decisions[2].Status != SignificanceEligible {
		t.Fatalf("compound decision = %#v, want eligible due to independent marathon action", decisions[2])
	}
	if got, want := decisions[2].DerivationContent, "開始準備馬拉松比賽"; got != want {
		t.Fatalf("compound derivation content = %q, want %q", got, want)
	}
}
