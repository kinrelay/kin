package socialcontext

import (
	"context"
	"strings"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorDoesNotTreatMarathonVolunteerAbandonmentAsParticipationReversal(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-run", Content: "最近開始準備馬拉松"},
		{ID: "activity-volunteer-stop", Content: "馬拉松比賽的志工工作我放棄了"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning == "" || !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want marathon participation preserved when abandonment targets volunteer work", got.Meaning)
	}
}

func TestDeterministicGeneratorPreservesMarathonIntentBeforeUnpunctuatedDescription(t *testing.T) {
	generator := NewDeterministicGenerator()
	raw := "最近開始準備第一次全程馬拉松訓練但目前沒有受傷"
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-run", Content: raw}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning == "" || !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want marathon participation preserved before unrelated trailing description", got.Meaning)
	}
}

func TestDeterministicGeneratorBindsObjectlessAbandonmentToNearestAntecedent(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-db", Content: "最近開始研究分散式系統"},
		{ID: "activity-run", Content: "最近開始準備馬拉松"},
		{ID: "activity-stop", Content: "後來放棄了"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning == "" || !strings.Contains(got.Meaning, "分散式系統") || strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want nearest marathon antecedent reversed while distributed-systems topic remains", got.Meaning)
	}
}

func TestDeterministicGeneratorBindsCompoundObjectlessAbandonmentToNearestAntecedent(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-compound",
		Content: "最近開始研究分散式系統，開始準備馬拉松，後來放棄了",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "分散式系統") || strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want nearest marathon antecedent abandoned while distributed-systems remains", got.Meaning)
	}
}

func TestDeterministicGeneratorDoesNotTreatPostposedMarathonLogisticsAsParticipationReversal(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-run", Content: "最近開始準備馬拉松"},
		{ID: "activity-photo", Content: "後來放棄馬拉松比賽的攝影工作"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want marathon participation preserved when abandonment targets photography work", got.Meaning)
	}
}

func TestDeterministicGeneratorRejectsCancellationAfterMarathonTarget(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-cancelled",
		Content: "最近開始準備馬拉松但目前取消參賽",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning != "" || len(got.Provenance) != 0 {
		t.Fatalf("Generate() = %#v, want cancelled participation suppressed", got)
	}
}
