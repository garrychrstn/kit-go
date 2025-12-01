package routes

import (
	"github.com/garrychrstn/kit-go/db"
	"github.com/garrychrstn/kit-go/src/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupAuthRoutes(router *gin.Engine, queries *db.Queries, pool *pgxpool.Pool) {
	api := router.Group("/v1/auth")
	con := controllers.NewAuthController(queries, pool)
	{
		api.POST("/login", con.Login)
		api.POST("/whoami", con.WhoAmI)
	}
}
