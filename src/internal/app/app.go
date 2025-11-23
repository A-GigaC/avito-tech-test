package app

import (
    "log"
    "internship-task/pr-review-service/internal/config"
    "internship-task/pr-review-service/internal/handler"
    "internship-task/pr-review-service/internal/migrations"
    "internship-task/pr-review-service/internal/repository"
    "internship-task/pr-review-service/internal/service"
    "internship-task/pr-review-service/pkg/database"

    "github.com/gin-gonic/gin"
)

func Run() {
    // Загрузка конфигурации
    cfg := config.Load()
    
    // Подключение к базе данных
    db, err := database.NewGorm(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    // Автомиграции
    if err := migrations.AutoMigrate(db); err != nil {
        log.Fatal("Failed to run migrations:", err)
    }
    
    // Инициализация репозиториев
    teamRepo := repository.NewTeamRepository(db)
    userRepo := repository.NewUserRepository(db)
    prRepo := repository.NewPullRequestRepository(db)
    
    // Инициализация сервисов
    teamService := service.NewTeamService(teamRepo, userRepo)
    userService := service.NewUserService(userRepo, prRepo)
    prService := service.NewPullRequestService(prRepo, userRepo, teamRepo)
    
    // Инициализация хендлеров
    teamHandler := handler.NewTeamHandler(teamService)
    userHandler := handler.NewUserHandler(userService)
    prHandler := handler.NewPullRequestHandler(prService)
    
    // Настройка роутера
    router := gin.Default()
    
    // Routes
    router.POST("/team/add", teamHandler.AddTeam)
    router.GET("/team/get", teamHandler.GetTeam)
    router.POST("/users/setIsActive", userHandler.SetUserActive)
    router.GET("/users/getReview", userHandler.GetReview)
    router.POST("/pullRequest/create", prHandler.CreatePullRequest)
    router.POST("/pullRequest/merge", prHandler.MergePullRequest)
    router.POST("/pullRequest/reassign", prHandler.ReassignReviewer)
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // Запуск сервера
    log.Printf("Server starting on port %s", cfg.Port)
    if err := router.Run(":" + cfg.Port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}