package main

import (
	"context"
	"elkonterm/config"
	"elkonterm/internal/repository/postgres"
	"elkonterm/pkg/logger"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pressly/goose/v3"
)

const dir = "./migrations"

var flags = flag.NewFlagSet("migrate", flag.ExitOnError)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	cfg := config.GetConfig()
	logger := logger.GetLogger()

	db, err := postgres.NewPostgresDB(cfg, logger)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	command := args[0]
	switch command {
	case "up", "down", "redo", "status":
		err = goose.RunContext(context.Background(), command, db.DB, dir, args...)
	case "create":
		if len(args) < 2 {
			fmt.Println("Usage: create [migration_name]")
			return
		}
		migrationName := args[1]
		goose.SetSequential(true)
		err = goose.Create(db.DB, dir, migrationName, "sql")
		if err != nil {
			log.Fatalf("Failed to create migration: %v\n", err)
		}
		fmt.Printf("Migration '%s' created successfully.\n", migrationName)
	default:
		err = goose.RunContext(context.Background(), "help", db.DB, dir, args...)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Println("Usage: make [OPTIONS] COMMAND")
	fmt.Println("Options:")
	fmt.Println("  -h, --help		Show this help message")
	fmt.Println("Commands:")
	fmt.Println("  up			Migrate the database to the most recent version available")
	fmt.Println("  down			Roll back the version by 1")
	fmt.Println("  redo			Roll back the most recently applied migration, then run it again")
	fmt.Println("  status		Print the status of all migrations")
	fmt.Println("  create [name]		Create a new migration file with the specified name")
}
