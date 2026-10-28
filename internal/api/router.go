package api

import "github.com/gin-gonic/gin"

func NewRouter(h *Handler) *gin.Engine {
	r := gin.Default()
	r.GET("/health", h.Health)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/records", h.Ingest)
		v1.GET("/history", h.History)
	}
	return r
}
