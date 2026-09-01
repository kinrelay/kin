package socialcontext

import (
	"context"
	"reflect"
	"testing"
	"time"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func TestDeriveContextCandidateUsesDeterministicTotalChronologyWhenTimestampsTie(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	occurred := time.Date(2026, 8, 31, 1, 0, 0, 0, time.UTC)
	contributed := time.Date(2026, 8, 31, 1, 1, 0, 0, time.UTC)
	activities := []ActivityForContext{
		{ID: "z-stop", OwnerID: ownerID, Content: "不再研究分散式系統", OccurredAt: occurred, ContributedAt: contributed},
		{ID: "a-start", OwnerID: ownerID, Content: "最近開始研究分散式系統與可靠性工程取捨", OccurredAt: occurred, ContributedAt: contributed},
	}

	got := deriveGeneratorInputIDs(t, activities)
	if want := []string{"a-start", "z-stop"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("generator input IDs = %#v, want deterministic total chronology %#v", got, want)
	}

	reversed := []ActivityForContext{activities[1], activities[0]}
	gotReversed := deriveGeneratorInputIDs(t, reversed)
	if !reflect.DeepEqual(gotReversed, got) {
		t.Fatalf("reverse reader order changed chronology: first=%#v reversed=%#v", got, gotReversed)
	}
}

func TestDeriveContextCandidateChronologyHandlesMissingTimestampsTransitively(t *testing.T) {
	ownerID, _ := domainidentity.NewID("owner-1")
	occurred := time.Date(2026, 8, 31, 2, 0, 0, 0, time.UTC)
	activities := []ActivityForContext{
		{ID: "c-known", OwnerID: ownerID, Content: "最近開始研究分散式系統與可靠性工程取捨", OccurredAt: occurred},
		{ID: "b-contributed", OwnerID: ownerID, Content: "最近開始準備第一次全程馬拉松訓練", ContributedAt: occurred},
		{ID: "a-unknown", OwnerID: ownerID, Content: "持續比較不同一致性模型的工程取捨"},
	}

	got := deriveGeneratorInputIDs(t, activities)
	if want := []string{"a-unknown", "b-contributed", "c-known"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("generator input IDs = %#v, want explicit availability order %#v", got, want)
	}
}

func deriveGeneratorInputIDs(t *testing.T, activities []ActivityForContext) []string {
	t.Helper()
	generator := &fakeContextGenerator{}
	uc := NewDeriveContextCandidate(fakeContextActivityReader{activities: activities}, generator, &fakeSocialContextRepository{})
	ids := make([]string, 0, len(activities))
	for _, activity := range activities {
		ids = append(ids, activity.ID)
	}
	_, err := uc.Execute(context.Background(), DeriveContextCandidateCommand{RequesterID: "owner-1", ActivityIDs: ids})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	got := make([]string, 0, len(generator.input.Activities))
	for _, activity := range generator.input.Activities {
		got = append(got, activity.ID)
	}
	return got
}
