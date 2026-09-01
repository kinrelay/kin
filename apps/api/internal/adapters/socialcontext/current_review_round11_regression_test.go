package socialcontext

import (
	"context"
	"strings"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorKeepsExplicitReversalAsBareAbandonmentAntecedentBarrier(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-distributed", Content: "最近開始研究分散式系統"},
		{ID: "activity-marathon-stop", Content: "停止準備馬拉松"},
		{ID: "activity-abandon", Content: "後來放棄了"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "分散式系統") {
		t.Fatalf("Generate() meaning = %q, want explicit marathon reversal to remain the nearest antecedent barrier for later bare abandonment", got.Meaning)
	}
}
