package socialcontext

import "testing"

func TestEvaluateSignificanceSplitsBlockedAbandonmentBeforeZhihouSupportedAction(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-distributed", Content: "最近開始深入研究分散式系統設計與可靠性工程取捨"},
		{ActivityID: "activity-unsupported", Content: "最近開始學習法文"},
		{ActivityID: "activity-compound", Content: "後來放棄了之後開始準備馬拉松比賽"},
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

func TestSignificanceReversalTopicIgnoresTemporallyFramedOtherSubject(t *testing.T) {
	for _, content := range []string{
		"後來朋友放棄馬拉松",
		"但後來同事停止研究分散式系統",
	} {
		if got := significanceReversalTopic(content); got != "" {
			t.Fatalf("significanceReversalTopic(%q) = %q, want empty because the reversal belongs to another subject", content, got)
		}
	}
}

func TestSupportedSignificanceAntecedentTrimsTrailingDescription(t *testing.T) {
	content := "開始準備馬拉松比賽但最近工作真的很忙"
	if got, want := supportedSignificanceAntecedentTopic(content), "marathon"; got != want {
		t.Fatalf("supportedSignificanceAntecedentTopic(%q) = %q, want %q", content, got, want)
	}
}
