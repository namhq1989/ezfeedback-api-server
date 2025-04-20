package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	apperrors "github.com/namhq1989/ezfeedback-api-server/internal/error"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

type Operations interface {
	GetPgDb() *sql.DB
	IsNoRowsError(err error) bool
}

type Database struct {
	pg *sql.DB
}

func NewDatabaseClient(pgConn string) *Database {
	return &Database{
		pg: initPostgres(pgConn),
	}
}

func initPostgres(conn string) *sql.DB {
	db, err := otelsql.Open("pgx", conn,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithDBName("ezfeedback"),
	)
	if err != nil {
		panic(err)
	}

	// config
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// verify the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}

	fmt.Printf("⚡️ [postgresql]: connected \n")

	return db
}

func (d Database) GetPgDb() *sql.DB {
	return d.pg
}

func (Database) IsNoRowsError(err error) bool {
	return errors.Is(err, qrm.ErrNoRows)
}

func (Database) IsDuplicatedError(err error) (bool, error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return true, apperrors.GetErrorUniqueViolationByConstraintName(pgErr.ConstraintName)
		}
	}
	return false, nil
}
