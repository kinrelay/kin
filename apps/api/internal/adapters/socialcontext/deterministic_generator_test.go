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
	first, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-db", Content: "最近開始深入研究分散式系統設計"}}})
	if err != nil {
		t.Fatalf("Generate(database) error = %v", err)
	}
	second, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-run", Content: "最近開始準備第一次全程馬拉松訓練"}}})
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
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-unmatched", Content: raw}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if strings.Contains(got.Meaning, raw) {
		t.Fatalf("Generate() meaning = %q, must not contain unmatched raw Activity content %q", got.Meaning, raw)
	}
	if got.Meaning != "" {
		t.Fatalf("Generate() meaning = %q, want blank meaning so candidate validation rejects an unsafe unmatched signal", got.Meaning)
	}
	if want := []string(nil); !reflect.DeepEqual(got.Provenance, want) {
		t.Fatalf("Generate() provenance = %#v, want no provenance for unsupported signal", got.Provenance)
	}
}

func TestDeterministicGeneratorAbstractsSingleSignalBeyondLightParaphrase(t *testing.T) {
	generator := NewDeterministicGenerator()
	raw := "最近開始深入研究分散式系統設計"
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-db", Content: raw}}})
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
	if want := []string{"activity-db"}; !reflect.DeepEqual(got.Provenance, want) {
		t.Fatalf("Generate() provenance = %#v, want only activities that contributed to meaning %#v", got.Provenance, want)
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
		"最近已經不再研究分散式系統了",
		"我不再深入研究分散式系統",
		"最近已經不想研究分散式系統了",
		"完成第一次全程馬拉松比賽，沒有準備",
	} {
		t.Run(raw, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-negated", Content: raw}}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if got.Meaning != "" {
				t.Fatalf("Generate() meaning = %q for negated/contrastive signal %q, want blank meaning", got.Meaning, raw)
			}
			if len(got.Provenance) != 0 {
				t.Fatalf("Generate() provenance = %#v for negated/contrastive signal %q, want none", got.Provenance, raw)
			}
		})
	}
}

func TestDeterministicGeneratorRejectsMarathonPreparationWithoutParticipationIntent(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, raw := range []string{
		"開始準備放棄馬拉松訓練",
		"開始準備馬拉松賽事的志工物資",
		"開始準備馬拉松比賽的志工物資",
	} {
		t.Run(raw, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-run", Content: raw}}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if got.Meaning != "" || len(got.Provenance) != 0 {
				t.Fatalf("Generate() = %#v for non-participation preparation %q, want blank meaning and provenance", got, raw)
			}
		})
	}
}

func TestDeterministicGeneratorScopesIntentToIndividualClauses(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, raw := range []string{"最近開始準備第一次全程馬拉松訓練，目前沒有受傷", "停止熬夜，開始準備馬拉松"} {
		t.Run(raw, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-run", Content: raw}}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if got.Meaning == "" || !strings.Contains(got.Meaning, "耐力運動") {
				t.Fatalf("Generate() meaning = %q for affirmative marathon clause %q, want derived endurance context", got.Meaning, raw)
			}
		})
	}
}

func TestDeterministicGeneratorRejectsLaterReversalOfRecognizedIntent(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, raw := range []string{
		"最近開始研究分散式系統，但後來不再研究",
		"最近開始準備馬拉松，但後來沒有準備",
		"最近開始研究分散式系統但後來不再研究",
		"最近開始準備馬拉松但後來沒有準備",
		"最近開始準備馬拉松訓練但後來放棄了",
	} {
		t.Run(raw, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-reversed", Content: raw}}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if got.Meaning != "" || len(got.Provenance) != 0 {
				t.Fatalf("Generate() = %#v for reversed intent %q, want blank meaning and provenance", got, raw)
			}
		})
	}
}

func TestDeterministicGeneratorContinuesAfterReversedClauseToIndependentSafeTopic(t *testing.T) {
	generator := NewDeterministicGenerator()
	raw := "最近開始研究分散式系統，但後來不再研究，最近開始準備馬拉松"
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-mixed", Content: raw}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning == "" || !strings.Contains(got.Meaning, "耐力運動") || strings.Contains(got.Meaning, "分散式系統") {
		t.Fatalf("Generate() meaning = %q, want only the later independent safe marathon topic", got.Meaning)
	}
}

func TestDeterministicGeneratorRejectsExplicitMarathonNonParticipationSuffix(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, raw := range []string{
		"最近開始準備馬拉松比賽但後來不參加",
		"最近開始準備馬拉松比賽，但後來不參賽",
	} {
		t.Run(raw, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-run", Content: raw}}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if got.Meaning != "" || len(got.Provenance) != 0 {
				t.Fatalf("Generate() = %#v for explicit non-participation %q, want blank meaning and provenance", got, raw)
			}
		})
	}
}

func TestDeterministicGeneratorDoesNotBindUnrelatedReversalToRecognizedTopic(t *testing.T) {
	generator := NewDeterministicGenerator()
	for _, raw := range []string{
		"最近開始研究分散式系統，後來停止研究英文",
		"最近開始研究分散式系統但後來停止研究英文",
	} {
		t.Run(raw, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-db", Content: raw}}})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if got.Meaning == "" || !strings.Contains(got.Meaning, "分散式系統") {
				t.Fatalf("Generate() meaning = %q, want distributed-systems topic preserved when reversal targets English in %q", got.Meaning, raw)
			}
		})
	}
}

func TestDeterministicGeneratorLimitsUnpunctuatedReversalToItsObject(t *testing.T) {
	generator := NewDeterministicGenerator()
	raw := "最近開始研究分散式系統但後來停止研究英文並持續比較一致性模型"
	got, err := generator.Generate(context.Background(), appsc.ContextGenerationInput{Activities: []appsc.ContextGenerationActivity{{ID: "activity-db", Content: raw}}})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if got.Meaning == "" || !strings.Contains(got.Meaning, "分散式系統") {
		t.Fatalf("Generate() meaning = %q, want supported distributed-systems context preserved when reversal object is English", got.Meaning)
	}
}
