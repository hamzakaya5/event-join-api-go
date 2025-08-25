package database

import (
	"context"
	"fmt"
	"gamegos_case/models"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps the pgx connection pool
type DB struct {
	Pool *pgxpool.Pool
}

var DBConn *DB

func ConnectToPostgre() error {

	cfg := models.PostgreConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),     //"myuser",
		Password: os.Getenv("DB_PASSWORD"), //"mypassword",
		Database: os.Getenv("DB_NAME"),     //"mydb",
		MaxConns: 10,
		Port:     5432,
	}

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		fmt.Print(err, "err")
		return err
	}

	// Set maximum connections
	poolConfig.MaxConns = int32(cfg.MaxConns)
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {

		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("cannot ping database: %w", err)
	}

	DBConn = &DB{Pool: pool}
	fmt.Println("Connected to PostgreSQL database!")
	return nil
}
