package db

import (
	"UserService/internal/config"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"time"
)

var dsn = "postgres://%s:%s@%s:%s/%s?sslmode=%s"

func ConnectPsql(cfg *config.MainDatabase) *sqlx.DB {
	db, err := sqlx.Connect("postgres", GetDSN(cfg))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Minute * time.Duration(cfg.MaxLifetimeConnectionsMin))

	return db
}

func GetDSN(cfg *config.MainDatabase) string {
	psqlDsn := fmt.Sprintf(dsn, cfg.Username, cfg.Password,
		cfg.Host, cfg.Port,
		cfg.Database, cfg.SslMode)

	return psqlDsn
}

type TxManager struct {
	DB                    *sqlx.DB
	MaxTransactionRetries int
}

func NewTxManager(db *sqlx.DB, maxTransactionRetries int) *TxManager {
	return &TxManager{
		DB:                    db,
		MaxTransactionRetries: maxTransactionRetries,
	}
}

func (tm *TxManager) RunSerializable(
	ctx context.Context,
	fn func(tx *sqlx.Tx) error,
) error {
	return tm.transaction(
		ctx,
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
		},
		tm.DB,
		fn,
	)
}

func (tm *TxManager) RunReadCommited(
	ctx context.Context,
	fn func(tx *sqlx.Tx) error,
) error {
	return tm.transaction(
		ctx,
		&sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
		},
		tm.DB,
		fn,
	)
}

func (tm *TxManager) transaction(
	ctx context.Context,
	opts *sql.TxOptions,
	db *sqlx.DB,
	fn func(tx *sqlx.Tx) error,
) error {
	maxAttempts := tm.MaxTransactionRetries
	if maxAttempts < 0 {
		maxAttempts = 3
	}

	var err error

	for i := 0; i < maxAttempts; i++ {
		err = func() error {
			tx, err := db.BeginTxx(ctx, opts)
			if err != nil {
				return err
			}

			defer func() {
				_ = tx.Rollback()
			}()

			if err := fn(tx); err != nil {
				return err
			}

			return tx.Commit()
		}()

		if err == nil {
			return nil
		}

		if !tm.isErrorRetryable(err) {
			return err
		}
	}

	return err
}

func (tm *TxManager) isErrorRetryable(err error) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code.Name() {
		case "serialization_failure",
			"deadlock_detected":
			return true
		}
	}
	return false
}
