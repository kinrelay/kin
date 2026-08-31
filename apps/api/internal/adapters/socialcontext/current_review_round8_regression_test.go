package socialcontext

import (
	"context"
	"strings"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorKeepsTechnicalHowItWorksSignalEligible(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-1",
		Content: "最近開始研究分散式系統如何工作",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "分散式系統") {
		t.Fatalf("Generate() meaning = %q, want technical verb 工作 to remain a distributed-systems interest", got.Meaning)
	}
}

func TestDeterministicGeneratorDoesNotBindBareAbandonmentPastLocalUnsupportedAntecedent(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-topic", Content: "最近開始研究分散式系統"},
		{ID: "activity-french", Content: "最近開始學習法文，後來放棄了"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "分散式系統") {
		t.Fatalf("Generate() meaning = %q, want local unsupported antecedent to prevent fallback abandonment of prior distributed-systems state", got.Meaning)
	}
}

func TestDeterministicGeneratorDelimitsMarathonIntentBeforeUnrelatedTrailingClause(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-1",
		Content: "最近開始準備馬拉松比賽但最近工作真的很忙",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want unrelated trailing clause excluded from marathon participation object", got.Meaning)
	}
}

func TestDeterministicGeneratorTreatsHouLaiAsCompleteTemporalBoundary(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-1",
		Content: "停止研究分散式系統後來開始準備馬拉松比賽",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if strings.Contains(got.Meaning, "分散式系統") || !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want reversal followed by later marathon intent", got.Meaning)
	}
}
