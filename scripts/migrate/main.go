package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rizky/go-scaffold/internal/config"
	"github.com/rizky/go-scaffold/pkg/database"
)

func main() {
	cfg := config.MustLoad("")

	ctx := context.Background()

	pool, err := database.NewPostgresPool(ctx, cfg.PostgresConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect to postgres: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	q := database.NewQuerier(pool)

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: go run ./scripts/migrate <up|down>")
		os.Exit(1)
	}
	direction := args[0]

	migrationsDir := "migrations"
	if len(args) > 1 {
		migrationsDir = args[1]
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read migrations dir: %v\n", err)
		os.Exit(1)
	}

	var targets []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), "."+direction+".sql") {
			targets = append(targets, f.Name())
		}
	}
	sort.Strings(targets)

	for _, name := range targets {
		path := filepath.Join(migrationsDir, name)
		sql, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", name, err)
			os.Exit(1)
		}

		if _, err := q.Exec(ctx, string(sql)); err != nil {
			fmt.Fprintf(os.Stderr, "execute %s: %v\n", name, err)
			os.Exit(1)
		}
		fmt.Printf("applied: %s\n", name)
	}

	fmt.Println("migrations complete")
}
