package socialcontext

import (
	"context"
	"strings"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorPreservesTechnicalDistributedSystemsWorkTerms(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, content := range []string{
		"最近開始研究分散式系統的工作原理",
		"最近開始研究分散式系統工作負載",
	} {
		t.Run(content, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
				ID:      "activity-1",
				Content: content,
			}}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if !strings.Contains(got.Meaning, "分散式系統") {
				t.Fatalf("Generate() meaning = %q, want technical 工作 term to remain a distributed-systems interest", got.Meaning)
			}
		})
	}
}

func TestDeterministicGeneratorTreatsSubjectBeforeTemporalFrameAsObjectlessAbandonment(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-topic", Content: "開始研究分散式系統"},
		{ID: "activity-stop", Content: "我後來放棄了，最近工作真的很忙碌"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning != "" {
		t.Fatalf("Generate() meaning = %q, want subject+temporal framed bare abandonment to retract nearest antecedent", got.Meaning)
	}
}
