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
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	// Load env file if it exists (local dev)
	isLocal := godotenv.Load(".env.local") == nil

	if isLocal {
		fmt.Println("================================ LOCAL DEVELOPMENT MODE ==================================")
	} else {
		fmt.Printf("================================%s MODE==================================\n", os.Getenv("ENVI"))
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	dbPool, err := pgxpool.New(ctx, dbURL)
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
