package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ginerator/base/config"
	"github.com/ginerator/base/errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bunotel"
)

const (
	TxContextKey string = "db-tx-key"
)

type BunPostgresDatabaseClient struct {
	DB            *bun.DB
	config        *config.DbConfig
	MigrationsDir string
}

func (client *BunPostgresDatabaseClient) getPostgresURL() string {
	return "postgres://" + client.config.Username + ":" + url.QueryEscape(client.config.Password) + "@" + client.config.Host + ":" + client.config.Port + "/" + client.config.Name + "?sslmode=disable"
}

func NewBunPostgresDatabaseClient(config *config.DbConfig, migrationsDir string) *BunPostgresDatabaseClient {
	_, err := os.Stat(migrationsDir)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("[POSTGRES CLIENT] - New - Migration folder: %s doesn't exist.", migrationsDir))
	}
	client := &BunPostgresDatabaseClient{
		MigrationsDir: migrationsDir,
	}
	client.config = config
	client.Connect()
	log.Info().Msg("Database client initialized.")
	return client
}

func (client *BunPostgresDatabaseClient) Connect() error {
	maxConnections, _ := strconv.Atoi(client.config.MaxOpenConns)
	maxIdleConnections, _ := strconv.Atoi(client.config.MaxIdleConns)
	connectionString := client.getPostgresURL()
	log.Info().Msgf("Connecting to database: %s", connectionString)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionString)))
	sqldb.SetMaxOpenConns(maxConnections)
	sqldb.SetMaxIdleConns(maxIdleConnections)
	client.DB = bun.NewDB(sqldb, pgdialect.New())
	client.DB.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName(client.config.Name)))
	err := client.DB.Ping()
	if err != nil {
		log.Error().Err(err).Msg("[POSTGRES CLIENT] - Connect - Error connecting")
	}
	return err
}

func (client *BunPostgresDatabaseClient) MigrateUp() {
	m, err := migrate.New(fmt.Sprintf("file://%s", client.MigrationsDir), client.getPostgresURL())
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Panic().Err(err).Msg("[POSTGRES CLIENT] - Connect - Error running migrations")
	}
}

func (client *BunPostgresDatabaseClient) MigrateDown() {
	m, err := migrate.New(fmt.Sprintf("file://%s", client.MigrationsDir), client.getPostgresURL())
	err = m.Down()
	log.Info().Msg("MigrateDown: Applying migration")
	if err != nil && err != migrate.ErrNoChange {
		log.Panic().Err(err).Msg("[POSTGRES CLIENT] - Connect - Error running migrations down")
	}
}

func (client *BunPostgresDatabaseClient) IsConnected() (bool, error) {
	err := client.DB.Ping()
	if err != nil {
		log.Error().Err(err).Msg("[POSTGRES CLIENT] - IsConnected - Checking connection open")
		return false, err
	}
	return true, nil
}

func (client *BunPostgresDatabaseClient) Close() {
	client.DB.Close()
}

func (repo *BunPostgresDatabaseClient) BeginTransaction(ctx *gin.Context) (context.Context, error) {
	tx, err := repo.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Error().Err(err).Msg("[POSTGRES CLIENT] - BeginTransaction - Could not begin transaction")
		return ctx, err
	}

	return context.WithValue(ctx.Request.Context(), TxContextKey, tx), nil
}

func (repo *BunPostgresDatabaseClient) ResolveTransaction(ctx *gin.Context, err error) error {
	tx := repo.getTx(ctx)
	if tx == nil {
		log.Error().Err(err).Msg("[POSTGRES CLIENT] - ResolveTransaction - Transaction is null")
		return errors.NewUnkownDatabaseError(fmt.Errorf("null transaction"))
	}

	if err == nil {
		return tx.Commit()
	}

	log.Error().Err(err).Msg("[POSTGRES CLIENT] - ResolveTransaction - Transaction rolledback")
	return tx.Rollback()
}

func (repo *BunPostgresDatabaseClient) getTx(ctx *gin.Context) *bun.Tx {
	value := ctx.Value(TxContextKey)
	if value == nil {
		return nil
	}

	tx := value.(bun.Tx)
	return &tx
}

func (repo *BunPostgresDatabaseClient) getDB(ctx *gin.Context) bun.IDB {
	tx := repo.getTx(ctx)
	if tx == nil {
		return repo.DB
	}

	return *tx
}
