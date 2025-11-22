package handler

import (
	dto "internship-task/pr-review-service/dto"
	service "internship-task/pr-review-service/service"
	http "net/http"

	gin "github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) AddTeam(c *gin.Context) {
	var req dto.TeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teamResp, err := h.teamService.CreateTeam(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrTeamExists {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "TEAM_EXISTS",
					Message: "team_name already exists",
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

	c.JSON(http.StatusCreated, gin.H{"team": teamResp})
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "BAD_REQUEST",
				Message: "team_name is required",
			},
		})
		return
	}

	teamResp, err := h.teamService.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "NOT_FOUND",
				Message: "team not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, teamResp)
}
