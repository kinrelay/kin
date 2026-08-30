package socialcontext

import "testing"

func TestEvaluateSignificanceMarksMeaningfulSignalEligibleWithoutCreatingContext(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-1", Content: "最近持續研究 distributed systems 的 consistency trade-offs"},
	})

	if len(decisions) != 1 {
		t.Fatalf("decision count = %d, want 1", len(decisions))
	}
	decision := decisions[0]
	if decision.ActivityID != "activity-1" || decision.Status != SignificanceEligible || decision.Reason != SuppressionNone {
		t.Fatalf("decision = %#v, want eligible without suppression reason", decision)
	}
}

func TestEvaluateSignificanceSuppressesBlankOrInvalidNormalizedSignal(t *testing.T) {
	tests := []SignificanceSignal{
		{ActivityID: "", Content: "meaningful content that should still fail without id"},
		{ActivityID: "activity-2", Content: " \n\t "},
	}

	decisions := EvaluateSignificance(tests)
	if len(decisions) != len(tests) {
		t.Fatalf("decision count = %d, want %d", len(decisions), len(tests))
	}
	for i, decision := range decisions {
		if decision.Status != SignificanceSuppressed || decision.Reason != SuppressionInvalidSignal {
			t.Fatalf("decision[%d] = %#v, want invalid-signal suppression", i, decision)
		}
	}
}

func TestEvaluateSignificanceSuppressesLowSignalContentDeterministically(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-short", Content: "看影片"},
		{ActivityID: "activity-boundary", Content: "123456789012"},
	})

	if got, want := MinimumMeaningfulRunes, 12; got != want {
		t.Fatalf("MinimumMeaningfulRunes = %d, want %d", got, want)
	}
	if decisions[0].Status != SignificanceSuppressed || decisions[0].Reason != SuppressionLowSignal {
		t.Fatalf("short decision = %#v, want low-signal suppression", decisions[0])
	}
	if decisions[1].Status != SignificanceEligible || decisions[1].Reason != SuppressionNone {
		t.Fatalf("boundary decision = %#v, want eligible", decisions[1])
	}
}

func TestEvaluateSignificanceSuppressesLaterExactNormalizedDuplicates(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: "activity-1", Content: "Reading   Distributed Systems Papers"},
		{ActivityID: "activity-2", Content: " reading distributed systems papers "},
		{ActivityID: "activity-3", Content: "Comparing consensus protocols and trade-offs"},
	})

	if decisions[0].Status != SignificanceEligible {
		t.Fatalf("first decision = %#v, want eligible", decisions[0])
	}
	if decisions[1].Status != SignificanceSuppressed || decisions[1].Reason != SuppressionDuplicate {
		t.Fatalf("duplicate decision = %#v, want duplicate suppression", decisions[1])
	}
	if decisions[2].Status != SignificanceEligible {
		t.Fatalf("distinct decision = %#v, want eligible", decisions[2])
	}
}

func TestEvaluateSignificancePreservesInputOrderAndActivityIdentity(t *testing.T) {
	decisions := EvaluateSignificance([]SignificanceSignal{
		{ActivityID: " activity-b ", Content: "A sufficiently meaningful activity signal"},
		{ActivityID: "activity-a", Content: "Another sufficiently meaningful activity"},
	})

	if decisions[0].ActivityID != "activity-b" || decisions[1].ActivityID != "activity-a" {
		t.Fatalf("activity identity/order = %#v, want normalized input order", decisions)
	}
}
