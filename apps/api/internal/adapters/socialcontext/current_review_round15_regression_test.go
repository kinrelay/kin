package socialcontext

import (
	"context"
	"testing"

	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorPreservesTrailingUnsupportedClauseAsNextActivityBarrier(t *testing.T) {
	generator := NewDeterministicGenerator()
	generated, err := generator.Generate(context.Background(), applicationsocialcontext.ContextGenerationInput{
		Activities: []applicationsocialcontext.ContextGenerationActivity{
			{ID: "activity-mixed", Content: "最近開始研究分散式系統，最近開始學習法文並且每天練習口說"},
			{ID: "activity-abandon", Content: "後來放棄了"},
		},
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	wantMeaning := "近期關注" + distributedSystemsTopic
	if generated.Meaning != wantMeaning {
		t.Fatalf("Generate().Meaning = %q, want %q", generated.Meaning, wantMeaning)
	}
	if len(generated.Provenance) != 1 || generated.Provenance[0] != "activity-mixed" {
		t.Fatalf("Generate().Provenance = %#v, want [activity-mixed]", generated.Provenance)
	}
}
