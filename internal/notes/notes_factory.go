package notes

import (
	"fmt"
	"time"

	"github.com/germandv/apio/internal/id"
	"github.com/germandv/apio/internal/tools"
)

// FromReq produces a NoteAggregate from a Request.
func FromReq(
	title string,
	content string,
	tagIDs []string,
	createdBy string,
	createdAt time.Time,
) (NoteAggregate, error) {
	noteTitle, err := ParseTitle(title)
	if err != nil {
		return NoteAggregate{}, err
	}

	noteContent, err := ParseContent(content)
	if err != nil {
		return NoteAggregate{}, err
	}

	creatorID, err := id.Parse(createdBy)
	if err != nil {
		return NoteAggregate{}, err
	}

	uid := id.New()

	noteEntity := NoteEntity{
		ID:        uid,
		Title:     noteTitle,
		Content:   noteContent,
		CreatedBy: creatorID,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	noteTags := make([]NoteTagEntity, 0, len(tagIDs))
	for _, t := range tagIDs {
		tagID, err := id.Parse(t)
		if err != nil {
			return NoteAggregate{}, err
		}
		noteTags = append(noteTags, NoteTagEntity{ID: tagID, Name: ""})
	}

	n := NoteAggregate{
		NoteEntity: noteEntity,
		Tags:       noteTags,
	}

	return n, nil
}

// FromDB produces a NoteAggregate from data in the DB format.
func FromDB(
	dbID string,
	dbTitle string,
	dbContent string,
	dbCreatedBy string,
	dbCreatedAt time.Time,
	dbUpdatedAt time.Time,
	dbTags []struct {
		ID   string
		Name string
	},
) (NoteAggregate, error) {
	uid, err := id.Parse(dbID)
	if err != nil {
		return NoteAggregate{}, err
	}

	title, err := ParseTitle(dbTitle)
	if err != nil {
		return NoteAggregate{}, err
	}

	content, err := ParseContent(dbContent)
	if err != nil {
		return NoteAggregate{}, err
	}

	creatorID, err := id.Parse(dbCreatedBy)
	if err != nil {
		return NoteAggregate{}, err
	}

	noteEntity := NoteEntity{
		ID:        uid,
		Title:     title,
		Content:   content,
		CreatedBy: creatorID,
		CreatedAt: dbCreatedAt,
		UpdatedAt: dbUpdatedAt,
	}

	corruptedTags := make([]string, 0)

	noteTags := tools.Map(dbTags, func(t struct {
		ID   string
		Name string
	}) NoteTagEntity {
		uid, err := id.Parse(t.ID)
		if err != nil {
			corruptedTags = append(corruptedTags, fmt.Sprintf("tag %v is invalid", t))
			return NoteTagEntity{ID: id.Zero(), Name: ""}
		}
		return NoteTagEntity{ID: uid, Name: t.Name}
	})

	if len(corruptedTags) > 0 {
		return NoteAggregate{}, fmt.Errorf("corrupted data, could not parse one or more tags %v", corruptedTags)
	}

	n := NoteAggregate{
		NoteEntity: noteEntity,
		Tags:       noteTags,
	}

	return n, nil
}
