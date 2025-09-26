package database

import (
	"aicvevaluator/internal/config"

	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var oncePostgres sync.Once
var pgSqlx *sqlx.DB

func InitPostgresql(ctx context.Context, conf *config.Config) {
	oncePostgres.Do(func() {

		val := url.Values{}
		val.Add("sslmode", "disable")
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", conf.DB.User, conf.DB.Password, conf.DB.Host, conf.DB.Port, conf.DB.Name, val.Encode())
		autoMigrate(&log.Logger, dsn)

		// Create *sql.DB using pgx driver
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to postgresql database")
		}

		// Wrap into sqlx.DB once
		sqlxDB := sqlx.NewDb(db, "pgx")

		// Configure pooling on the underlying *sql.DB
		sqlxDB.SetConnMaxLifetime(conf.DB.ConnectionLifetime)
		sqlxDB.SetMaxOpenConns(conf.DB.MaxOpen)
		sqlxDB.SetMaxIdleConns(conf.DB.MaxIdle)
		// SetConnMaxIdleTime is available on Go 1.15+. If you need it:
		if conf.DB.ConnectionIdle > 0 {
			sqlxDB.DB.SetConnMaxIdleTime(conf.DB.ConnectionIdle)
		}

		// Ping using sqlx
		if err = sqlxDB.PingContext(ctx); err != nil {
			log.Fatal().Err(err).Msg("Failed to ping postgresql database")
		}

		log.Info().Msg("Postgresql database connected")

		pgSqlx = sqlxDB
	})
}

func autoMigrate(log *zerolog.Logger, dsn string) {
	baseDir := "database/migrations"
	files, err := os.ReadDir(baseDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Warn().Msg("Migration directory does not exist, skipping migration")
			return
		}

		log.Fatal().Err(err).Msg("Failed to read migration directory")
	}

	if len(files) == 0 {
		log.Info().Msg("No migration files found, skipping migration")
		return
	}

	m, err := migrate.New("file://"+baseDir, dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create migration")
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("Failed to migrate")
	}
}

func GetPostgresql() *sqlx.DB {
	return pgSqlx
}
