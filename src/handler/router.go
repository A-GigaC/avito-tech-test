package handler

import (
    "github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, 
    teamHandler *TeamHandler,
    userHandler *UserHandler, 
    prHandler *PullRequestHandler) {
    
    teamRoutes := router.Group("/team")
    {
        teamRoutes.POST("/add", teamHandler.AddTeam)
        teamRoutes.GET("/get", teamHandler.GetTeam)
    }
    
    userRoutes := router.Group("/users")
    {
        userRoutes.POST("/setIsActive", userHandler.SetUserActive)
        userRoutes.GET("/getReview", userHandler.GetReview)
    }
    
    prRoutes := router.Group("/pullRequest")
    {
        prRoutes.POST("/create", prHandler.CreatePullRequest)
        prRoutes.POST("/merge", prHandler.MergePullRequest)
        prRoutes.POST("/reassign", prHandler.ReassignReviewer)
    }
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
}