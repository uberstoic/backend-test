package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Serhio/backend-vk-test/internal/models"
	"github.com/Serhio/backend-vk-test/internal/service"
	"github.com/gin-gonic/gin"
)

type createAdInput struct {
	Title    string  `json:"title" binding:"required"`
	Text     string  `json:"text" binding:"required"`
	ImageURL string  `json:"image_url" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
}

func (h *Handler) createAd(c *gin.Context) {
	userID, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	var input createAdInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ad := &models.Ad{
		Title:    input.Title,
		Text:     input.Text,
		ImageURL: input.ImageURL,
		Price:    input.Price,
		UserID:   userID.(int64),
	}

	createdAd, err := h.services.Ad.CreateAd(c.Request.Context(), ad)
	if err != nil {
		if errors.Is(err, service.ErrTitleTooLong) || errors.Is(err, service.ErrTextTooLong) || errors.Is(err, service.ErrPriceInvalid) || errors.Is(err, service.ErrImageURLInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, createdAd)
}

func (h *Handler) getAds(c *gin.Context) {
	// try to get user ID from context, if authorized
	var currentUserID int64 = -1 // default value for unauthorized
	if id, ok := c.Get(userCtx); ok {
		if userID, ok := id.(int64); ok {
			currentUserID = userID
		}
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := c.DefaultQuery("sort_dir", "desc")
	minPrice, _ := strconv.ParseFloat(c.DefaultQuery("min_price", "0"), 64)
	maxPrice, _ := strconv.ParseFloat(c.DefaultQuery("max_price", "0"), 64)

	// validate sort parameters
	if !(sortBy == "created_at" || sortBy == "price") {
		sortBy = "created_at"
	}
	if !(sortDir == "asc" || sortDir == "desc") {
		sortDir = "desc"
	}

	ads, err := h.services.Ad.GetAds(c.Request.Context(), page, limit, sortBy, sortDir, minPrice, maxPrice, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, ads)
}
