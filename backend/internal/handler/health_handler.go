package handler

import (
	"github.com/aryasatyawa/bayarin/internal/pkg/database"
	"github.com/aryasatyawa/bayarin/internal/pkg/redis"
	"github.com/aryasatyawa/bayarin/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	db    *database.PostgresDB
	redis *redis.RedisClient
}

func NewHealthHandler(db *database.PostgresDB, redis *redis.RedisClient) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

// Health godoc
// @Summary Health check
// @Description Check API health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	// Check database
	dbStatus := "ok"
	if err := h.db.Health(c.Request.Context()); err != nil {
		dbStatus = "error"
	}

	// Check Redis
	redisStatus := "ok"
	if err := h.redis.Health(c.Request.Context()); err != nil {
		redisStatus = "error"
	}

	response.Success(c, "API is running", gin.H{
		"status":   "ok",
		"database": dbStatus,
		"redis":    redisStatus,
	})
}
