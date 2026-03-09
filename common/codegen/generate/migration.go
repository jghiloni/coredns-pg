package generate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"

	gormpg "gorm.io/driver/postgres"

	"github.com/jghiloni/coredns-pg/common/migrations"
)

type TemporaryDB struct {
	tempDBContainer *postgres.PostgresContainer
	dbURL           *string
}

func CreateTemporaryDB(ctx context.Context) (*TemporaryDB, error) {
	var err error

	t := new(TemporaryDB)

	slog.Info("Creating Temporary DB")

	if t.tempDBContainer, err = postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithUsername("postgres"),
		postgres.WithPassword(uuid.NewString()),
		postgres.WithDatabase("coredns"),
		postgres.BasicWaitStrategies(),
		testcontainers.WithLogger(slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo)),
	); err != nil {
		return nil, fmt.Errorf("could not create test db: %w", err)
	}

	url, err := t.tempDBContainer.ConnectionString(ctx)
	if err != nil {
		t.Close()
		return nil, fmt.Errorf("could not get db connection string: %w", err)
	}

	t.dbURL = new(url)

	err = migrations.RunMigrations(ctx, *t.dbURL)
	if err != nil {
		errFmt := "could not run migrations on new database: %w"
		args := []any{err}

		e2 := t.Close()
		if e2 != nil {
			errFmt = errFmt + ". Additionally, an error occurred closing the DB: %w"

			args = append(args, e2)
		}

		err = fmt.Errorf(errFmt, args...)
	}

	return t, err
}

var ErrDBNotCreated = errors.New("temporary db has not been created yet")

func (t *TemporaryDB) DB() (*gorm.DB, error) {
	if t.dbURL == nil {
		return nil, ErrDBNotCreated
	}

	return gorm.Open(gormpg.Open(*t.dbURL))
}

func (t *TemporaryDB) Close() error {
	if t.tempDBContainer == nil {
		return nil
	}

	err := testcontainers.TerminateContainer(t.tempDBContainer)
	t.tempDBContainer = nil

	return err
}
