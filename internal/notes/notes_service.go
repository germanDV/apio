package notes

import (
	"time"

	"github.com/germandv/apio/internal/errs"
	"github.com/germandv/apio/internal/id"
)

// IService defines the interface that a Note Service must implement.
type IService interface {
	Create(title, content string, tagIDs []string, createdAt time.Time) (id.ID, error)
	GetAll() ([]NoteAggregate, error)
}

// Service implements the Note Service.
type Service struct {
	repo Repository
}

// NewService create a Service that satisfies the Note Service interface.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(title, content string, tagIDs []string, createdAt time.Time) (id.ID, error) {
	noteTitle, err := ParseTitle(title)
	if err != nil {
		return id.Zero(), err
	}

	noteContent, err := ParseContent(content)
	if err != nil {
		return id.Zero(), err
	}

	uid := id.New()
	note := NoteAggregate{
		NoteEntity: NoteEntity{
			ID:        uid,
			Title:     noteTitle,
			Content:   noteContent,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		},
	}

	noteTags := make([]NoteTagEntity, 0, len(tagIDs))
	for _, t := range tagIDs {
		tagID, err := id.Parse(t)
		if err != nil {
			return id.Zero(), err
		}
		noteTags = append(noteTags, NoteTagEntity{ID: tagID, Name: ""})
	}

	allTagsFound, err := s.repo.TagsExist(noteTags)
	if err != nil {
		return id.Zero(), err
	}
	if !allTagsFound {
		return id.Zero(), errs.ErrTagNotFound
	}
	note.Tags = noteTags

	err = s.repo.Save(note)
	if err != nil {
		return id.Zero(), err
	}

	return uid, nil
}

func (s *Service) GetAll() ([]NoteAggregate, error) {
	notes, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	return notes, nil
}
