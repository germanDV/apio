package notes

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/germandv/apio/internal/errs"
)

const (
	minTitleLength   = 1
	maxTitleLength   = 64
	maxContentLength = 2_000
)

// Title is a Value Object that represents the Title of the Note.
type Title struct {
	value string
}

// ParseTitle creates a Title Value Object ensuring it respects the invariants.
func ParseTitle(title string) (Title, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Title{}, errs.ErrEmptyTitle
	}
	if utf8.RuneCountInString(title) < minTitleLength {
		return Title{}, fmt.Errorf("%w: title must be at least %d char(s)", errs.ErrMinLen, minTitleLength)
	}
	if utf8.RuneCountInString(title) > maxTitleLength {
		return Title{}, fmt.Errorf("%w: title must be less than %d chars", errs.ErrMaxLen, maxTitleLength)
	}
	return Title{value: title}, nil
}

func (title Title) String() string {
	return title.value
}

// Content is a Value Object that represents the Content of the Note.
type Content struct {
	value string
}

// ParseContent creates a Content Value Object ensuring it respects the invariants.
func ParseContent(content string) (Content, error) {
	if utf8.RuneCountInString(content) > maxContentLength {
		return Content{}, fmt.Errorf("%w: content must be less than %d chars", errs.ErrMaxLen, maxContentLength)
	}
	return Content{value: content}, nil
}

func (c Content) String() string {
	return c.value
}
