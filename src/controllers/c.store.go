package controllers

import (
	"log"

	"github.com/dirental/core/db"
	"github.com/dirental/core/src/helpers"
	helperdb "github.com/dirental/core/src/helpers/db"
	middleware "github.com/dirental/core/src/middlewares"
	"github.com/dirental/core/src/types/requests"
	"github.com/gin-gonic/gin"
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
func (c *IStoreController) CreateStore(ctx *gin.Context) {
	data, err := helpers.ValidateRequest[requests.IRequestCreateStore](ctx)
	if err != nil {
		return
	}
	auth, err := middleware.GetAuthorized(ctx)
	if err != nil {
		log.Print(err)
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := helperdb.StringToUUID(auth.UserID)
	if err != nil {
		log.Print(err)
		ctx.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := c.queries.GetUser(ctx, userID)
	if err != nil {
		log.Print(err)
		ctx.JSON(404, gin.H{"error": "User not found", "auth": auth})
		return
	}

	store := db.CreateStoreParams{
		Name:           data.Name,
		Phone:          data.Phone,
		Coordinate:     helperdb.SafeString(&data.Coordinate),
		Address:        "Feature Needed",
		IsActive:       false,
		Contacts:       data.Contacts,
		Description:    helperdb.SafeString(&data.Description),
		Category:       data.Category,
		TermAndService: data.TermAndService,
		OfOwner:        user.ID,
		Logo:           helperdb.SafeString(&data.Logo),
	}

	newStore, err := c.queries.CreateStore(ctx, store)
	if err != nil {
		log.Print(err)
		ctx.JSON(500, gin.H{"error": "Failed to create store"})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Store created successfully",
		"tenant":  newStore,
		"user":    user,
	})
}
