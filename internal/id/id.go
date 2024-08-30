package id

import (
	"fmt"

	"github.com/germandv/apio/internal/errs"
	"github.com/google/uuid"
)

// ID is a Value Object that represents the ID of the Tag.
type ID struct {
	value string
}

// New creates a new ID Value Object.
func New() ID {
	id, _ := uuid.NewV7()
	return ID{value: id.String()}
}

// Parse creates an ID Value Object ensuring it is valid.
func Parse(id string) (ID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ID{}, fmt.Errorf("error parsing id %s: %w", id, errs.ErrInvalidID)
	}
	return ID{value: parsedID.String()}, nil
}

func (id ID) String() string {
	return id.value
}

// Zero is the zero-value representation of an ID.
func Zero() ID {
	return ID{value: uuid.Nil.String()}
}
