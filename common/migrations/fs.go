package migrations

import (
	"context"
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed "0*.sql"
var all embed.FS

func RunMigrations(ctx context.Context, dbURL string) error {
	goose.SetDialect(string(goose.DialectPostgres))

	db, err := sql.Open("pgx/v5", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	goose.SetBaseFS(all)

	return goose.UpContext(ctx, db, ".")
}
