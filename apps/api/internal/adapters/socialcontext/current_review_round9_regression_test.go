package socialcontext

import (
	"context"
	"strings"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorBindsTrailedBareAbandonmentToNearestTopic(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-distributed", Content: "最近開始研究分散式系統"},
		{ID: "activity-marathon", Content: "最近開始準備馬拉松"},
		{ID: "activity-abandon", Content: "後來放棄了但最近工作真的很忙"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "分散式系統") {
		t.Fatalf("Generate() meaning = %q, want trailed bare abandonment to preserve the older distributed-systems topic", got.Meaning)
	}
	if strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want trailed bare abandonment to retract only the nearest marathon topic", got.Meaning)
	}
}
