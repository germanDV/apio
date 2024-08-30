package tags_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/germandv/apio/internal/errs"
	"github.com/germandv/apio/internal/memorydb"
	"github.com/germandv/apio/internal/tags"
)

func TestTagService(t *testing.T) {
	repo := memorydb.NewTagsRepository()
	svc := tags.NewService(repo)

	t.Run("name_missing", func(t *testing.T) {
		t.Parallel()
		_, err := svc.Create("")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, errs.ErrEmptyName) {
			t.Fatalf("expected error ErrEmptyName, got %s", err)
		}
	})

	t.Run("name_too_long", func(t *testing.T) {
		t.Parallel()
		_, err := svc.Create("Arnold Schwarzenegger Arnold Schwarzenegger Arnold Schwarzenegger Arnold Schwarzenegger")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, errs.ErrMaxLen) {
			t.Fatalf("expected error ErrMaxLen, got %s", err)
		}
	})

	t.Run("create_and_retrieve_tags", func(t *testing.T) {
		t.Parallel()

		ts, err := svc.GetAll()
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
		if len(ts) != 0 {
			t.Fatalf("expected tags to be an empty slice, got %v", ts)
		}

		_, err = svc.Create("meditation")
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		_, err = svc.Create("music")
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		ts, err = svc.GetAll()
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		hasMeditationTag := slices.ContainsFunc(ts, func(e tags.TagAggregate) bool {
			return e.Name.String() == "meditation"
		})
		if !hasMeditationTag {
			t.Fatalf("missing meditation tag, got %v", ts)
		}

		hasMusicTag := slices.ContainsFunc(ts, func(e tags.TagAggregate) bool {
			return e.Name.String() == "music"
		})
		if !hasMusicTag {
			t.Fatalf("missing music tag, got %v", ts)
		}

		_, err = svc.Create("music")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, errs.ErrDuplicateTag) {
			t.Fatalf("expected error ErrDuplicateTag, got %s", err)
		}
	})
}
