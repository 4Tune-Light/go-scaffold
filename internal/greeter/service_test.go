package greeter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	saveErr error
}

func (m *mockRepository) Save(ctx context.Context, name, greeting string) error {
	return m.saveErr
}

func (m *mockRepository) FindByID(ctx context.Context, id string) (*Greeting, error) {
	return nil, nil
}

func TestGreet_Success(t *testing.T) {
	svc := NewService(&mockRepository{}, nil)
	result, err := svc.Greet(context.Background(), "John")
	assert.NoError(t, err)
	assert.Equal(t, "Hello, John!", result.Message)
}

func TestGreet_TrimsSpaces(t *testing.T) {
	svc := NewService(&mockRepository{}, nil)
	result, err := svc.Greet(context.Background(), "  Jane  ")
	assert.NoError(t, err)
	assert.Equal(t, "Hello, Jane!", result.Message)
}

func TestGreet_EmptyName(t *testing.T) {
	svc := NewService(&mockRepository{}, nil)
	_, err := svc.Greet(context.Background(), "  ")
	assert.ErrorIs(t, err, ErrNameRequired)
}

func TestGreet_RepoError(t *testing.T) {
	svc := NewService(&mockRepository{saveErr: assert.AnError}, nil)
	_, err := svc.Greet(context.Background(), "John")
	assert.Error(t, err)
}
