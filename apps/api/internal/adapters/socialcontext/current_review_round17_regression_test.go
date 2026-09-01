package socialcontext

import (
	"context"
	"testing"

	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorIgnoresAnotherSubjectAfterCompoundConnector(t *testing.T) {
	generator := NewDeterministicGenerator()
	generated, err := generator.Generate(context.Background(), applicationsocialcontext.ContextGenerationInput{
		Activities: []applicationsocialcontext.ContextGenerationActivity{
			{ID: "activity-compound", Content: "開始準備馬拉松但朋友放棄馬拉松"},
		},
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	wantMeaning := "近期關注" + marathonTopic
	if generated.Meaning != wantMeaning {
		t.Fatalf("Generate().Meaning = %q, want %q", generated.Meaning, wantMeaning)
	}
	if len(generated.Provenance) != 1 || generated.Provenance[0] != "activity-compound" {
		t.Fatalf("Generate().Provenance = %#v, want [activity-compound]", generated.Provenance)
	}
}
