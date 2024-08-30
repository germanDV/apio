package memorydb

import (
	"fmt"
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
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NoteTagRow struct {
	NoteID string
	TagID  string
}

type data struct {
	tags      map[string]TagRow
	notes     map[string]NoteRow
	note_tags []NoteTagRow
}

type DB struct {
	mu   sync.RWMutex
	data data
}

var database = &DB{
	mu: sync.RWMutex{},
	data: data{
		tags:      make(map[string]TagRow),
		notes:     make(map[string]NoteRow),
		note_tags: make([]NoteTagRow, 0),
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
	fmt.Printf("notes: %v\n", db.data.notes)
	fmt.Printf("tags: %v\n", db.data.tags)
	return slices.Collect(maps.Values(db.data.tags)), nil
}

func (db *DB) SaveNote(id, title, content string, createdAt, updatedAt time.Time) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data.notes[id] = NoteRow{
		ID:        id,
		Title:     title,
		Content:   content,
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
		db.data.note_tags = append(db.data.note_tags, row)
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

	fmt.Printf("notes: %v\n", db.data.notes)
	fmt.Printf("tags: %v\n", db.data.tags)

	rows := slices.DeleteFunc(slices.Clone(db.data.note_tags), func(r NoteTagRow) bool {
		return r.NoteID != id
	})

	ts := make([]TagRow, 0, len(rows))
	for _, row := range rows {
		tag, ok := db.data.tags[row.TagID]
		if !ok {
			fmt.Println("tag not found")
		} else {
			ts = append(ts, tag)
		}
	}

	return ts, nil
}

func (db *DB) CheckTagExistence(tagID string) (bool, error) {
	_, ok := db.data.tags[tagID]
	return ok, nil
}
