package memorydb

import (
	"fmt"

	"github.com/germandv/apio/internal/id"
	"github.com/germandv/apio/internal/notes"
	"github.com/germandv/apio/internal/tools"
)

// NotesRepository implements the Repository interface for Notes.
type NotesRepository struct {
	db *DB
}

// NewNotesRepository creates a NotesRepository that satisfies the Repository interface for Notes.
func NewNotesRepository() notes.Repository {
	return &NotesRepository{
		db: database,
	}
}

func (r *NotesRepository) Save(note notes.NoteAggregate) error {
	err := r.db.SaveNote(
		note.ID.String(),
		note.Title.String(),
		note.Content.String(),
		note.CreatedAt,
		note.UpdatedAt,
	)
	if err != nil {
		return err
	}

	ids := tools.Map(note.Tags, func(n notes.NoteTagEntity) string {
		return n.ID.String()
	})

	err = r.db.SaveNoteTags(note.ID.String(), ids)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotesRepository) List() ([]notes.NoteAggregate, error) {
	rows, err := r.db.GetNotes()
	if err != nil {
		return nil, err
	}

	entries := make([]notes.NoteAggregate, 0, len(rows))
	for _, row := range rows {
		tagRows, err := r.db.GetNoteTags(row.ID)
		if err != nil {
			return nil, err
		}

		corruptedTags := make([]TagRow, 0)
		tags := tools.Map(tagRows, func(t TagRow) notes.NoteTagEntity {
			uid, err := id.Parse(t.ID)
			if err != nil {
				corruptedTags = append(corruptedTags, t)
				return notes.NoteTagEntity{ID: id.Zero(), Name: ""}
			}
			return notes.NoteTagEntity{ID: uid, Name: t.Name}
		})
		if len(corruptedTags) > 0 {
			return nil, fmt.Errorf("corrupted data, could not parse one ore more tag row %v", corruptedTags)
		}

		uid, err := id.Parse(row.ID)
		if err != nil {
			return nil, err
		}

		title, err := notes.ParseTitle(row.Title)
		if err != nil {
			return nil, err
		}

		content, err := notes.ParseContent(row.Content)
		if err != nil {
			return nil, err
		}

		entries = append(entries, notes.NoteAggregate{
			NoteEntity: notes.NoteEntity{
				ID:        uid,
				Title:     title,
				Content:   content,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
			},
			Tags: tags,
		})
	}

	return entries, nil
}

func (r *NotesRepository) TagsExist(tags []notes.NoteTagEntity) (bool, error) {
	for _, t := range tags {
		found, err := r.db.CheckTagExistence(t.ID.String())
		if err != nil {
			return false, err
		}
		if !found {
			return false, nil
		}
	}
	return true, nil
}
