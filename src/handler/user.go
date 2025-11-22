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
	// Получаем тело запроса
	var req dto.SetUserActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем статус активности пользователя
	user, err := h.userService.SetUserActive(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *UserHandler) GetReview(c *gin.Context) {
	userID := c.Query("user_id")

	userReviews, err := h.userService.GetUsersReviews(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, userReviews)
}
