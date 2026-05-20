package greeter

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rizky/go-scaffold/pkg/database"
)

type Repository interface {
	Save(ctx context.Context, name, greeting string) error
	FindByID(ctx context.Context, id string) (*Greeting, error)
}

type repository struct {
	q *database.Querier
}

func NewRepository(q *database.Querier) Repository {
	return &repository{q: q}
}

func (r *repository) Save(ctx context.Context, name, greeting string) error {
	query := `INSERT INTO greetings (name, message) VALUES ($1, $2)`
	_, err := r.q.Exec(ctx, query, name, greeting)
	return err
}

func (r *repository) FindByID(ctx context.Context, id string) (*Greeting, error) {
	query := `SELECT id, name, message, created_at FROM greetings WHERE id = $1`
	g := &Greeting{}
	err := r.q.QueryRow(ctx, query, id).Scan(&g.ID, &g.Name, &g.Message, &g.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return g, nil
}
