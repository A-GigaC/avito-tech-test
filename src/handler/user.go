package handler

import (
	dto "internship-task/pr-review-service/dto"
	service "internship-task/pr-review-service/service"
	http "net/http"

	gin "github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) SetUserActive(c *gin.Context) {
	var req dto.SetUserActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.SetUserActive(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorInternal)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) GetReview(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorBadRequest)
		return
	}

	userReviews, err := h.userService.GetUsersReviews(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorInternal)
		return
	}

	c.JSON(http.StatusOK, userReviews)
}