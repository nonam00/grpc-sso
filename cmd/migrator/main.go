package main

import (
	  "errors"
	  "flag"
	  "fmt"

	  "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    var connectionUrl, migrationsPath, migrationsTable string
  
    flag.StringVar(&connectionUrl, "connection-url", "", "postgres url for db")
    flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
    flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations")
    flag.Parse()
  
    if connectionUrl == "" {
        panic("connection-url is required")
    }

    if migrationsPath == "" {
        panic("migrations-path is required") 
    }

    m, err := migrate.New(
        "file://"+migrationsPath,
        fmt.Sprintf("%s/%s", connectionUrl, migrationsTable),
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

