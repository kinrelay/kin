package socialcontext

import "testing"

func TestSignificanceReversalTopicIgnoresPluralOtherSubjectAfterCompoundConnector(t *testing.T) {
	for _, subject := range []string{"我的朋友們", "朋友們", "同事們", "家人們"} {
		content := "開始準備馬拉松但" + subject + "放棄馬拉松"
		if got := significanceReversalTopic(content); got != "" {
			t.Fatalf("significanceReversalTopic(%q) = %q, want empty because the reversal belongs to another subject", content, got)
		}
	}
}
