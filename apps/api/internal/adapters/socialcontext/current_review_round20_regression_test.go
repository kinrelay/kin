package socialcontext

import "testing"

func TestCompoundObjectlessAbandonmentKeepsNegatedContinuationBoundToBatchTopic(t *testing.T) {
	recognized := []recognizedSignal{{activityID: "activity-marathon", topic: marathonTopic}}
	bound := compoundObjectlessAbandonmentTopics(
		"我不會取消參加馬拉松，後來放棄了",
		recognized,
	)
	if len(bound) != 1 || bound[0] != marathonTopic {
		t.Fatalf("bound topics = %#v, want [%q] because the negated continuation still refers to the prior marathon topic", bound, marathonTopic)
	}
}
