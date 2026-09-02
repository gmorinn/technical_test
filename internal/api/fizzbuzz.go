package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gmorinn/technical_test/internal/db"
	"github.com/gmorinn/technical_test/internal/fizzbuzz"
)

type fizzBuzzParams struct {
	Int1  int    `form:"int1" binding:"required,min=1,max=1000"`
	Int2  int    `form:"int2" binding:"required,min=1,max=1000"`
	Limit int    `form:"limit" binding:"required,min=1,max=1000"`
	Str1  string `form:"str1" binding:"required"`
	Str2  string `form:"str2" binding:"required"`
}

// FizzBuzz godoc
//
//	@Summary		FizzBuzz
//	@Description	Returns a fizzbuzz list from 1 to limit: multiples of int1 become str1, multiples of int2 become str2, multiples of both become str1+str2.
//	@Tags			fizzbuzz
//	@Produce		json
//	@Param			int1	query		int		true	"First divisor"						minimum(1)	maximum(1000)
//	@Param			int2	query		int		true	"Second divisor"					minimum(1)	maximum(1000)
//	@Param			limit	query		int		true	"Upper bound, inclusive"			minimum(1)	maximum(1000)
//	@Param			str1	query		string	true	"Replacement for int1 multiples"	example(fizz)
//	@Param			str2	query		string	true	"Replacement for int2 multiples"	example(buzz)
//	@Success		200		{object}	map[string][]string
//	@Failure		400		{object}	map[string]string
//	@Router			/fizzbuzz [get]
func (h *Handler) FizzBuzz(c *gin.Context) {
	var p fizzBuzzParams
	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := fizzbuzz.Generate(p.Int1, p.Int2, p.Limit, p.Str1, p.Str2)

	if err := h.DB.UpsertFizzBuzzStat(c.Request.Context(), db.UpsertFizzBuzzStatParams{
		Int1:  int32(p.Int1),
		Int2:  int32(p.Int2),
		Limit: int32(p.Limit),
		Str1:  p.Str1,
		Str2:  p.Str2,
	}); err != nil {
		log.Printf("fizzbuzz: recording stat failed: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// FizzBuzzStats godoc
//
//	@Summary		FizzBuzz statistics
//	@Description	Returns the most requested fizzbuzz parameters and how many times they have been asked for. Responds with {"message": "no requests yet"} when nothing has been recorded.
//	@Tags			fizzbuzz
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]string
//	@Router			/fizzbuzz/stats [get]
func (h *Handler) FizzBuzzStats(c *gin.Context) {
	row, err := h.DB.GetTopFizzBuzzStat(c.Request.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusOK, gin.H{"message": "no requests yet"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"int1":  row.Int1,
		"int2":  row.Int2,
		"limit": row.Limit,
		"str1":  row.Str1,
		"str2":  row.Str2,
		"hits":  row.Hits,
	})
}
