package handler

import (
	"errors"
	"net/http"

	"github.com/Serhio/backend-vk-test/internal/service"
	"github.com/Serhio/backend-vk-test/internal/storage/postgres"
	"github.com/gin-gonic/gin"
)

type registerInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) register(c *gin.Context) {
	var input registerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.services.User.Register(c.Request.Context(), input.Login, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) || errors.Is(err, service.ErrLoginFormat) || errors.Is(err, service.ErrPasswordFormat) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) login(c *gin.Context) {
	var input registerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.services.User.Login(c.Request.Context(), input.Login, input.Password)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) || errors.Is(err, service.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
