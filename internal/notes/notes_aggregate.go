package notes

import (
	"time"

	"github.com/germandv/apio/internal/id"
)

// NoteEntity is an Entity that represents a Note within the domain.
type NoteEntity struct {
	ID        id.ID
	Title     Title
	Content   Content
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NoteEntity is an Entity that represents Tags associated to a Note.
type NoteTagEntity struct {
	ID   id.ID
	Name string
}

// NoteAggregate is an Aggregate that represents a Note to the outside world.
type NoteAggregate struct {
	NoteEntity
	Tags []NoteTagEntity
}
