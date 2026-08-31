package socialcontext

import (
	"context"
	"strings"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorReconcilesEveryCompoundObjectlessAbandonment(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-old-db", Content: "最近開始研究分散式系統"},
		{ID: "activity-old-run", Content: "最近開始準備馬拉松"},
		{ID: "activity-new", Content: "開始研究分散式系統，後來放棄了，開始準備馬拉松，後來放棄了"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning != "" || len(got.Provenance) != 0 {
		t.Fatalf("Generate() = %#v, want both compound objectless abandonments to retract prior batch topics", got)
	}
}

func TestDeterministicGeneratorBindsLeadingObjectlessAbandonmentToPriorBatchState(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-run", Content: "最近開始準備馬拉松"},
		{ID: "activity-stop", Content: "後來放棄了，最近工作很忙"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning != "" || len(got.Provenance) != 0 {
		t.Fatalf("Generate() = %#v, want leading objectless abandonment to retract nearest prior batch topic", got)
	}
}

func TestDeterministicGeneratorIgnoresNegatedReversalPatterns(t *testing.T) {
	generator := NewDeterministicGenerator()
	cases := []struct {
		name string
		base string
		keep string
		want string
	}{
		{name: "distributed stopping", base: "最近開始深入研究分散式系統", keep: "不會停止研究分散式系統", want: "分散式系統"},
		{name: "marathon cancellation", base: "最近開始準備馬拉松", keep: "我不會取消參加馬拉松比賽", want: "耐力運動"},
		{name: "marathon no abandonment", base: "最近開始準備馬拉松", keep: "我沒有放棄馬拉松比賽", want: "耐力運動"},
		{name: "marathon no cancellation", base: "最近開始準備馬拉松", keep: "我沒有取消參加馬拉松比賽", want: "耐力運動"},
		{name: "distributed never stopping", base: "最近開始深入研究分散式系統", keep: "我從未停止研究分散式系統", want: "分散式系統"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
				{ID: "activity-base", Content: tc.base},
				{ID: "activity-keep", Content: tc.keep},
			}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if !strings.Contains(got.Meaning, tc.want) {
				t.Fatalf("Generate() meaning = %q, want negated reversal to preserve %q topic", got.Meaning, tc.want)
			}
		})
	}
}

func TestDeterministicGeneratorPreservesAffirmativeIntentWithNegatedCancellation(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-run",
		Content: "最近開始準備馬拉松但不會取消參賽",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want negated cancellation to preserve affirmative marathon intent", got.Meaning)
	}
}

func TestDeterministicGeneratorBindsBareAbandonmentBeforeLaterExplicitReversal(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-compound",
		Content: "開始研究分散式系統，開始準備馬拉松，後來放棄了，停止準備馬拉松",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "分散式系統") || strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want bare abandonment bound to marathon state at that clause and distributed-systems preserved", got.Meaning)
	}
}

func TestDeterministicGeneratorRecognizesModifiedMarathonRaceIntent(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-run",
		Content: "最近開始準備第一次全程馬拉松比賽",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want modified marathon race participation recognized", got.Meaning)
	}
}
