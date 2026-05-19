package health

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Checker struct {
	pg  *pgxpool.Pool
	rdb *redis.Client
}

func New(pg *pgxpool.Pool, rdb *redis.Client) *Checker {
	return &Checker{pg: pg, rdb: rdb}
}

type Status struct {
	Status     string            `json:"status"`
	Components map[string]string `json:"components"`
}

func (c *Checker) Check(ctx context.Context) Status {
	comps := make(map[string]string)

	if c.pg != nil {
		if err := c.pg.Ping(ctx); err != nil {
			comps["postgres"] = "down"
		} else {
			comps["postgres"] = "up"
		}
	}

	if c.rdb != nil {
		if err := c.rdb.Ping(ctx).Err(); err != nil {
			comps["redis"] = "down"
		} else {
			comps["redis"] = "up"
		}
	}

	overall := "ok"
	for _, v := range comps {
		if v == "down" {
			overall = "degraded"
			break
		}
	}

	if len(comps) == 0 {
		overall = "ok"
	}

	return Status{Status: overall, Components: comps}
}
