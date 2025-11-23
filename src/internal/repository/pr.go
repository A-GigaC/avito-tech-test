	package repository

	import (
		"context"
		model "internship-task/pr-review-service/internal/model"
		gorm "gorm.io/gorm"
	)

	type PullRequestRepository interface {
		Create(ctx context.Context, pr *model.PullRequestDB) error
		GetByID(ctx context.Context, prID string) (*model.PullRequestDB, error)
		GetByReviewerID(ctx context.Context, reviewerID string) ([]model.PullRequestDB, error)
		Exists(ctx context.Context, prID string) (bool, error)
		Update(ctx context.Context, pr *model.PullRequestDB) error
	}

	type prRepo struct {
		db *gorm.DB
	}

	func NewPullRequestRepository(db *gorm.DB) PullRequestRepository {
		return &prRepo{db: db}
	}

	func (r *prRepo) Create(ctx context.Context, pr *model.PullRequestDB) error {
		return r.db.WithContext(ctx).Create(pr).Error
	}

	func (r *prRepo) GetByID(ctx context.Context, prID string) (*model.PullRequestDB, error) {
		var pr model.PullRequestDB
		err := r.db.WithContext(ctx).Where("pr_id = ?", prID).First(&pr).Error
		if err != nil {
			return nil, err
		}
		return &pr, nil
	}

	func (r *prRepo) GetByReviewerID(ctx context.Context, reviewerID string) ([]model.PullRequestDB, error) {
		var prs []model.PullRequestDB
		err := r.db.WithContext(ctx).Where("assigned_reviewers::jsonb @> ?", `["`+reviewerID+`"]`).Find(&prs).Error
		if err != nil {
			return nil, err
		}
		return prs, nil
	}

	func (r *prRepo) Exists(ctx context.Context, prID string) (bool, error) {
		var count int64
		err := r.db.Model(&model.PullRequestDB{}).WithContext(ctx).Where("pr_id = ?", prID).Count(&count).Error
		if err != nil {
			return false, err
		}
		return count > 0, nil
	}

	func (r *prRepo) Update(ctx context.Context, pr *model.PullRequestDB) error {
		return r.db.WithContext(ctx).Save(pr).Error
	}
