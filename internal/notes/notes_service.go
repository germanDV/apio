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
	n, err := FromReq(title, content, tagIDs, createdAt)
	if err != nil {
		return id.Zero(), err
	}

	// All Tags must be created in advance.
	allTagsFound, err := s.repo.TagsExist(n.Tags)
	if err != nil {
		return id.Zero(), err
	}
	if !allTagsFound {
		return id.Zero(), errs.ErrTagNotFound
	}

	err = s.repo.Save(n)
	if err != nil {
		return id.Zero(), err
	}

	return n.ID, nil
}

func (s *Service) GetAll() ([]NoteAggregate, error) {
	notes, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	return notes, nil
}
