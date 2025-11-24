package dto

// TeamRequest - запрос на создание/обновление команды
type TeamRequest struct {
    TeamName string         `json:"team_name" binding:"required,min=1,max=255"`
    Members  []UserRequest  `json:"members" binding:"required,min=1,dive"`
}

// UserRequest - данные пользователя в запросе
type UserRequest struct {
    UserID   string `json:"user_id" binding:"required,min=1,max=255"`
    Username string `json:"username" binding:"required,min=1,max=255"`
    IsActive bool   `json:"is_active"`
}

// CreatePRRequest - запрос на создание PR
type CreatePRRequest struct {
    PRID     string `json:"pull_request_id" binding:"required,min=1,max=255"`
    PRName   string `json:"pull_request_name" binding:"required,min=1,max=255"`
    AuthorID string `json:"author_id" binding:"required,min=1,max=255"`
}

// MergePRRequest - запрос на мерж PR
type MergePRRequest struct {
    PRID string `json:"pull_request_id" binding:"required,min=1,max=255"`
}

// ReassignReviewerRequest - запрос на переназначение ревьювера
type ReassignReviewerRequest struct {
    PRID      string `json:"pull_request_id" binding:"required,min=1,max=255"`
    OldUserID string `json:"old_reviewer_id" binding:"required,min=1,max=255"`
}

// SetUserActiveRequest - запрос на изменение активности пользователя
type SetUserActiveRequest struct {
    UserID  string `json:"user_id" binding:"required,min=1,max=255"`
    IsActive bool  `json:"is_active"`
}