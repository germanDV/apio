package memorydb

import (
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
		note.CreatedBy.String(),
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

		mappedTags := tools.Map(tagRows, func(t TagRow) struct {
			ID   string
			Name string
		} {
			return struct {
				ID   string
				Name string
			}{
				ID:   t.ID,
				Name: t.Name,
			}
		})

		n, err := notes.FromDB(row.ID, row.Title, row.Content, row.CreatedBy, row.CreatedAt, row.UpdatedAt, mappedTags)
		if err != nil {
			return nil, err
		}
		entries = append(entries, n)
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
