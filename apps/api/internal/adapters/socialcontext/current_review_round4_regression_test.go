package socialcontext

import (
	"context"
	"strings"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorPreservesFramedNegatedStopAndCancellation(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, content := range []string{
		"開始準備馬拉松但後來不會停止準備馬拉松",
		"開始準備馬拉松但後來不會取消參賽",
	} {
		t.Run(content, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-run", Content: content}}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if !strings.Contains(got.Meaning, "耐力運動") {
				t.Fatalf("Generate() meaning = %q, want framed negated reversal to preserve marathon intent", got.Meaning)
			}
		})
	}
}

func TestDeterministicGeneratorBindsSubjectOnlyAbandonmentToNearestAntecedent(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-db", Content: "開始研究分散式系統"},
		{ID: "activity-run", Content: "開始準備馬拉松"},
		{ID: "activity-abandon", Content: "我放棄了，最近工作真的非常忙碌"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "分散式系統") || strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want subject-only abandonment to retract only nearest marathon antecedent", got.Meaning)
	}
}

func TestDeterministicGeneratorAcceptsTerminalParticleOnMarathonReversal(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-run", Content: "開始準備馬拉松"},
		{ID: "activity-stop", Content: "經過慎重考慮之後我放棄馬拉松了"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning != "" {
		t.Fatalf("Generate() meaning = %q, want terminal particle reversal to retract marathon", got.Meaning)
	}
}

func TestDeterministicGeneratorStopsReversalObjectAtUnrelatedTrailingClause(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-run", Content: "開始準備馬拉松"},
		{ID: "activity-stop", Content: "後來放棄馬拉松比賽但最近工作真的很忙"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning != "" {
		t.Fatalf("Generate() meaning = %q, want explicit marathon abandonment before unrelated trailing clause to retract marathon", got.Meaning)
	}
}

func TestDeterministicGeneratorTreatsNoLongerAbandonmentAsNegatedReversal(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-run", Content: "開始準備馬拉松"},
		{ID: "activity-continue", Content: "我不再放棄馬拉松比賽"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want 不再放棄 to preserve existing marathon intent", got.Meaning)
	}
}
