package main

import (
	"database/sql"
	"fmt"
	"os"
	"seeds/auth"
	"seeds/dataset"
	"seeds/pages"
	"seeds/widgets"
	"strconv"
	"strings"
)

type Env struct {
	AuthAdminBaseUrl     string
	PGDatabaseUrl        string
	PGMaxOpenConnections string
	PGMaxIdleConnections string
	AdminSecrets         string
	ApiBaseUrl           string
	Mode                 string
}

func getEnv() Env {
	return Env{
		AuthAdminBaseUrl:     os.Getenv("AUTH_ADMIN_HOST"),
		PGDatabaseUrl:        os.Getenv("PG_DATABASE_URL_APP"),
		PGMaxOpenConnections: os.Getenv("PG_MAX_OPEN_CONNECTIONS"),
		PGMaxIdleConnections: os.Getenv("PG_MAX_IDLE_CONNECTIONS"),
		AdminSecrets:         os.Getenv("CONTROL_PLANE_ADMIN_SECRETS"),
		ApiBaseUrl:           os.Getenv("API_BASE_URL"),
		Mode:                 os.Getenv("MODE"),
	}
}

func main() {
	env := getEnv()

	// TODO: use mode when it is updated for everyone
	if !strings.Contains(env.PGDatabaseUrl, "postgres:5432") {
		fmt.Println("Skipping seeds for non-local environment")
		return
	}

	maxOpen, err := strconv.Atoi(env.PGMaxOpenConnections)
	if err != nil {
		fmt.Printf("Error parsing PG_MAX_OPEN_CONNECTIONS: %v\n", err)
		return
	}
	maxIdle, err := strconv.Atoi(env.PGMaxIdleConnections)
	if err != nil {
		fmt.Printf("Error parsing PG_MAX_IDLE_CONNECTIONS: %v\n", err)
		return
	}

	pgClient, err := newPostgresClient(env.PGDatabaseUrl, maxIdle, maxOpen)
	if err != nil {
		fmt.Println("error initializing db connection")
		return
	}
	defer pgClient.Close()

	adminSecret := strings.Split(env.AdminSecrets, ",")[0]
	auth.RunAuthSeeds(env.AuthAdminBaseUrl, env.ApiBaseUrl, adminSecret, pgClient)

	fmt.Println("Running dataset seeds...")

	dataset.RunDatasetSeeds(pgClient)

	pages.RunDatasetSeeds(pgClient)

	widgets.RunWidgetsSeeds(pgClient)
}

func newPostgresClient(dbURL string, maxIdleConns, maxOpenConns int) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}
