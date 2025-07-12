package postgres

import (
	"context"
	"crypto/rand"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/ffss92/webhookd/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
)

const (
	testUser     = "test-user"
	testPassword = "test-pwd"
	testDatabase = "test-db"

	postgresImage = "postgres:17-alpine"
)

type TestInstance struct {
	pool       *dockertest.Pool
	container  *dockertest.Resource
	skipReason string
	url        *url.URL

	mu sync.Mutex
	db *sql.DB
}

func MustTestInstance() *TestInstance {
	ti, err := NewTestInstance()
	if err != nil {
		log.Fatal(err)
	}
	return ti
}

func NewTestInstance() (*TestInstance, error) {
	if !flag.Parsed() {
		flag.Parse()
	}

	if testing.Short() {
		return &TestInstance{
			skipReason: "Skipping database tests (-short flag)",
		}, nil
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	repository, tag, ok := strings.Cut(postgresImage, ":")
	if !ok {
		return nil, fmt.Errorf("invalid docker image: %q", postgresImage)
	}

	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: repository,
		Tag:        tag,
		Env: []string{
			"POSTGRES_USER=" + testUser,
			"POSTGRES_PASSWORD=" + testPassword,
			"POSTGRES_DB=" + testDatabase,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start db container: %w", err)
	}

	connURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(testUser, testPassword),
		Host:   container.GetHostPort("5432/tcp"),
		Path:   testDatabase,
	}

	var db *sql.DB
	err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("pgx", connURL.String())
		if err != nil {
			return err
		}
		db.SetMaxIdleConns(1)
		db.SetMaxOpenConns(1)
		if err := db.Ping(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	if err := migrations.Up(db); err != nil {
		return nil, fmt.Errorf("failed to migrate db: %w", err)
	}

	return &TestInstance{
		pool:      pool,
		container: container,
		url:       connURL,
		db:        db,
	}, nil
}

func (i *TestInstance) Close() error {
	if i.skipReason != "" {
		return nil
	}

	err := i.db.Close()
	if err != nil {
		if purgeErr := i.pool.Purge(i.container); purgeErr != nil {
			return fmt.Errorf("failed to purge container (%v): %w", purgeErr, err)
		}
		return fmt.Errorf("failed to close db: %w", err)
	}

	if err := i.pool.Purge(i.container); err != nil {
		return fmt.Errorf("failed to purge container: %w", err)
	}
	return nil
}

func (i *TestInstance) NewPool(tb testing.TB) *pgxpool.Pool {
	tb.Helper()

	if i.skipReason != "" {
		tb.Skip(i.skipReason)
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	dbName := rand.Text()
	query := fmt.Sprintf("CREATE DATABASE %q WITH TEMPLATE %q", dbName, testDatabase)
	_, err := i.db.Exec(query)
	if err != nil {
		tb.Fatal(err)
	}

	ctx := context.Background()
	connURL := i.url.ResolveReference(&url.URL{Path: dbName})
	pool, err := pgxpool.New(ctx, connURL.String())
	if err != nil {
		tb.Fatal(err)
	}

	tb.Cleanup(func() {
		pool.Close()

		i.mu.Lock()
		defer i.mu.Unlock()

		query := fmt.Sprintf("DROP DATABASE IF EXISTS %q WITH (FORCE)", dbName)
		_, err := i.db.Exec(query)
		if err != nil {
			tb.Errorf("failed to drop database: %v", err)
		}
	})

	return pool
}
