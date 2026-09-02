package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gmorinn/technical_test/internal/config"
	"github.com/gmorinn/technical_test/internal/db"
)

type FizzBuzzStore interface {
	UpsertFizzBuzzStat(ctx context.Context, arg db.UpsertFizzBuzzStatParams) error
	GetTopFizzBuzzStat(ctx context.Context) (db.GetTopFizzBuzzStatRow, error)
}

type Handler struct {
	Config *config.API
	DB     FizzBuzzStore
}

func New(cfg *config.API, store FizzBuzzStore) *Handler {
	return &Handler{Config: cfg, DB: store}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {

	v1 := r.Group("/api/v1")

	fizzbuzz := v1.Group("/fizzbuzz")
	{
		fizzbuzz.GET("", h.FizzBuzz)
		fizzbuzz.GET("/stats", h.FizzBuzzStats)
	}
}
