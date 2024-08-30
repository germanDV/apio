package tags

import (
	"fmt"
	"unicode/utf8"

	"github.com/germandv/apio/internal/errs"
)

const (
	maxNameLength = 64
)

// Name is a Value Object that represents the name of the Tag.
type Name struct {
	value string
}

// ParseName creates a Name Value Object ensuring it is valid.
func ParseName(name string) (Name, error) {
	if name == "" {
		return Name{}, errs.ErrEmptyName
	}
	if utf8.RuneCountInString(name) > maxNameLength {
		return Name{}, fmt.Errorf("%w: name must be less than %d chars", errs.ErrMaxLen, maxNameLength)
	}
	return Name{value: name}, nil
}

func (name Name) String() string {
	return name.value
}
