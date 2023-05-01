package store

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib" //nolint:blank-imports
	"go.uber.org/zap"
)

//go:generate options-gen -out-filename=client_psql_options.gen.go -from-struct=PSQLOptions
type PSQLOptions struct {
	address  string `option:"mandatory" validate:"required,hostname_port"`
	username string `option:"mandatory" validate:"required"`
	password string `option:"mandatory" validate:"required"`
	database string `option:"mandatory" validate:"required"`
	debug    bool
}

func NewPSQLClient(opts PSQLOptions) (*Client, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	lg := zap.L().Named("psql")

	db, err := NewPgxDB(NewPgxOptions(opts.address, opts.username, opts.password, opts.database))
	if err != nil {
		return nil, fmt.Errorf("init db driver: %v", err)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)

	clientOpts := []Option{
		Driver(drv),
	}

	if opts.debug {
		clientOpts = append(clientOpts, Debug())
	}

	client := NewClient(clientOpts...)
	client.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			start := time.Now()
			defer func() {
				lg.Debug("query", zap.Duration("time", time.Since(start)))
			}()
			return next.Mutate(ctx, m)
		})
	})

	return client, nil
}

//go:generate options-gen -out-filename=client_psql_pgx_options.gen.go -from-struct=PgxOptions
type PgxOptions struct {
	address  string `option:"mandatory" validate:"required,hostname_port"`
	username string `option:"mandatory" validate:"required"`
	password string `option:"mandatory" validate:"required"`
	database string `option:"mandatory" validate:"required"`
}

func NewPgxDB(opts PgxOptions) (*sql.DB, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate pgx options: %v", err)
	}

	ds := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(opts.username, opts.password),
		Host:   opts.address,
		Path:   opts.database,
	}

	db, err := sql.Open("pgx", ds.String())
	if err != nil {
		return nil, fmt.Errorf("open connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping to db: %v", err)
	}

	return db, nil
}
