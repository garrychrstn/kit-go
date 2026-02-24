package controllers

import (
	"time"

	"github.com/dirental/core/db"
	"github.com/dirental/core/src/helpers"
	helperdb "github.com/dirental/core/src/helpers/db"
	"github.com/dirental/core/src/types/requests"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IStoreController struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func StoreController(queries *db.Queries, pool *pgxpool.Pool) *IStoreController {
	return &IStoreController{
		queries: queries,
		pool:    pool,
	}
}
func (c *IStoreController) SetupStore(ctx *gin.Context) {
	data, err := helpers.ValidateRequest[requests.IRequestCreateStore](ctx)

	if err != nil {
		return
	}

	newTenant := db.Store{
		Name:           data.Name,
		Phone:          data.Phone,
		Coordinate:     helperdb.SafeString(&data.Coordinate),
		Address:        "Feature Needed",
		IsActive:       false,
		Contacts:       data.Contacts,
		CreatedAt:      pgtype.Timestamptz{Time: time.Now()},
		Description:    helperdb.SafeString(&data.Description),
		Category:       data.Category,
		TermAndService: data.TermAndService,
	}
	ctx.JSON(200, gin.H{"message": "Store created successfully", "tenant": newTenant})
}
