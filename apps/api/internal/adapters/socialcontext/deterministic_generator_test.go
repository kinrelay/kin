package socialcontext

import (
	"context"
	"reflect"
	"testing"

	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorProducesDerivedMeaningAndAuthorizedProvenance(t *testing.T) {
	generator := NewDeterministicGenerator()
	input := applicationsocialcontext.ContextGenerationInput{Activities: []applicationsocialcontext.ContextGenerationActivity{
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
	if want := []string{"activity-1", "activity-2"}; !reflect.DeepEqual(got.Provenance, want) {
		t.Fatalf("Generate() provenance = %#v, want %#v", got.Provenance, want)
	}
}
