package memorydb

import (
	"github.com/germandv/apio/internal/id"
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
		uid, err := id.Parse(row.ID)
		if err != nil {
			return nil, err
		}

		name, err := tags.ParseName(row.Name)
		if err != nil {
			return nil, err
		}

		entries = append(entries, tags.TagAggregate{
			TagEntity: tags.TagEntity{ID: uid, Name: name},
		})
	}

	return entries, nil
}
