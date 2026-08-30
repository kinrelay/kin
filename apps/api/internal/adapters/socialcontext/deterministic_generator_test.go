package socialcontext

import (
	"context"
	"reflect"
	"strings"
	"testing"

	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorProducesDerivedMeaningAndAuthorizedProvenance(t *testing.T) {
	generator := NewDeterministicGenerator()
	input := applicationsocialcontext.ContextGenerationInput{Activities: []applicationssocialcontext.ContextGenerationActivity{
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
	first, err := generator.Generate(context.Background(), applicationssocialcontext.ContextGenerationInput{Activities: []applicationssocialcontext.ContextGenerationActivity{
		{ID: "activity-db", Content: "最近開始深入研究分散式系統設計"},
	}})
	if err != nil {
		t.Fatalf("Generate(database) error = %v", err)
	}
	second, err := generator.Generate(context.Background(), applicationssocialcontext.ContextGenerationInput{Activities: []applicationssocialcontext.ContextGenerationActivity{
		{ID: "activity-run", Content: "最近開始準備第一次全程馬拉松訓練"},
	}})
	if err != nil {
		t.Fatalf("Generate(marathon) error = %v", err)
	}
	if first.Meaning == second.Meaning {
		t.Fatalf("different signals produced identical meaning %q", first.Meaning)
	}
}
