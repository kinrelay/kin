package socialcontext

import "testing"

func TestSignificanceReversalTopicIgnoresAnotherSubjectAfterCompoundConnector(t *testing.T) {
	if got := significanceReversalTopic("開始準備馬拉松但朋友放棄馬拉松"); got != "" {
		t.Fatalf("significanceReversalTopic() = %q, want empty because the reversal belongs to a friend", got)
	}
}
