package socialcontext

import (
	"context"
	"reflect"
	"strings"
	"testing"

	appsc "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
)

func TestDeterministicGeneratorRejectsModifierBearingDistributedSystemsReversal(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-db",
		Content: "最近開始研究分散式系統，但後來不再深入研究分散式系統",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning != "" || len(got.Provenance) != 0 {
		t.Fatalf("Generate() = %#v, want modifier-bearing reversal to retract distributed-systems topic", got)
	}
}

func TestDeterministicGeneratorRecognizesDistributedSystemsAbandonment(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{
		{ID: "activity-start", Content: "最近開始研究分散式系統與可靠性工程取捨"},
		{ID: "activity-stop", Content: "放棄研究分散式系統"},
	}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning != "" || len(got.Provenance) != 0 {
		t.Fatalf("Generate() = %#v, want abandonment to retract distributed-systems topic", got)
	}
}

func TestDeterministicGeneratorDoesNotTreatNounPrefixAsActionBoundary(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-db",
		Content: "最近開始研究分散式系統，後來停止研究並購策略",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning == "" || !strings.Contains(got.Meaning, "分散式系統") {
		t.Fatalf("Generate() meaning = %q, want unrelated M&A-strategy reversal to preserve distributed-systems topic", got.Meaning)
	}
}

func TestDeterministicGeneratorAccumulatesSupportedActionsWithinCompoundClause(t *testing.T) {
	generator := NewDeterministicGenerator()
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{
		ID:      "activity-multi",
		Content: "最近開始研究分散式系統並開始準備馬拉松",
	}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !strings.Contains(got.Meaning, "分散式系統") || !strings.Contains(got.Meaning, "耐力運動") {
		t.Fatalf("Generate() meaning = %q, want both supported topics from compound clause", got.Meaning)
	}
	if want := []string{"activity-multi"}; !reflect.DeepEqual(got.Provenance, want) {
		t.Fatalf("Generate() provenance = %#v, want %#v", got.Provenance, want)
	}
}
