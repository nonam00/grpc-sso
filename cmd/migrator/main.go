package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    var connectionStr, migrationsPath string
  
    flag.StringVar(&connectionStr, "connection-string", "", "postgres connection string for db")
    flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
    flag.Parse()
  
    if connectionStr == "" {
        panic("connection-string is required")
    }

    if migrationsPath == "" {
        panic("migrations-path is required") 
    }

    db, err := sql.Open("postgres", connectionStr)

    if err != nil {
        panic(err)
    }

    defer db.Close()

    driver, err := postgres.WithInstance(db, &postgres.Config{})

    m, err := migrate.NewWithDatabaseInstance(
        "file://"+migrationsPath,
        "postgres",
        driver,
    )

    if err != nil {
        panic(err)
    }

    if err := m.Up(); err != nil {
        if errors.Is(err, migrate.ErrNoChange) {
            fmt.Println("no migrations to apply")
            return
        }
        panic(err)
    }

    fmt.Println("migrations applied successfully")
}

