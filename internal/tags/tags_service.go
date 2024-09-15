package tags

import "github.com/germandv/apio/internal/id"

// IService defines the interface that a Tag Service must implement.
type IService interface {
	Create(name string) (id.ID, error)
	GetAll() ([]TagAggregate, error)
}

// Service implements the Tag Service.
type Service struct {
	repo Repository
}

// NewService create a Service that satisfies the Tag Service interface.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(name string) (id.ID, error) {
	t, err := FromReq(name)
	if err != nil {
		return id.Zero(), err
	}

	err = s.repo.Save(t)
	if err != nil {
		return id.Zero(), err
	}

	return t.ID, nil
}

func (s *Service) GetAll() ([]TagAggregate, error) {
	tags, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	return tags, nil
}
