package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"test_task/internal/models"
)

func (h *Handler) createSubscription(c *gin.Context) {
	var input models.Subscription
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, err)
		return
	}
	id, err := h.services.Subscription.Create(c.Request.Context(), &input)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) listSubscriptions(c *gin.Context) {
	subs, err := h.services.Subscription.List(c.Request.Context())
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, subs)
}

func (h *Handler) getSubscriptionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	sub, err := h.services.Subscription.GetByID(c.Request.Context(), id)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, sub)
}

func (h *Handler) updateSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	var input models.UpdateSubscriptionInput
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, err)
		return
	}

	if err := h.services.Subscription.Update(c.Request.Context(), id, &input); err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) deleteSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	if err := h.services.Subscription.Delete(c.Request.Context(), id); err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) getSummary(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	from := c.Query("from")
	to := c.Query("to")

	if userID == "" || serviceName == "" || from == "" || to == "" {
		h.newErrorResponse(c, gin.Error{Err: http.ErrMissingFile})
		return
	}

	sum, err := h.services.Subscription.SumByPeriod(c.Request.Context(), userID, serviceName, from, to)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": sum})
}
