package socialcontext

import "testing"

func TestContextCandidateRequiresMeaningAndProvenance(t *testing.T) {
	tests := []struct {
		name       string
		meaning    string
		provenance []string
		wantErr    error
	}{
		{name: "blank meaning", meaning: "   ", provenance: []string{"activity-1"}, wantErr: ErrBlankContextMeaning},
		{name: "missing provenance", meaning: "最近開始深入研究分散式系統設計", provenance: nil, wantErr: ErrMissingContextProvenance},
		{name: "blank provenance id", meaning: "最近開始深入研究分散式系統設計", provenance: []string{" "}, wantErr: ErrMissingContextProvenance},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewContextCandidate(tt.meaning, tt.provenance)
			if err != tt.wantErr {
				t.Fatalf("NewContextCandidate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestPromoteContextCandidateRejectsPureSourceReplay(t *testing.T) {
	candidate, err := NewContextCandidate("最近開始深入研究分散式系統設計", []string{"activity-1"})
	if err != nil {
		t.Fatalf("NewContextCandidate() error = %v", err)
	}

	_, err = PromoteContextCandidate(candidate, []SourceActivity{{ID: "activity-1", Content: "  最近開始深入研究分散式系統設計  "}})
	if err != ErrSourceReplay {
		t.Fatalf("PromoteContextCandidate() error = %v, want %v", err, ErrSourceReplay)
	}
}

func TestPromoteContextCandidateRejectsReplayOfAuthorizedSourceOutsideDeclaredProvenance(t *testing.T) {
	candidate, err := NewContextCandidate("整理了一份一致性模型筆記", []string{"activity-1"})
	if err != nil {
		t.Fatalf("NewContextCandidate() error = %v", err)
	}

	_, err = PromoteContextCandidate(candidate, []SourceActivity{
		{ID: "activity-1", Content: "讀完一篇關於 Raft 的文章"},
		{ID: "activity-2", Content: "整理了一份一致性模型筆記"},
	})
	if err != ErrSourceReplay {
		t.Fatalf("PromoteContextCandidate() error = %v, want %v", err, ErrSourceReplay)
	}
}

func TestPromoteContextCandidateRejectsZeroValueCandidate(t *testing.T) {
	_, err := PromoteContextCandidate(ContextCandidate{}, []SourceActivity{{ID: "activity-1", Content: "讀完一篇關於 Raft 的文章"}})
	if err != ErrBlankContextMeaning {
		t.Fatalf("PromoteContextCandidate() error = %v, want %v", err, ErrBlankContextMeaning)
	}
}

func TestPromoteContextCandidateCreatesPrivateSocialContext(t *testing.T) {
	candidate, err := NewContextCandidate("最近對分散式系統的可靠性與取捨特別有興趣", []string{"activity-1", "activity-2"})
	if err != nil {
		t.Fatalf("NewContextCandidate() error = %v", err)
	}

	context, err := PromoteContextCandidate(candidate, []SourceActivity{
		{ID: "activity-1", Content: "讀完一篇關於 Raft 的文章"},
		{ID: "activity-2", Content: "整理了一份一致性模型筆記"},
	})
	if err != nil {
		t.Fatalf("PromoteContextCandidate() error = %v", err)
	}
	if got, want := context.Meaning(), "最近對分散式系統的可靠性與取捨特別有興趣"; got != want {
		t.Fatalf("Meaning() = %q, want %q", got, want)
	}
	if !context.IsPrivate() {
		t.Fatal("promoted SocialContext must remain owner-private in MVP 2")
	}
	if got := context.Provenance(); len(got) != 2 || got[0] != "activity-1" || got[1] != "activity-2" {
		t.Fatalf("Provenance() = %#v, want source activity ids", got)
	}
}

func TestPromoteContextCandidateRejectsMissingAuthorizedSource(t *testing.T) {
	candidate, err := NewContextCandidate("最近對分散式系統的可靠性與取捨特別有興趣", []string{"activity-1", "activity-2"})
	if err != nil {
		t.Fatalf("NewContextCandidate() error = %v", err)
	}

	_, err = PromoteContextCandidate(candidate, []SourceActivity{{ID: "activity-1", Content: "讀完一篇關於 Raft 的文章"}})
	if err != ErrMissingContextProvenance {
		t.Fatalf("PromoteContextCandidate() error = %v, want %v", err, ErrMissingContextProvenance)
	}
}
