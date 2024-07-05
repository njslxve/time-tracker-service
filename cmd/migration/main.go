package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/njslxve/time-tracker-service/internal/config"
	"github.com/njslxve/time-tracker-service/migrations"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	goose.SetBaseFS(migrations.EmbedFS)

	err = goose.Up(conn, ".")
	if err != nil {
		log.Fatal(err)
	}
}
