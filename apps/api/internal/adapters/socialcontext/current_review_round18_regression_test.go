package socialcontext

import (
	"context"
	"testing"

	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorIgnoresPluralOtherSubjectAfterCompoundConnector(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, subject := range []string{"我的朋友們", "朋友們", "同事們", "家人們"} {
		generated, err := generator.Generate(context.Background(), applicationsocialcontext.ContextGenerationInput{
			Activities: []applicationsocialcontext.ContextGenerationActivity{
				{ID: "activity-compound", Content: "開始準備馬拉松但" + subject + "放棄馬拉松"},
			},
		})
		if err != nil {
			t.Fatalf("Generate(%q) error = %v", subject, err)
		}
		wantMeaning := "近期關注" + marathonTopic
		if generated.Meaning != wantMeaning {
			t.Fatalf("Generate(%q).Meaning = %q, want %q", subject, generated.Meaning, wantMeaning)
		}
		if len(generated.Provenance) != 1 || generated.Provenance[0] != "activity-compound" {
			t.Fatalf("Generate(%q).Provenance = %#v, want [activity-compound]", subject, generated.Provenance)
		}
	}
}
