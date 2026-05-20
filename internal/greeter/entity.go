package greeter

import "time"

type Greeting struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func (g *Greeting) IsValid() error {
	if g.Name == "" {
		return ErrNameRequired
	}
	return nil
}
