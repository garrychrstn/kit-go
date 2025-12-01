//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println("No .env file found, using system environment")
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Error: migration name is required")
			fmt.Println("Usage: go run mgr.go create <migration_name>")
			os.Exit(1)
		}
		createMigration(os.Args[2])

	case "do":
		if len(os.Args) < 3 {
			fmt.Println("Error: migration action is required (up/down)")
			fmt.Println("Usage: go run mgr.go do <up|down>")
			os.Exit(1)
		}
		runMigration(os.Args[2])

	case "force":
		if len(os.Args) < 3 {
			fmt.Println("Error: version number is required")
			fmt.Println("Usage: go run mgr.go force <version>")
			os.Exit(1)
		}
		runReset(os.Args[2])

	case "version":
		showVersion()

	default:
		fmt.Printf("Error: unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run mgr.go create <migration_name>  - Create a new migration file")
	fmt.Println("  go run mgr.go do <up|down>             - Run migrations up or down")
	fmt.Println("  go run mgr.go force <version>          - Force database to specific version")
	fmt.Println("  go run mgr.go version                  - Show current migration version")
}

func createMigration(migrationName string) {
	fmt.Printf("Creating migration: %s\n", migrationName)

	// Check if migrations directory exists
	if _, err := os.Stat("./db/migrations"); os.IsNotExist(err) {
		if err := os.MkdirAll("./db/migrations", 0755); err != nil {
			fmt.Printf("Error creating migrations directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Get list of existing migrations to determine next sequence number
	files, err := os.ReadDir("./db/migrations")
	if err != nil {
		fmt.Printf("Error reading migrations directory: %v\n", err)
		os.Exit(1)
	}

	nextSeq := len(files)/2 + 1 // Each migration has 2 files (up/down)

	upFile := fmt.Sprintf("./db/migrations/%06d_%s.up.sql", nextSeq, migrationName)
	downFile := fmt.Sprintf("./db/migrations/%06d_%s.down.sql", nextSeq, migrationName)

	if err := os.WriteFile(upFile, []byte(""), 0644); err != nil {
		fmt.Printf("Error creating up migration: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(downFile, []byte(""), 0644); err != nil {
		fmt.Printf("Error creating down migration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Migration '%s' created successfully!\n", migrationName)
	fmt.Printf("  - %s\n", upFile)
	fmt.Printf("  - %s\n", downFile)
}

func runMigration(action string) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		fmt.Println("Error: DB_URL environment variable not set")
		os.Exit(1)
	}

	fmt.Println("Action:", action)
	fmt.Println("DB:", dbURL)
	fmt.Println("Running migrations...")

	m, err := migrate.New("file://db/migrations", dbURL)
	if err != nil {
		fmt.Printf("Error creating migrate instance: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	switch action {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		fmt.Printf("Error: unknown action '%s'. Use 'up' or 'down'\n", action)
		os.Exit(1)
	}

	if err != nil && err != migrate.ErrNoChange {
		fmt.Printf("Error running migrations: %v\n", err)
		os.Exit(1)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("No changes - migrations already up to date")
	} else {
		fmt.Println("Migrations completed successfully!")
	}
}

func runReset(toVersion string) {
	version, err := strconv.Atoi(toVersion)
	if err != nil {
		fmt.Printf("Error: invalid version number '%s'\n", toVersion)
		os.Exit(1)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		fmt.Println("Error: DB_URL environment variable not set")
		os.Exit(1)
	}

	fmt.Printf("Forcing migration to version %d\n", version)

	m, err := migrate.New("file://db/migrations", dbURL)
	if err != nil {
		fmt.Printf("Error creating migrate instance: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	if err := m.Force(version); err != nil {
		fmt.Printf("Error forcing migration: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Reset completed successfully!")
}

func showVersion() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		fmt.Println("Error: DB_URL environment variable not set")
		os.Exit(1)
	}

	m, err := migrate.New("file://db/migrations", dbURL)
	if err != nil {
		fmt.Printf("Error creating migrate instance: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil {
		fmt.Printf("Error getting version: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Current migration version: %d\n", version)
	if dirty {
		fmt.Println("Warning: Database is in dirty state")
	}
}
