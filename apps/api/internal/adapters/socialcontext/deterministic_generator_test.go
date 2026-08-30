package socialcontext

import (
	"context"
	"reflect"
	"strings"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorProducesDerivedMeaningAndAuthorizedProvenance(t *testing.T) {
	generator := NewDeterministicGenerator()
	input := appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-1", Content: "最近開始深入研究分散式系統設計"},
		{ID: "activity-2", Content: "持續比較不同一致性模型的工程取捨"},
	}}

	got, err := generator.Generate(context.Background(), input)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning == "" {
		t.Fatal("Generate() meaning is blank")
	}
	for _, activity := range input.Activities {
		if got.Meaning == activity.Content {
			t.Fatalf("Generate() replayed raw Activity content %q", activity.Content)
		}
	}
	if !strings.Contains(got.Meaning, "分散式系統") || !strings.Contains(got.Meaning, "一致性模型") {
		t.Fatalf("Generate() meaning = %q, want meaning grounded in approved signal content", got.Meaning)
	}
	if want := []string{"activity-1", "activity-2"}; !reflect.DeepEqual(got.Provenance, want) {
		t.Fatalf("Generate() provenance = %#v, want %#v", got.Provenance, want)
	}
}

func TestDeterministicGeneratorDistinguishesDifferentSignals(t *testing.T) {
	generator := NewDeterministicGenerator()
	first, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-db", Content: "最近開始深入研究分散式系統設計"},
	}})
	if err != nil {
		t.Fatalf("Generate(database) error = %v", err)
	}
	second, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-run", Content: "最近開始準備第一次全程馬拉松訓練"},
	}})
	if err != nil {
		t.Fatalf("Generate(marathon) error = %v", err)
	}
	if first.Meaning == second.Meaning {
		t.Fatalf("different signals produced identical meaning %q", first.Meaning)
	}
}

func TestDeterministicGeneratorDeclinesUnmatchedSignalInsteadOfReplayingRawContent(t *testing.T) {
	generator := NewDeterministicGenerator()
	raw := "完成第一次全程馬拉松比賽"

	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-unmatched", Content: raw},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if strings.Contains(got.Meaning, raw) {
		t.Fatalf("Generate() meaning = %q, must not contain unmatched raw Activity content %q", got.Meaning, raw)
	}
	if got.Meaning != "" {
		t.Fatalf("Generate() meaning = %q, want blank meaning so candidate validation rejects an unsafe unmatched signal", got.Meaning)
	}
	if want := []string{"activity-unmatched"}; !reflect.DeepEqual(got.Provenance, want) {
		t.Fatalf("Generate() provenance = %#v, want %#v", got.Provenance, want)
	}
}

func TestDeterministicGeneratorAbstractsSingleSignalBeyondLightParaphrase(t *testing.T) {
	generator := NewDeterministicGenerator()
	raw := "最近開始深入研究分散式系統設計"

	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-db", Content: raw},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning == "" || !strings.Contains(got.Meaning, "分散式系統") || !strings.Contains(got.Meaning, "可靠性") || !strings.Contains(got.Meaning, "工程取捨") {
		t.Fatalf("Generate() meaning = %q, want a higher-level social meaning grounded in the signal", got.Meaning)
	}
	if strings.Contains(got.Meaning, "深入研究分散式系統設計") {
		t.Fatalf("Generate() meaning = %q, must not retain the source action phrase", got.Meaning)
	}
}

func TestDeterministicGeneratorKeepsRecognizedTopicsWhenAnotherSignalIsUnsupported(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-unsupported", Content: "完成第一次全程馬拉松比賽"},
		{ID: "activity-db", Content: "最近開始深入研究分散式系統設計"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning == "" || !strings.Contains(got.Meaning, "分散式系統") {
		t.Fatalf("Generate() meaning = %q, want recognized topic preserved despite unsupported signal", got.Meaning)
	}
}

func TestDeterministicGeneratorCanonicalizesTopicOrder(t *testing.T) {
	generator := NewDeterministicGenerator()
	firstInput := appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-db", Content: "最近開始深入研究分散式系統設計"},
		{ID: "activity-run", Content: "最近開始準備第一次全程馬拉松訓練"},
	}}
	reversedInput := appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-run", Content: "最近開始準備第一次全程馬拉松訓練"},
		{ID: "activity-db", Content: "最近開始深入研究分散式系統設計"},
	}}

	first, err := generator.Generate(context.Background(), firstInput)
	if err != nil {
		t.Fatalf("Generate(first) error = %v", err)
	}
	reversed, err := generator.Generate(context.Background(), reversedInput)
	if err != nil {
		t.Fatalf("Generate(reversed) error = %v", err)
	}
	if first.Meaning != reversed.Meaning {
		t.Fatalf("meanings differ by input order: first=%q reversed=%q", first.Meaning, reversed.Meaning)
	}
}

func TestDeterministicGeneratorRejectsNegatedOrContrastiveKeywordMatches(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, raw := range []string{
		"我不研究分散式系統",
		"完成第一次全程馬拉松比賽，沒有準備",
	} {
		t.Run(raw, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
				{ID: "activity-negated", Content: raw},
			}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if got.Meaning != "" {
				t.Fatalf("Generate() meaning = %q for negated/contrastive signal %q, want blank meaning", got.Meaning, raw)
			}
		})
	}
}
