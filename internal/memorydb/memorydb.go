package memorydb

import (
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/germandv/apio/internal/errs"
)

type TagRow struct {
	ID   string
	Name string
}

type NoteRow struct {
	ID        string
	Title     string
	Content   string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NoteTagRow struct {
	NoteID string
	TagID  string
}

type data struct {
	tags     map[string]TagRow
	notes    map[string]NoteRow
	noteTags []NoteTagRow
}

type DB struct {
	mu   sync.RWMutex
	data data
}

var database = &DB{
	mu: sync.RWMutex{},
	data: data{
		tags:     make(map[string]TagRow),
		notes:    make(map[string]NoteRow),
		noteTags: make([]NoteTagRow, 0),
	},
}

func (db *DB) SaveTag(id string, name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, t := range db.data.tags {
		if t.Name == name {
			return errs.ErrDuplicateTag
		}
	}

	db.data.tags[id] = TagRow{id, name}
	return nil
}

func (db *DB) GetTags() ([]TagRow, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return slices.Collect(maps.Values(db.data.tags)), nil
}

func (db *DB) SaveNote(id, title, content, createdBy string, createdAt, updatedAt time.Time) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data.notes[id] = NoteRow{
		ID:        id,
		Title:     title,
		Content:   content,
		CreatedBy: createdBy,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return nil
}

func (db *DB) SaveNoteTags(noteID string, tagIDs []string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	for _, tagID := range tagIDs {
		row := NoteTagRow{NoteID: noteID, TagID: tagID}
		db.data.noteTags = append(db.data.noteTags, row)
	}
	return nil
}

func (db *DB) GetNotes() ([]NoteRow, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	items := slices.Collect(maps.Values(db.data.notes))
	return items, nil
}

func (db *DB) GetNoteTags(id string) ([]TagRow, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	rows := slices.DeleteFunc(slices.Clone(db.data.noteTags), func(r NoteTagRow) bool {
		return r.NoteID != id
	})

	ts := make([]TagRow, 0, len(rows))
	for _, row := range rows {
		tag, ok := db.data.tags[row.TagID]
		if ok {
			ts = append(ts, tag)
		}
	}

	return ts, nil
}

func (db *DB) CheckTagExistence(tagID string) (bool, error) {
	_, ok := db.data.tags[tagID]
	return ok, nil
}

func (db *DB) CountNotesByTag(tagID string) (int, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	count := 0
	for _, r := range db.data.noteTags {
		if r.TagID == tagID {
			count++
		}
	}

	return count, nil
}
