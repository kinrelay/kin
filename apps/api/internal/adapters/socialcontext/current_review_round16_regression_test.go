package socialcontext

import (
	"context"
	"testing"

	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorDoesNotApplyAnotherPersonsExplicitReversalToOwner(t *testing.T) {
	generator := NewDeterministicGenerator()
	generated, err := generator.Generate(context.Background(), applicationsocialcontext.ContextGenerationInput{
		Activities: []applicationsocialcontext.ContextGenerationActivity{
			{ID: "activity-marathon", Content: "最近開始準備第一次全程馬拉松比賽"},
			{ID: "activity-friend-abandon", Content: "朋友後來放棄了馬拉松比賽"},
		},
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	wantMeaning := "近期關注" + marathonTopic
	if generated.Meaning != wantMeaning {
		t.Fatalf("Generate().Meaning = %q, want %q", generated.Meaning, wantMeaning)
	}
	if len(generated.Provenance) != 1 || generated.Provenance[0] != "activity-marathon" {
		t.Fatalf("Generate().Provenance = %#v, want [activity-marathon]", generated.Provenance)
	}
}
