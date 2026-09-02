package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Registers the "postgres" driver used by sql.Open below. This package is
	// the one that names the driver, so it is the one that must guarantee it
	// exists rather than relying on some other package importing it first.
	_ "github.com/lib/pq"
)

const (
	// connectTimeout bounds establishing a connection. It is passed to lib/pq as
	// connect_timeout because that is what puts a deadline on the socket for the
	// whole startup handshake; a context alone is not enough, since pq only
	// applies the context to the TCP dial and then reads the server's greeting
	// with no deadline at all.
	//
	// Without this the process can hang forever with no output: Docker's port
	// proxy accepts the TCP connection before Postgres is listening, so the dial
	// succeeds and the handshake read blocks indefinitely.
	connectTimeout = 5 * time.Second

	// pingTimeout is the outer backstop. It is deliberately looser than
	// connectTimeout so that pq's own network error surfaces rather than a bare
	// "context deadline exceeded", which says nothing about what went wrong.
	pingTimeout = connectTimeout + 2*time.Second
)

// DB owns the connection pool and exposes the generated queries. The pool is a
// named field rather than an embedded one so that Exec, Query, Close and the
// rest of *sql.DB are not promoted onto every consumer of this type — callers
// get the queries and nothing else.
type DB struct {
	pool *sql.DB
	*Queries
}

func NewDB(user, password, host, database, tz string, port int) (*DB, error) {
	source := fmt.Sprintf(
		"user=%s password=%s host=%s port=%v dbname=%s sslmode=disable TimeZone=%s connect_timeout=%d",
		user, password, host, port, database, tz, int(connectTimeout.Seconds()),
	)

	pool, err := sql.Open("postgres", source)
	if err != nil {
		return nil, fmt.Errorf("opening postgres connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connecting to postgres at %s:%d: %w", host, port, err)
	}

	return &DB{
		pool:    pool,
		Queries: New(pool),
	}, nil
}

func (d *DB) Close() error {
	return d.pool.Close()
}

func (d *DB) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := d.pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	if err := fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx: err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
