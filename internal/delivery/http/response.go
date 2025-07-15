package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"test_task/internal/models"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func (h *Handler) newErrorResponse(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrInvalidInput):
		c.AbortWithStatus(http.StatusBadRequest)
	case errors.Is(err, models.ErrNotFound):
		c.AbortWithStatus(http.StatusNotFound)
	case errors.Is(err, models.ErrUnauthorized):
		c.AbortWithStatus(http.StatusUnauthorized)
	default:
		h.logger.Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
