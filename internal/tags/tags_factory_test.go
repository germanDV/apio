package tags

import (
	"errors"
	"testing"

	"github.com/germandv/apio/internal/errs"
	"github.com/germandv/apio/internal/id"
	"github.com/germandv/apio/internal/tools"
)

func TestFromReq(t *testing.T) {
	t.Run("no_name", func(t *testing.T) {
		_, err := FromReq("")
		if err == nil {
			t.Error("expected an error, got nil")
		}
		if !errors.Is(err, errs.ErrEmptyName) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}
	})

	t.Run("name_too_long", func(t *testing.T) {
		_, err := FromReq(tools.RandomStr(maxNameLength + 1))
		if err == nil {
			t.Error("expected an error, got nil")
		}
		if !errors.Is(err, errs.ErrMaxLen) {
			t.Errorf("expected ErrMaxLen, got %v", err)
		}
	})

	t.Run("valid", func(t *testing.T) {
		input := tools.RandomStr(8)
		tag, err := FromReq(input)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if tag.Name.String() != input {
			t.Errorf("expected tag name to be %q, got %q", input, tag.Name.String())
		}
	})
}

func TestFromDB(t *testing.T) {
	goodID := id.New().String()

	t.Run("no_name", func(t *testing.T) {
		_, err := FromDB(goodID, "")
		if err == nil {
			t.Error("expected an error, got nil")
		}
		if !errors.Is(err, errs.ErrEmptyName) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}
	})

	t.Run("name_too_long", func(t *testing.T) {
		_, err := FromDB(goodID, tools.RandomStr(maxNameLength+1))
		if err == nil {
			t.Error("expected an error, got nil")
		}
		if !errors.Is(err, errs.ErrMaxLen) {
			t.Errorf("expected ErrMaxLen, got %v", err)
		}
	})

	t.Run("no_id", func(t *testing.T) {
		_, err := FromDB("", tools.RandomStr(6))
		if err == nil {
			t.Error("expected an error, got nil")
		}
		if !errors.Is(err, errs.ErrInvalidID) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}
	})

	t.Run("invalid_id", func(t *testing.T) {
		_, err := FromDB("not-a-uuidv7", tools.RandomStr(6))
		if err == nil {
			t.Error("expected an error, got nil")
		}
		if !errors.Is(err, errs.ErrInvalidID) {
			t.Errorf("expected ErrEmptyName, got %v", err)
		}
	})

	t.Run("valid", func(t *testing.T) {
		name := tools.RandomStr(8)
		tag, err := FromDB(goodID, name)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if tag.ID.String() != goodID {
			t.Errorf("expected tag ID to be %q, got %q", goodID, tag.ID.String())
		}
		if tag.Name.String() != name {
			t.Errorf("expected tag name to be %q, got %q", name, tag.Name.String())
		}
	})
}
