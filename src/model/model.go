package model

import (
    "time"
)

// TeamDB - модель команды в базе данных
type TeamDB struct {
    TeamName  string    `gorm:"primaryKey;size:255"`
    CreatedAt time.Time
    UpdatedAt time.Time
    // FK
    Members   []UserDB  `gorm:"foreignKey:TeamName;references:TeamName"`
}

// UserDB - модель пользователя в базе данных
type UserDB struct {
    UserID   string `gorm:"primaryKey;size:255"`
    Username string `gorm:"size:255;not null"`
    TeamName string `gorm:"size:255;not null"`
    IsActive bool   `gorm:"default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
    // FK
    Team TeamDB `gorm:"foreignKey:TeamName;references:TeamName"`
}

// PullRequestDB - модель Pull Request в базе данных
type PullRequestDB struct {
    PRID      string    `gorm:"primaryKey;size:255"`
    PRName    string    `gorm:"size:255;not null"`
    AuthorID  string    `gorm:"size:255;not null"`
    Status    PRStatus  `gorm:"type:varchar(20);default:OPEN"`
    CreatedAt time.Time
    MergedAt  time.Time
    // Простой список ID пользователей (0-2 элемента)
    AssignedReviewers []string `gorm:"type:jsonb"`
    // FK
    Author UserDB `gorm:"foreignKey:AuthorID;references:UserID"`
}

// PRStatus - статус Pull Request
type PRStatus string

const (
    StatusOpen   PRStatus = "OPEN"
    StatusMerged PRStatus = "MERGED"
)