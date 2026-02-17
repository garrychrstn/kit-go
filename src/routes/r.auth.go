package routes

import (
	"github.com/dirental/core/db"
	"github.com/dirental/core/src/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupAuthRoutes(router *gin.Engine, queries *db.Queries, pool *pgxpool.Pool) {
	api := router.Group("/v1/auth")
	con := controllers.AuthController(queries, pool)
	{
		api.POST("/login", con.Login)
	}
}
