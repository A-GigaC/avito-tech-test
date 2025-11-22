package mappers

import (
	dto "internship-task/pr-review-service/dto"
	model "internship-task/pr-review-service/model"
	"time"
)

// ToTeamDB преобразует TeamRequest в TeamDB и []UserDB
func ToTeamDB(req *dto.TeamRequest) (*model.TeamDB, []model.UserDB) {
	teamDB := &model.TeamDB{
		TeamName: req.TeamName,
	}

	usersDB := make([]model.UserDB, len(req.Members))
	for i, member := range req.Members {
		usersDB[i] = model.UserDB{
			UserID:   member.UserID,
			Username: member.Username,
			TeamName: req.TeamName,
			IsActive: member.IsActive,
		}
	}

	return teamDB, usersDB
}

// ToTeamResponse преобразует TeamDB и []UserDB в TeamResponse
func ToTeamResponse(teamDB *model.TeamDB, usersDB []model.UserDB) *dto.TeamResponse {
	members := make([]dto.UserResponse, len(usersDB))
	for i, user := range usersDB {
		members[i] = dto.UserResponse{
			UserID:   user.UserID,
			Username: user.Username,
			IsActive: user.IsActive,
		}
	}

	return &dto.TeamResponse{
		TeamName: teamDB.TeamName,
		Members:  members,
	}
}

// ToPRDB преобразует CreatePRRequest в PullRequestDB
func ToPRDB(req *dto.CreatePRRequest, reviewers []string) *model.PullRequestDB {
	return &model.PullRequestDB{
		PRID:              req.PRID,
		PRName:            req.PRName,
		AuthorID:          req.AuthorID,
		Status:            model.StatusOpen,
		AssignedReviewers: reviewers,
		CreatedAt:         time.Now(),
	}
}

// ToPRResponse преобразует PullRequestDB в PRResponse
func ToPRResponse(prDB *model.PullRequestDB) *dto.PRResponse {
	return &dto.PRResponse{
		PRID:              prDB.PRID,
		PRName:            prDB.PRName,
		AuthorID:          prDB.AuthorID,
		Status:            string(prDB.Status),
		AssignedReviewers: prDB.AssignedReviewers,
		CreatedAt:         prDB.CreatedAt,
		MergedAt:          prDB.MergedAt,
	}
}

// ToPRShortResponse преобразует PullRequestDB в PRShortResponse
func ToPRShortResponse(prDB *model.PullRequestDB) dto.PRShortResponse {
	return dto.PRShortResponse{
		PRID:     prDB.PRID,
		PRName:   prDB.PRName,
		AuthorID: prDB.AuthorID,
		Status:   string(prDB.Status),
	}
}

// ToUserResponse преобразует UserDB в UserResponse
func ToUserResponse(userDB *model.UserDB) *dto.UserResponse {
	return &dto.UserResponse{
		UserID:   userDB.UserID,
		Username: userDB.Username,
		IsActive: userDB.IsActive,
	}
}

// ToUserWithTeamResponse преобразует UserDB в UserWithTeamResponse
func ToUserWithTeamResponse(userDB *model.UserDB) *dto.UserWithTeamResponse {
	return &dto.UserWithTeamResponse{
		UserID:   userDB.UserID,
		Username: userDB.Username,
		TeamName: userDB.TeamName,
		IsActive: userDB.IsActive,
	}
}
