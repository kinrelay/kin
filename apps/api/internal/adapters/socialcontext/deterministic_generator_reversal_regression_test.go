package socialcontext

import (
	"context"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorFindsLaterSameTopicReversalAfterUnrelatedReversalObject(t *testing.T) {
	generator := NewDeterministicGenerator()
	raw := "最近開始研究分散式系統但後來停止研究英文並停止研究分散式系統"

	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-db", Content: raw}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning != "" || len(got.Provenance) != 0 {
		t.Fatalf("Generate() = %#v, want later same-topic reversal to suppress the context", got)
	}
}

func TestDeterministicGeneratorTreatsSentenceFinalParticleAsObjectlessReversal(t *testing.T) {
	generator := NewDeterministicGenerator()
	raw := "最近開始研究分散式系統但後來停止研究了"

	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-db", Content: raw}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning != "" || len(got.Provenance) != 0 {
		t.Fatalf("Generate() = %#v, want objectless reversal with sentence-final particle to suppress the context", got)
	}
}

func TestDeterministicGeneratorReconcilesTopicReversalAcrossActivities(t *testing.T) {
	generator := NewDeterministicGenerator()
	input := appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-start", Content: "最近開始研究分散式系統與可靠性工程取捨"},
		{ID: "activity-stop", Content: "最近已經不再研究分散式系統了，改為研究英文"},
	}}

	got, err := generator.Generate(context.Background(), input)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning != "" || len(got.Provenance) != 0 {
		t.Fatalf("Generate() = %#v, want a later cross-activity reversal to suppress the stale topic", got)
	}
}
