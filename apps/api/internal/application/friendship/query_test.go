package friendship

import (
	"context"
	"errors"
	"testing"

	domainidentity "github.com/kinrelay/kin/apps/api/internal/domain/identity"
)

type friendshipReadRepositoryFake struct {
	result FriendshipReadModel
	found  bool
	err    error
	calls  int
}

func (f *friendshipReadRepositoryFake) FindActiveBetween(_ context.Context, _, _ domainidentity.ID) (FriendshipReadModel, bool, error) {
	f.calls++
	return f.result, f.found, f.err
}

func TestGetFriendshipReturnsSameActiveRelationshipForEitherParticipant(t *testing.T) {
	for _, requester := range []string{"user-a", "user-b"} {
		t.Run(requester, func(t *testing.T) {
			repository := &friendshipReadRepositoryFake{
				result: FriendshipReadModel{FirstParticipantID: "user-a", SecondParticipantID: "user-b", Active: true},
				found:  true,
			}
			query := NewGetFriendship(repository)

			got, found, err := query.Execute(context.Background(), GetFriendshipQuery{
				RequesterID:         requester,
				FirstParticipantID:  "user-a",
				SecondParticipantID: "user-b",
			})
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if !found {
				t.Fatal("Execute() found = false, want true")
			}
			if !got.Active || got.FirstParticipantID != "user-a" || got.SecondParticipantID != "user-b" {
				t.Fatalf("Execute() = %#v, want active relationship user-a/user-b", got)
			}
		})
	}
}

func TestGetFriendshipDoesNotExposePendingRelationship(t *testing.T) {
	repository := &friendshipReadRepositoryFake{found: false}
	query := NewGetFriendship(repository)

	_, found, err := query.Execute(context.Background(), GetFriendshipQuery{
		RequesterID:         "user-a",
		FirstParticipantID:  "user-a",
		SecondParticipantID: "user-b",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if found {
		t.Fatal("Execute() found = true, want false for pending/not-active friendship")
	}
}

func TestGetFriendshipRejectsNonParticipantBeforeReadingPrivateState(t *testing.T) {
	repository := &friendshipReadRepositoryFake{found: true}
	query := NewGetFriendship(repository)

	_, _, err := query.Execute(context.Background(), GetFriendshipQuery{
		RequesterID:         "user-c",
		FirstParticipantID:  "user-a",
		SecondParticipantID: "user-b",
	})
	if !errors.Is(err, ErrFriendshipQueryUnauthorized) {
		t.Fatalf("Execute() error = %v, want %v", err, ErrFriendshipQueryUnauthorized)
	}
	if repository.calls != 0 {
		t.Fatalf("repository calls = %d, want 0 for unauthorized requester", repository.calls)
	}
}
