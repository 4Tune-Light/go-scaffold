package greeter

import (
	"context"
	"fmt"
	"strings"
)

type Repository interface {
	Save(ctx context.Context, name, greeting string) error
}

type Service interface {
	Greet(ctx context.Context, name string) (string, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Greet(ctx context.Context, name string) (string, error) {
	greeting := fmt.Sprintf("Hello, %s!", strings.TrimSpace(name))
	if err := s.repo.Save(ctx, name, greeting); err != nil {
		return "", err
	}
	return greeting, nil
}
