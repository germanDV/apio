package memorydb

import (
	"github.com/germandv/apio/internal/tags"
)

// TagsRepository implements the Repository interface for Tags.
type TagsRepository struct {
	db *DB
}

// NewTagsRepository creates a TagsRepository that satisfies the Repository interface for Tags.
func NewTagsRepository() tags.Repository {
	return &TagsRepository{
		db: database,
	}
}

func (r *TagsRepository) Save(tag tags.TagAggregate) error {
	return r.db.SaveTag(tag.ID.String(), tag.Name.String())
}

func (r *TagsRepository) List() ([]tags.TagAggregate, error) {
	rows, err := r.db.GetTags()
	if err != nil {
		return nil, err
	}

	entries := make([]tags.TagAggregate, 0, len(rows))
	for _, row := range rows {
		noteCount, err := r.db.CountNotesByTag(row.ID)
		if err != nil {
			return nil, err
		}

		t, err := tags.FromDB(row.ID, row.Name, noteCount)
		if err != nil {
			return nil, err
		}

		entries = append(entries, t)
	}

	return entries, nil
}
