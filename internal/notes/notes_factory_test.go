package notes

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/germandv/apio/internal/errs"
	"github.com/germandv/apio/internal/id"
	"github.com/germandv/apio/internal/tools"
)

func TestFromReq(t *testing.T) {
	content := tools.RandomStr(100)
	createdBy := id.New().String()

	t.Run("no_title", func(t *testing.T) {
		_, err := FromReq("", content, []string{}, createdBy, time.Now())
		if err == nil {
			t.Error("expected an error, got nil")
		}
		if !errors.Is(err, errs.ErrEmptyTitle) {
			t.Errorf("expected ErrEmptyTitle, got %v", err)
		}
	})

	t.Run("title_too_short", func(t *testing.T) {
		_, err := FromReq(tools.RandomStr(minTitleLength-1), content, []string{}, createdBy, time.Now())
		if err == nil {
			t.Error("expected an error, got nil")
		}
		if !errors.Is(err, errs.ErrMinLen) {
			t.Errorf("expected ErrMinLen, got %v", err)
		}
	})

	t.Run("title_too_long", func(t *testing.T) {
		_, err := FromReq(tools.RandomStr(maxTitleLength+1), content, []string{}, createdBy, time.Now())
		if err == nil {
			t.Error("expected an error, got nil")
		}
		if !errors.Is(err, errs.ErrMaxLen) {
			t.Errorf("expected ErrMaxLen, got %v", err)
		}
	})

	t.Run("valid", func(t *testing.T) {
		title := tools.RandomStr(minTitleLength + 4)
		n, err := FromReq(title, content, []string{}, createdBy, time.Now())
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if n.Title.String() != title {
			t.Errorf("expected note title to be %q, got %q", title, n.Title.String())
		}
		if n.Content.String() != content {
			t.Errorf("expected note content to be %q, got %q", content, n.Content.String())
		}
		if len(n.Tags) != 0 {
			t.Errorf("expected note to have 0 tags, got %d", len(n.Tags))
		}
	})
}

func TestFromDB(t *testing.T) {
	createdBy := id.New().String()

	t.Run("tag_with_invalid_id", func(t *testing.T) {
		dbTags := []struct {
			ID   string
			Name string
		}{
			{ID: id.New().String(), Name: tools.RandomStr(4)},
			{ID: "invalid-id", Name: tools.RandomStr(4)},
		}

		_, err := FromDB(id.New().String(), "title", "content", createdBy, time.Now(), time.Now(), dbTags)
		if err == nil {
			t.Error("expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "could not parse one or more tags") {
			t.Errorf("expected error msg to contain could not parse one or more tags, got: %s", err)
		}
	})

	t.Run("valid", func(t *testing.T) {
		dbTags := []struct {
			ID   string
			Name string
		}{
			{ID: id.New().String(), Name: tools.RandomStr(4)},
			{ID: id.New().String(), Name: tools.RandomStr(4)},
		}

		n, err := FromDB(id.New().String(), "title", "content", createdBy, time.Now(), time.Now(), dbTags)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if n.Title.String() != "title" {
			t.Errorf("expected note title to be title, got %q", n.Title.String())
		}
		if len(n.Tags) != 2 {
			t.Errorf("expected note to have 2 tags, got %d", len(n.Tags))
		}
	})
}
