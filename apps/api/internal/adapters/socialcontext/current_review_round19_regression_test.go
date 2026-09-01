package socialcontext

import "testing"

func TestSummarizeSignalsKeepsNegatedContinuationBoundToTopic(t *testing.T) {
	topics, trailingBarrier := summarizeSignalsWithBarrier("開始準備馬拉松但不會停止準備馬拉松")
	if len(topics) != 1 || topics[0] != marathonTopic {
		t.Fatalf("topics = %#v, want [%q]", topics, marathonTopic)
	}
	if trailingBarrier {
		t.Fatal("trailingBarrier = true, want false because the negated continuation remains bound to the marathon topic")
	}
}

func TestCompoundObjectlessAbandonmentDoesNotJumpPastLocalExplicitReversal(t *testing.T) {
	recognized := []recognizedSignal{{activityID: "activity-marathon", topic: marathonTopic}}
	bound := compoundObjectlessAbandonmentTopics(
		"開始研究分散式系統，停止研究分散式系統，後來放棄了",
		recognized,
	)
	if len(bound) != 0 {
		t.Fatalf("bound topics = %#v, want none because the local distributed-systems reversal is the nearest semantic barrier", bound)
	}
}

func TestSummarizeSignalsRecognizesOwnerSubjectAfterActionConnector(t *testing.T) {
	topics := summarizeSignals("停止研究分散式系統然後我開始準備馬拉松比賽")
	if len(topics) != 1 || topics[0] != marathonTopic {
		t.Fatalf("topics = %#v, want [%q]", topics, marathonTopic)
	}
}
