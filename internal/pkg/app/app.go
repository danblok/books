package app

import (
	"log"
	"os"

	"github.com/danblok/books/internal/app/db"
	"github.com/danblok/books/internal/app/handlers"
	"github.com/danblok/books/internal/app/repos"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type App struct {
	handlers *handlers.BooksHandlers
	repos    *repos.BooksRepo
}

func New() (*App, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	db, err := db.New(db.Config{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DB:       os.Getenv("POSTGRES_DB"),
	})
	if err != nil {
		return nil, err
	}

	repos := repos.New(db)
	handlers := handlers.New(repos)

	return &App{
		repos:    repos,
		handlers: handlers,
	}, nil
}

func (a *App) Run() {
	g := gin.Default()
	a.handlers.RegisterHandlers("/api/", g)
	log.Fatalln(g.Run(":" + os.Getenv("API_PORT")))
}
