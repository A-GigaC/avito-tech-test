package handler

import (
	dto "internship-task/pr-review-service/internal/dto"
	service "internship-task/pr-review-service/internal/service"
	http "net/http"

	gin "github.com/gin-gonic/gin"
)

type PullRequestHandler struct {
	pullRequestService service.PullRequestService
}

func NewPullRequestHandler(pullRequestService service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{
		pullRequestService: pullRequestService,
	}
}

func (h *PullRequestHandler) CreatePullRequest(c *gin.Context) {
	var req dto.CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pullRequestResp, err := h.pullRequestService.CreatePullRequest(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrPullRequestExists {
			c.JSON(http.StatusConflict, dto.ErrorPRExists)
			return
		} else if err == service.ErrResourceNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorInternal)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"pullRequest": pullRequestResp})
}

func (h *PullRequestHandler) MergePullRequest(c *gin.Context) {
	var req dto.MergePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pullRequestResp, err := h.pullRequestService.MergePullRequest(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrResourceNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorInternal)
		return
	}

	c.JSON(http.StatusOK, gin.H{"pullRequest": pullRequestResp})
}

func (h *PullRequestHandler) ReassignReviewer(c *gin.Context) {
	var req dto.ReassignReviewerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reassignResp, err := h.pullRequestService.ReassignReviewer(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrResourceNotFound:
			c.JSON(http.StatusNotFound, dto.ErrorNotFound)
		case service.ErrPRMerged:
			c.JSON(http.StatusConflict, dto.ErrorPRMerged)
		case service.ErrReviewerNotAssigned:
			c.JSON(http.StatusConflict, dto.ErrorNotAssigned)
		case service.ErrNoCandidate:
			c.JSON(http.StatusConflict, dto.ErrorNoCandidate)
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorInternal)
		}
		return
	}

	c.JSON(http.StatusOK, reassignResp)
}
