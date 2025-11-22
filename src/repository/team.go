package repository

import (
	"context"
	model "internship-task/pr-review-service/model"

	gorm "gorm.io/gorm"
)

type TeamRepository interface {
	Create(ctx context.Context, team *model.TeamDB) error
	GetByName(ctx context.Context, teamName string) (*model.TeamDB, error)
	Exists(ctx context.Context, teamName string) (bool, error)
}

type teamRepo struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepo{db: db}
}

func (r *teamRepo) Create(ctx context.Context, team *model.TeamDB) error {
	return r.db.WithContext(ctx).Create(team).Error
}

func (r *teamRepo) GetByName(ctx context.Context, teamName string) (*model.TeamDB, error) {
	var team model.TeamDB
	// Preload Members, чтобы получить и пользователей команды
	err := r.db.WithContext(ctx).Preload("Members").Where("team_name = ?", teamName).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepo) Exists(ctx context.Context, teamName string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.TeamDB{}).Where("team_name = ?", teamName).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
