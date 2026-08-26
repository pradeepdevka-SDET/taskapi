package main

import (
	"database/sql"
	"log"
	"os"
	"taskapi/internal/handler"
	"taskapi/internal/repository"
	"taskapi/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found,using real environment variables")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("open db:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("connect db:", err)
	}
	log.Println("Connected to database!")
	// -- BUILD the chain bottom-up: repos-> services->handlers->---
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	taskService := service.NewTaskService(taskRepo)

	authHandler := handler.NewAuthHandler(userRepo, jwtSecret)
	taskHandler := handler.NewTaskHandler(taskService)

	// --- ROUTER --
	r := gin.Default()

	r.POST("/signup", authHandler.Signup)
	r.POST("/login", authHandler.Login)

	tasks := r.Group("/tasks")
	tasks.Use(handler.AuthMiddleware(jwtSecret)) //every /tasks route needs a valid token
	{
		tasks.POST("", taskHandler.Create)
		tasks.GET("", taskHandler.List)
		tasks.GET("/:id", taskHandler.Get)
		tasks.PATCH("/:id", taskHandler.Update)
		tasks.DELETE("/:id", taskHandler.Delete)
	}
	log.Println("listerning on :" + port)
	r.Run(":" + port)
}
