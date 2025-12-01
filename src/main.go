package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/garrychrstn/kit-go/db"
	"github.com/garrychrstn/kit-go/src/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func runMigrations(dbURL string) error {
	m, err := migrate.New("file://db/migrations", dbURL)
	if err != nil {
		return fmt.Errorf("migrate.New failed: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("m.Up failed: %w", err)
	}

	fmt.Println("✓ Migrations completed successfully")
	return nil
}
func main() {
	ctx := context.Background()
	dbUrl := ""
	if err := godotenv.Load(".env.local"); err != nil {
		fmt.Printf("================================%s MODE==================================\n", os.Getenv("ENVI"))
		dbUrl = os.Getenv("DB_URL")
		fmt.Printf("Running migrations to %s...\n", dbUrl)
		if err := runMigrations(dbUrl); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		fmt.Println("================================DEVELOPMENT MODE==================================")
		dbUrl = os.Getenv("DB_URL")
	}
	dbPool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer dbPool.Close()

	queries := db.New(dbPool)
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"PUT", "DELETE", "GET", "POST", "PATCH"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	routes.SetupAuthRoutes(r, queries, dbPool)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if os.Getenv("ENVI") == "local" {
		r.Run("127.0.0.1:" + port)
	} else {
		r.Run(":" + port)
	}
}
