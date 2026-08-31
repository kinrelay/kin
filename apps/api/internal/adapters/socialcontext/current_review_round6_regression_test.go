package socialcontext

import (
	"context"
	"strings"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorRejectsDistributedSystemsRoleDutyIntent(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, content := range []string{
		"最近開始研究分散式系統研討會的志工工作",
		"最近開始研究分散式系統課程的助教工作",
	} {
		t.Run(content, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-1", Content: content}}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if got.Meaning != "" {
				t.Fatalf("Generate() meaning = %q, want role-duty logistics to stay outside distributed-systems interest", got.Meaning)
			}
		})
	}
}

func TestDeterministicGeneratorSplitsTemporalTransitionsBeforeSupportedActions(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, content := range []string{
		"停止研究分散式系統後開始準備馬拉松比賽",
		"停止研究分散式系統之後開始準備馬拉松比賽",
		"停止研究分散式系統然後開始準備馬拉松比賽",
	} {
		t.Run(content, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-1", Content: content}}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if !strings.Contains(got.Meaning, "耐力運動") || strings.Contains(got.Meaning, "分散式系統") {
				t.Fatalf("Generate() meaning = %q, want temporal transition to retain only the newer marathon action", got.Meaning)
			}
		})
	}
}
