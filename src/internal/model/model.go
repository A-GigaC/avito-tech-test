package model

import (
    "time"
    "fmt"
    "encoding/json"
    "database/sql/driver"
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
    //Team TeamDB `gorm:"foreignKey:TeamName;references:TeamName"`
    Team TeamDB `gorm:"-"`
}

type PullRequestDB struct {
	PRID               string    `gorm:"column:pr_id;primaryKey"`
	PRName             string    `gorm:"column:pr_name"`
	AuthorID           string    `gorm:"column:author_id"`
	Status             PRStatus    `gorm:"column:status"`
	CreatedAt          time.Time `gorm:"column:created_at"`
	MergedAt           time.Time `gorm:"column:merged_at"`
	AssignedReviewers  JSONArray `gorm:"column:assigned_reviewers;type:json"`
}

type JSONArray []string

func (a *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*a = JSONArray{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON value: %v", value)
	}
	
	return json.Unmarshal(bytes, a)
}

func (a JSONArray) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return json.Marshal(a)
}

type PRStatus string

const (
    StatusOpen   PRStatus = "OPEN"
    StatusMerged PRStatus = "MERGED"
)