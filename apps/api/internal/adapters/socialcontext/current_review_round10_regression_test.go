package socialcontext

import (
	"context"
	"strings"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorDoesNotBindBareAbandonmentAcrossEligibleUnsupportedAntecedent(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-distributed", Content: "最近開始研究分散式系統"},
		{ID: "activity-french", Content: "最近開始學習法文並且每天練習口說"},
		{ID: "activity-abandon", Content: "後來放棄了"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "分散式系統") {
		t.Fatalf("Generate() meaning = %q, want unsupported French activity to block bare abandonment from retracting older distributed-systems topic", got.Meaning)
	}
}
