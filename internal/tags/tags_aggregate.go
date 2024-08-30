package tags

import "github.com/germandv/apio/internal/id"

// TagEntity is an Entity that represents a Tag within the domain.
type TagEntity struct {
	ID   id.ID
	Name Name
}

// TagAggregate is an Aggregate that represents a Tag to the outside world.
type TagAggregate struct {
	TagEntity
	NoteCount int
}
