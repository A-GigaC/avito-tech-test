package service

import (
	"context"
	"errors"
	dto "internship-task/pr-review-service/internal/dto"
	mapper "internship-task/pr-review-service/internal/mapper"
	repository "internship-task/pr-review-service/internal/repository"
)

var (
	ErrTeamExists = errors.New("team already exists")
)

type TeamService interface {
	CreateTeam(ctx context.Context, req *dto.TeamRequest) (*dto.TeamResponse, error)
	GetTeam(ctx context.Context, teamName string) (*dto.TeamResponse, error)
}

type teamService struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository) TeamService {
	return &teamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *teamService) CreateTeam(ctx context.Context, req *dto.TeamRequest) (*dto.TeamResponse, error) {
	exists, err := s.teamRepo.Exists(ctx, req.TeamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrTeamExists
	}

	teamDB, usersDB := mapper.ToTeamDB(req)

	if err := s.teamRepo.Create(ctx, teamDB); err != nil {
		return nil, err
	}

	for _, user := range usersDB {
		if err := s.userRepo.CreateOrUpdate(ctx, &user); err != nil {
			return nil, err
		}
	}

	return mapper.ToTeamResponse(teamDB, usersDB), nil
}

func (s *teamService) GetTeam(ctx context.Context, teamName string) (*dto.TeamResponse, error) {
	teamDB, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	return mapper.ToTeamResponse(teamDB, teamDB.Members), nil
}
