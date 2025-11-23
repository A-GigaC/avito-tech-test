package migrations

import (
    "internship-task/pr-review-service/internal/model"
    "gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
    // Создаем таблицы в правильном порядке, только если они не существуют
    
    if !db.Migrator().HasTable(&model.TeamDB{}) {
        if err := db.Migrator().CreateTable(&model.TeamDB{}); err != nil {
            return err
        }
    }
    
    if !db.Migrator().HasTable(&model.UserDB{}) {
        if err := db.Migrator().CreateTable(&model.UserDB{}); err != nil {
            return err
        }
    }
    
    if !db.Migrator().HasTable(&model.PullRequestDB{}) {
        if err := db.Migrator().CreateTable(&model.PullRequestDB{}); err != nil {
            return err
        }
    }
    
    return nil
}