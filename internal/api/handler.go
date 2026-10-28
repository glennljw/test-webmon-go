package api

import (
	"mini-webmon/internal/record"
	"mini-webmon/internal/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store store.Store
}

func NewHandler(s store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) History(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url, "records": h.store.History(url)})
}

type ingestRequest struct {
	Line string `json:"line" binding:"required"`
}

func (h *Handler) Ingest(c *gin.Context) {
	var req ingestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rec, err := record.ParseLine(req.Line)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.store.Add(rec)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
