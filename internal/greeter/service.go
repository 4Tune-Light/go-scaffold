package greeter

import (
	"context"
	"fmt"
	"strings"

	greeterdto "github.com/rizky/go-scaffold/internal/greeter/dto"
	"github.com/rizky/go-scaffold/pkg/database"
)

type Service interface {
	Greet(ctx context.Context, name string) (*greeterdto.GreetResponse, error)
}

type service struct {
	repo Repository
	tx   *database.Transactor
}

func NewService(repo Repository, tx *database.Transactor) Service {
	return &service{repo: repo, tx: tx}
}

func (s *service) Greet(ctx context.Context, name string) (*greeterdto.GreetResponse, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrNameRequired
	}

	greeting := fmt.Sprintf("Hello, %s!", name)

	g := &Greeting{Name: name, Message: greeting}
	if err := g.IsValid(); err != nil {
		return nil, err
	}

	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.repo.Save(txCtx, g.Name, g.Message)
	})
	if err != nil {
		return nil, err
	}

	return &greeterdto.GreetResponse{Message: greeting}, nil
}
