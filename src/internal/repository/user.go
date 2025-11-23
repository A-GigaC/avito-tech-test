package repository

import (
	"context"
	model "internship-task/pr-review-service/internal/model"

	gorm "gorm.io/gorm"
)

type UserRepository interface {
	CreateOrUpdate(ctx context.Context, user *model.UserDB) error
	GetByID(ctx context.Context, userID string) (*model.UserDB, error)
	SetActive(ctx context.Context, userID string, isActive bool) error
	GetActiveTeamMembers(ctx context.Context, teamName string, excludeUserID string) ([]model.UserDB, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateOrUpdate(ctx context.Context, user *model.UserDB) error {
	// Используем Save, который обновляет, если существует, или создает нового
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepo) GetByID(ctx context.Context, userID string) (*model.UserDB, error) {
	var user model.UserDB
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) SetActive(ctx context.Context, userID string, isActive bool) error {
	return r.db.WithContext(ctx).Model(&model.UserDB{}).Where("user_id = ?", userID).Update("is_active", isActive).Error
}

func (r *userRepo) GetActiveTeamMembers(ctx context.Context, teamName string, excludeUserID string) ([]model.UserDB, error) {
	var users []model.UserDB
	err := r.db.WithContext(ctx).Where("team_name = ? AND is_active = ? AND user_id != ?", teamName, true, excludeUserID).Find(&users).Error
	return users, err
}
