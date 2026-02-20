package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"github.com/user/fintech-expense-tracker/internal/db"
	"github.com/user/fintech-expense-tracker/internal/handler"
	"github.com/user/fintech-expense-tracker/internal/repository"
	"github.com/user/fintech-expense-tracker/internal/service"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" { dsn = "host=localhost user=postgres password=postgres dbname=fintech_tracker port=5432 sslmode=disable" }
	database := db.Init(dsn)
	repo := repository.NewGORMRepository(database)
	svc := service.NewExpenseService(repo)
	h := handler.NewHandler(svc)
	r := gin.Default()

	// Global Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health Check Diagnostic
	r.GET("/health", func(c *gin.Context) {
		status := "UP"
		dbErr := db.Ping(database)
		if dbErr != nil {
			status = "DATABASE_DOWN"
		}
		c.JSON(200, gin.H{
			"status":   status,
			"database": status != "DATABASE_DOWN",
			"version":  "1.0.0",
		})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/users", h.CreateUser)
		api.POST("/groups", h.CreateGroup)
		api.POST("/groups/:id/members", h.AddUserToGroup)
		api.POST("/groups/:id/expenses", h.AddExpense)
		api.GET("/groups/:id/balances", h.GetBalances)
		api.GET("/groups/:id/settlements", h.GetSettlements)
	}
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	r.Run(":" + port)
}
