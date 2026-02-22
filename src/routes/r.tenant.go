package routes

import (
	"github.com/dirental/core/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TenantRoutes(router *gin.Engine, queries *db.Queries, pool *pgxpool.Pool) {
	// api := router.Group("/v1/tenant")
	// con := controllers.AuthController(queries, pool)

}
