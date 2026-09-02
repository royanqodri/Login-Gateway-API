package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royanqodri/Login-Gateway-API/util"
)

func CountMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next() // process request first

		method := ctx.Request.Method // e.g. GET, POST
		path := ctx.FullPath()       // e.g. /hour-meter
		if path != "" {
			member := method + " " + path
			util.RedisClient.ZIncrBy(ctx, "api_usage", 1, member)
		}
	}
}

// Handler to get stats sorted (from Redis)
func StatsHandler(ctx *gin.Context) {
	// Get all members sorted by score (desc)
	res, err := util.RedisClient.ZRevRangeWithScores(ctx, "api_usage", 0, -1).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response
	type kv struct {
		Endpoint string  `json:"endpoint"`
		Count    float64 `json:"count"`
	}
	var data []kv
	for _, z := range res {
		data = append(data, kv{z.Member.(string), z.Score})
	}

	ctx.JSON(http.StatusOK, data)
}
