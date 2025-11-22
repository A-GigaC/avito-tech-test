package handler

import (
	dto "internship-task/pr-review-service/dto"
	service "internship-task/pr-review-service/service"
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
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "PR_EXISTS",
					Message: "PR id already exists",
				},
			})
			return
		} else if err == service.ErrResourceNotFound {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
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
			c.JSON(http.StatusConflict, dto.ErrorResponse{
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

	c.JSON(http.StatusCreated, gin.H{"pullRequest": pullRequestResp})
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
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "NOT_FOUND",
					Message: "PR or user not found",
				},
			})
		case service.ErrPRMerged:
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "PR_MERGED",
					Message: "cannot reassign on merged PR",
				},
			})
		case service.ErrReviewerNotAssigned:
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "NOT_ASSIGNED",
					Message: "reviewer is not assigned to this PR",
				},
			})
		case service.ErrNoCandidate:
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "NO_CANDIDATE",
					Message: "no active replacement candidate in team",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "INTERNAL_ERROR",
					Message: err.Error(),
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, reassignResp)
}
