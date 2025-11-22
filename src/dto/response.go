package dto

import "time"

// TeamResponse - ответ с данными команды
type TeamResponse struct {
    TeamName string         `json:"team_name"`
    Members  []UserResponse `json:"members"`
}

// UserResponse - данные пользователя в ответе
type UserResponse struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    IsActive bool   `json:"is_active"`
}

// UserWithTeamResponse - пользователь с информацией о команде
type UserWithTeamResponse struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    TeamName string `json:"team_name"`
    IsActive bool   `json:"is_active"`
}

// PRResponse - полные данные PR
type PRResponse struct {
    PRID              string     `json:"pull_request_id"`
    PRName            string     `json:"pull_request_name"`
    AuthorID          string     `json:"author_id"`
    Status            string     `json:"status"`
    AssignedReviewers []string   `json:"assigned_reviewers"`
    CreatedAt         time.Time  `json:"createdAt"`
    MergedAt          time.Time `json:"mergedAt,omitempty"`
}

// PRShortResponse - сокращенные данные PR
type PRShortResponse struct {
    PRID     string `json:"pull_request_id"`
    PRName   string `json:"pull_request_name"`
    AuthorID string `json:"author_id"`
    Status   string `json:"status"`
}

// ReassignResponse - ответ на переназначение ревьювера
type ReassignResponse struct {
    PR         PRResponse `json:"pr"`
    ReplacedBy string     `json:"replaced_by"`
}

// UserPRsResponse - ответ с PR пользователя
type UserPRsResponse struct {
    UserID       string            `json:"user_id"`
    PullRequests []PRShortResponse `json:"pull_requests"`
}

// ErrorResponse - стандартный ответ об ошибке
type ErrorResponse struct {
    Error struct {
        Code    string `json:"code"`
        Message string `json:"message"`
    } `json:"error"`
}

