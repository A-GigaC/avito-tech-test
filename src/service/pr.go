package service

import (
	"context"
	"errors"
	dto "internship-task/pr-review-service/dto"
	mapper "internship-task/pr-review-service/mapper"
	model "internship-task/pr-review-service/model"
	repository "internship-task/pr-review-service/repository"
	"math/rand"
	"time"
)

var (
	ErrPullRequestExists   = errors.New("PR id already exists")
	ErrResourceNotFound    = errors.New("resource not found")
	ErrPRMerged            = errors.New("cannot reassign on merged PR")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate         = errors.New("no active replacement candidate in team")
)

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, req *dto.CreatePRRequest) (*dto.PRResponse, error)
	MergePullRequest(ctx context.Context, req *dto.MergePRRequest) (*dto.PRResponse, error)
	ReassignReviewer(ctx context.Context, req *dto.ReassignReviewerRequest) (*dto.ReassignResponse, error)
}

type pullRequestService struct {
	prRepo   repository.PullRequestRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewPullRequestService(prRepo repository.PullRequestRepository, userRepo repository.UserRepository, teamRepo repository.TeamRepository) PullRequestService {
	return &pullRequestService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *pullRequestService) CreatePullRequest(ctx context.Context, req *dto.CreatePRRequest) (*dto.PRResponse, error) {
	// Проверяем существование pr с таким же id
	pr, err := s.prRepo.GetByID(ctx, req.PRID)
	if err != nil {
		return nil, err
	} else if pr.PRID == req.PRID {
		return nil, ErrPullRequestExists
	}
	// Проверяем, существует ли пользователь
	userDB, err := s.userRepo.GetByID(ctx, req.AuthorID)
	if err != nil {
		return nil, err
	} else if userDB.UserID != req.AuthorID {
		return nil, ErrResourceNotFound
	}
	// Проверяем существование команды (по идее она и так должна быть, но на всякий случай)
	teamExists, err := s.teamRepo.Exists(ctx, userDB.TeamName)
	if err != nil {
		return nil, err
	} else if teamExists == false {
		return nil, ErrResourceNotFound
	}
	// Получаем всех активных тиммейтов автора
	activeMembers, err := s.userRepo.GetActiveTeamMembers(ctx, userDB.TeamName, req.AuthorID)
	if err != nil {
		return nil, err
	}
	// Назначаем ревьюеров
	assignedReviewers := make([]string, 0, 2)
	if len(activeMembers) > 0 {
		indices := rand.Perm(len(activeMembers))
		count := min(2, len(indices))
		for i := 0; i < count; i++ {
			assignedReviewers = append(assignedReviewers, activeMembers[indices[i]].UserID)
		}
	}
	// Получаем объект PRDB
	prDB := mapper.ToPRDB(req, assignedReviewers)
	// наконец-то оздаем PR
	err = s.prRepo.Create(ctx, mapper.ToPRDB(req, assignedReviewers))
	if err != nil {
		return nil, err
	}

	// Формируем и возвращаем ответ
	return mapper.ToPRResponse(prDB), nil
}

func (s *pullRequestService) MergePullRequest(ctx context.Context, req *dto.MergePRRequest) (*dto.PRResponse, error) {
	// Проверяем существования pr с данным id
	prExists, err := s.prRepo.Exists(ctx, req.PRID)
	if err != nil {
		return nil, err
	} else if prExists == false {
		return nil, ErrResourceNotFound
	}
	// Получаем pr
	prDB, err := s.prRepo.GetByID(ctx, req.PRID)
	if err != nil {
		return nil, err
	}
	if prDB.Status != model.StatusMerged {
		prDB.Status = model.StatusMerged
		prDB.MergedAt = time.Now()
		err = s.prRepo.Update(ctx, prDB)
		if err != nil {
			return nil, err
		}
	}

	return mapper.ToPRResponse(prDB), nil
}

func (s *pullRequestService) ReassignReviewer(ctx context.Context, req *dto.ReassignReviewerRequest) (*dto.ReassignResponse, error) {
	// Получаем PR по ID
	prDB, err := s.prRepo.GetByID(ctx, req.PRID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	// Проверяем, что PR не в статусе MERGED
	if prDB.Status == model.StatusMerged {
		return nil, ErrPRMerged
	}

	// Проверяем, что старый ревьювер назначен на этот PR
	oldReviewerFound := false
	for _, reviewer := range prDB.AssignedReviewers {
		if reviewer == req.OldUserID {
			oldReviewerFound = true
			break
		}
	}
	if !oldReviewerFound {
		return nil, ErrReviewerNotAssigned
	}

	// Получаем автора PR
	author, err := s.userRepo.GetByID(ctx, prDB.AuthorID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	// Получаем всех активных тиммейтов, исключая автора
	activeMembers, err := s.userRepo.GetActiveTeamMembers(ctx, author.TeamName, prDB.AuthorID)
	if err != nil {
		return nil, err
	}

	// Фильтруем: убираем уже назначенных ревьюверов и старого ревьювера
	var candidates []model.UserDB
	for _, member := range activeMembers {
		// Пропускаем старого ревьювера
		if member.UserID == req.OldUserID {
			continue
		}

		// Проверяем, не назначен ли уже этот пользователь ревьювером
		alreadyAssigned := false
		for _, reviewer := range prDB.AssignedReviewers {
			if reviewer == member.UserID {
				alreadyAssigned = true
				break
			}
		}

		if !alreadyAssigned {
			candidates = append(candidates, member)
		}
	}

	// Проверяем, есть ли кандидаты
	if len(candidates) == 0 {
		return nil, ErrNoCandidate
	}

	// Выбираем случайного кандидата
	newReviewer := candidates[rand.Intn(len(candidates))]

	// Обновляем список ревьюверов в PR
	newReviewers := make([]string, len(prDB.AssignedReviewers))
	for i, reviewer := range prDB.AssignedReviewers {
		if reviewer == req.OldUserID {
			newReviewers[i] = newReviewer.UserID
		} else {
			newReviewers[i] = reviewer
		}
	}

	// Обновляем PR в базе
	prDB.AssignedReviewers = newReviewers
	err = s.prRepo.Update(ctx, prDB)
	if err != nil {
		return nil, err
	}

	// Формируем ответ
	prResponse := mapper.ToPRResponse(prDB)
	return &dto.ReassignResponse{
		PR:         *prResponse,
		ReplacedBy: newReviewer.UserID,
	}, nil
}
