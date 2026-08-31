package socialcontext

import (
	"context"
	"strings"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorTreatsDanShiAsCompleteCompoundBoundary(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-run",
		Content: "開始準備馬拉松但是不會停止準備馬拉松",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want 但是 boundary to preserve affirmative marathon intent", got.Meaning)
	}
}

func TestDeterministicGeneratorAcceptsOrdinaryMarathonOrdinalModifiers(t *testing.T) {
	generator := NewDeterministicGenerator()

	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-run",
		Content: "最近開始準備第二次全程馬拉松比賽",
	}}})
	if err != nil {
		t.Fatalf("Generate() affirmative error = %v", err)
	}
	if !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() affirmative meaning = %q, want ordinary ordinal marathon participation", got.Meaning)
	}

	got, err = generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-run", Content: "開始準備馬拉松"},
		{ID: "activity-stop", Content: "放棄第二次全程馬拉松比賽"},
	}})
	if err != nil {
		t.Fatalf("Generate() reversal error = %v", err)
	}
	if got.Meaning != "" {
		t.Fatalf("Generate() reversal meaning = %q, want ordinary ordinal marathon reversal to retract endurance topic", got.Meaning)
	}
}

func TestDeterministicGeneratorDoesNotBindDistributedSystemsWorkAsTopicReversal(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, reversal := range []string{
		"分散式系統研討會的志工工作我放棄了",
		"後來放棄分散式系統課程的助教工作",
	} {
		t.Run(reversal, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
				{ID: "activity-db", Content: "開始研究分散式系統"},
				{ID: "activity-work", Content: reversal},
			}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if !strings.Contains(got.Meaning, "分散式系統") {
				t.Fatalf("Generate() meaning = %q, want unrelated distributed-systems volunteer/employment duty to preserve topic", got.Meaning)
			}
		})
	}
}
