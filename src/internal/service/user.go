package service

import (
	"context"
	"errors"
	dto "internship-task/pr-review-service/internal/dto"
	mapper "internship-task/pr-review-service/internal/mapper"
	repository "internship-task/pr-review-service/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService interface {
	SetUserActive(ctx context.Context, req *dto.SetUserActiveRequest) (*dto.UserResponse, error)
	GetUsersReviews(ctx context.Context, userName string) (*dto.UserPRsResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
	prRepo   repository.PullRequestRepository
}

func NewUserService(userRepo repository.UserRepository, prRepo repository.PullRequestRepository) UserService {
	return &userService{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

func (s *userService) SetUserActive(ctx context.Context, req *dto.SetUserActiveRequest) (*dto.UserResponse, error) {
	// Проверяем, существует ли пользователь
	userDB, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Меняем статус активности
	if err = s.userRepo.SetActive(ctx, req.UserID, req.IsActive); err != nil {
		return nil, err
	}

	userDB.IsActive = req.IsActive
	return mapper.ToUserResponse(userDB), nil
}

func (s *userService) GetUsersReviews(ctx context.Context, userID string) (*dto.UserPRsResponse, error) {
	prs, err := s.prRepo.GetByReviewerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var prsShort = make([]dto.PRShortResponse, 0)
	for _, pr := range prs {
		prsShort = append(prsShort, mapper.ToPRShortResponse(&pr))
	}
	return &dto.UserPRsResponse{
		UserID:       userID,
		PullRequests: prsShort,
	}, nil
}
