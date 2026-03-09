package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/alecthomas/kong"
	"github.com/lmittmann/tint"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/jghiloni/coredns-pg/common/codegen/generate"
	"github.com/jghiloni/coredns-pg/common/config"
)

type generatorArgs struct {
	Database       config.DatabaseConfig `required:"" group:"Connection Info" xor:"Connection Info" embed:"" env_prefix:"COREDNS_POSTGRES_DATABASE_" prefix:"database."`
	UseTemporaryDB bool                  `required:"" group:"Connection Info" xor:"Connection Info"                                                                     help:"Use a temporary database to generate code. Requires a container runtime."`
}

func main() {
	var cliArgs generatorArgs

	k := kong.Parse(&cliArgs)

	logger := slog.New(tint.NewHandler(k.Stdout, &tint.Options{
		TimeFormat: time.DateTime,
	}))

	slog.SetDefault(logger)

	sigmon, stopMonitoring := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stopMonitoring()

	conn, closer, err := getDB(sigmon, cliArgs)
	defer closer.Close()

	if err != nil {
		logFatal("Could not open DB", err)
	}

	generator := gen.NewGenerator(gen.Config{
		OutPath:           "generated/db",
		ModelPkgPath:      "generated/tables",
		WithUnitTest:      false,
		FieldNullable:     false,
		FieldCoverable:    false,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithGeneric,
	})

	generator.UseDB(conn.Debug())

	if err = generate.ORMCode(sigmon, generator); err != nil {
		logFatal("Could not generate ORM code", err)
	}
}

func logFatal(msg string, err error) {
	slog.Error(msg, slog.Any("err", err))
	os.Exit(1)
}

type nopCloser bool

func (nopCloser) Close() error { return nil }

var globalNopCloser nopCloser

func getDB(ctx context.Context, args generatorArgs) (*gorm.DB, io.Closer, error) {
	// never let the closer return val be nil, just pass globalNopCloser
	if args.UseTemporaryDB {
		t, err := generate.CreateTemporaryDB(ctx)
		if err == nil {
			db, e2 := t.DB()
			return db, t, e2
		}

		return nil, globalNopCloser, err
	}

	db, err := gorm.Open(postgres.Open(args.Database.ConnectionString()))

	return db, globalNopCloser, err
}
