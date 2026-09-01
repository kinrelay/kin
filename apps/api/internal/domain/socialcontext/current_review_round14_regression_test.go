package socialcontext

import "testing"

func TestEvaluateSignificanceSanitizesConnectorDelimitedBlockedAbandonment(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-distributed", Content: "最近開始深入研究分散式系統設計"},
		{ActivityID: "activity-french", Content: "最近開始學習法文"},
		{ActivityID: "activity-compound", Content: "後來放棄了但開始準備馬拉松比賽"},
	})

	if got := decisions[2].Status; got != SignificanceEligible {
		t.Fatalf("compound status = %q, want eligible", got)
	}
	if got := decisions[2].DerivationContent; got != "開始準備馬拉松比賽" {
		t.Fatalf("compound derivation content = %q, want blocked abandonment removed", got)
	}
}

func TestEvaluateSignificanceKeepsGeneralizedMarathonOrdinalBesideBlockedAbandonment(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-distributed", Content: "最近開始深入研究分散式系統設計"},
		{ActivityID: "activity-french", Content: "最近開始學習法文"},
		{ActivityID: "activity-compound", Content: "後來放棄了，開始準備第二次全程馬拉松比賽"},
	})

	if got := decisions[2].Status; got != SignificanceEligible {
		t.Fatalf("compound status = %q, want eligible", got)
	}
	if got := decisions[2].DerivationContent; got != "開始準備第二次全程馬拉松比賽" {
		t.Fatalf("compound derivation content = %q, want generalized ordinal clause", got)
	}
}
