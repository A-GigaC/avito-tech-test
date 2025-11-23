package handler

import (
	dto "internship-task/pr-review-service/internal/dto"
	service "internship-task/pr-review-service/internal/service"
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
			c.JSON(http.StatusBadRequest, dto.ErrorTeamExists)
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorInternal)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"team": teamResp})
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorBadRequest)
		return
	}

	teamResp, err := h.teamService.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorNotFound)
		return
	}

	c.JSON(http.StatusOK, teamResp)
}
