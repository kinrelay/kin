package activity

import (
	"context"
	"reflect"
	"testing"
	"time"

	applicationsocialcontext "github.com/kinrelay/kin/apps/api/internal/application/socialcontext"
	domainactivity "github.com/kinrelay/kin/apps/api/internal/domain/activity"
	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

func TestMemoryReadRepositoryListsOnlyRequestedPrivateNormalizedActivitiesForContext(t *testing.T) {
	ctx := context.Background()
	writeRepository := NewMemoryRepository()
	readRepository := NewMemoryReadRepository(writeRepository)
	aliceID, err := domainidentity.NewID("alice")
	if err != nil {
		t.Fatalf("NewID(alice): %v", err)
	}
	baseTime := time.Date(2026, time.August, 30, 10, 0, 0, 0, time.UTC)

	makeActivity := func(id, owner, content string) domainactivity.Activity {
		ownerID, err := domainidentity.NewID(owner)
		if err != nil {
			t.Fatalf("NewID(%q): %v", owner, err)
		}
		normalized, err := domainactivity.NewContent(content)
		if err != nil {
			t.Fatalf("NewContent(%q): %v", content, err)
		}
		value, err := domainactivity.NewManual(id, ownerID, normalized, baseTime, baseTime)
		if err != nil {
			t.Fatalf("NewManual(%q): %v", id, err)
		}
		return value
	}

	for _, value := range []domainactivity.Activity{
		makeActivity("alice-first", "alice", "最近開始深入研究分散式系統設計"),
		makeActivity("alice-second", "alice", "持續比較不同一致性模型的工程取捨"),
		makeActivity("alice-unrequested", "alice", "另一筆不在本次 derivation request 的活動"),
		makeActivity("bob-requested", "bob", "另一位使用者的 private activity"),
	} {
		if err := writeRepository.Save(ctx, value); err != nil {
			t.Fatalf("Save(%q): %v", value.ID(), err)
		}
	}

	got, err := readRepository.ListOwnerPrivateNormalized(ctx, aliceID, []string{"alice-second", "alice-first", "alice-second", "bob-requested"})
	if err != nil {
		t.Fatalf("ListOwnerPrivateNormalized() error = %v", err)
	}
	want := []applicationsocialcontext.ActivityForContext{
		{ID: "alice-second", OwnerID: aliceID, Content: "持續比較不同一致性模型的工程取捨"},
		{ID: "alice-first", OwnerID: aliceID, Content: "最近開始深入研究分散式系統設計"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListOwnerPrivateNormalized() = %#v, want %#v", got, want)
	}
}
