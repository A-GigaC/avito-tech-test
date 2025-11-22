package service

import (
	"context"
	"errors"
	dto "internship-task/pr-review-service/dto"
	mapper "internship-task/pr-review-service/mapper"
	repository "internship-task/pr-review-service/repository"
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
	// Проверяем, существует ли команда
	exists, err := s.teamRepo.Exists(ctx, req.TeamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrTeamExists
	}

	// Преобразуем DTO в модели БД
	teamDB, usersDB := mapper.ToTeamDB(req)

	// Начинаем транзакцию, чтобы сохранить команду и пользователей
	// Для этого нам понадобится передать транзакцию в репозитории, но для упрощения предположим, что у нас нет транзакций
	// Вместо этого, мы сохраним команду, а затем пользователей

	// Сохраняем команду
	if err := s.teamRepo.Create(ctx, teamDB); err != nil {
		return nil, err
	}

	// Сохраняем пользователей
	for _, user := range usersDB {
		// Используем userRepo для сохранения каждого пользователя
		if err := s.userRepo.CreateOrUpdate(ctx, &user); err != nil {
			return nil, err
		}
	}

	// Преобразуем обратно в ответ
	return mapper.ToTeamResponse(teamDB, usersDB), nil
}

func (s *teamService) GetTeam(ctx context.Context, teamName string) (*dto.TeamResponse, error) {
	teamDB, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	// Преобразуем в ответ
	return mapper.ToTeamResponse(teamDB, teamDB.Members), nil
}
