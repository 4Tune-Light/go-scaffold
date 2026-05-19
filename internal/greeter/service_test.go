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

func TestGreet_Success(t *testing.T) {
	svc := NewService(&mockRepository{})
	result, err := svc.Greet(context.Background(), "John")
	assert.NoError(t, err)
	assert.Equal(t, "Hello, John!", result)
}

func TestGreet_TrimsSpaces(t *testing.T) {
	svc := NewService(&mockRepository{})
	result, err := svc.Greet(context.Background(), "  Jane  ")
	assert.NoError(t, err)
	assert.Equal(t, "Hello, Jane!", result)
}

func TestGreet_RepoError(t *testing.T) {
	svc := NewService(&mockRepository{saveErr: assert.AnError})
	_, err := svc.Greet(context.Background(), "John")
	assert.Error(t, err)
}
