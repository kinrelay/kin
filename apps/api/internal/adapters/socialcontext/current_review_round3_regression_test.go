package socialcontext

import (
	"context"
	"strings"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorPreservesAffirmativeIntentWithUnwillingCancellation(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-run",
		Content: "開始準備馬拉松但不願取消參賽",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want 不願取消參賽 to preserve affirmative marathon intent", got.Meaning)
	}
}

func TestDeterministicGeneratorPreservesAffirmativeIntentWithNegatedStop(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-run",
		Content: "開始準備馬拉松但不會停止準備馬拉松",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want negated stop to preserve affirmative marathon intent", got.Meaning)
	}
}

func TestDeterministicGeneratorPreservesAffirmativeIntentWithNegatedAbandonment(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-run",
		Content: "開始準備馬拉松但不會放棄馬拉松",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want negated abandonment to preserve affirmative marathon intent", got.Meaning)
	}
}

func TestDeterministicGeneratorRejectsSpectatorMarathonLogistics(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-ticket",
		Content: "開始準備馬拉松比賽的觀賽門票",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning != "" || len(got.Provenance) != 0 {
		t.Fatalf("Generate() = %#v, want spectator logistics rejected as non-participation", got)
	}
}

func TestDeterministicGeneratorKeepsExplicitReversalAsLocalBareAbandonmentBarrier(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-db", Content: "最近開始研究分散式系統"},
		{ID: "activity-run", Content: "最近開始準備馬拉松"},
		{ID: "activity-stop", Content: "停止準備馬拉松，後來放棄了"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "分散式系統") || strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() = %#v, want explicit marathon reversal to block the later bare abandonment from jumping back to older distributed-systems state", got)
	}
}

func TestDeterministicGeneratorPreservesAffirmativeIntentWithFramedNegatedAbandonment(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, content := range []string{
		"開始準備馬拉松但後來不會放棄馬拉松",
		"開始準備馬拉松但後來不願放棄馬拉松",
		"開始準備馬拉松但後來從未放棄馬拉松",
	} {
		t.Run(content, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
				ID:      "activity-run",
				Content: content,
			}}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if !strings.Contains(got.Meaning, "耐力運動") {
				t.Fatalf("Generate() meaning = %q, want framed negated abandonment to preserve affirmative marathon intent", got.Meaning)
			}
		})
	}
}

func TestDeterministicGeneratorRejectsSpectatorIntentEndingInMarathonRace(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, content := range []string{
		"開始準備觀看馬拉松比賽",
		"開始準備應援馬拉松比賽",
	} {
		t.Run(content, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
				ID:      "activity-spectator",
				Content: content,
			}}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if got.Meaning != "" || len(got.Provenance) != 0 {
				t.Fatalf("Generate() = %#v, want spectator intent rejected as non-participation", got)
			}
		})
	}
}

func TestDeterministicGeneratorDoesNotUseSpectatorIntentAsPreposedMarathonReversal(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-run", Content: "開始準備馬拉松"},
		{ID: "activity-spectator", Content: "觀看馬拉松比賽我放棄了"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want spectator abandonment to preserve existing marathon participation", got.Meaning)
	}
}
