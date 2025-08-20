package main

import (
	"fmt"
	"log"
	"net/http"

	"blog-system/internal/auth"
	"blog-system/internal/handler"
	"blog-system/internal/repository"
	"blog-system/internal/service"
	"blog-system/pkg/config"
	"blog-system/pkg/database"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	db, err := database.NewSQLiteDB(cfg.DBPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	sessionManager := auth.NewSessionManager(cfg.SessionSecret)

	repo := repository.NewSQLiteRepository(db)
	blogService := service.NewBlogService(repo)
	blogHandler := handler.NewBlogHandler(blogService)
	authHandler := handler.NewAuthHandler(sessionManager, cfg.AdminPassword)

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	authHandler.RegisterRoutes(api)
	blogHandler.RegisterRoutes(api, sessionManager)

	r.Use(corsMiddleware)

	fmt.Printf("Blog server starting on :%d\n", cfg.Port)
	fmt.Printf("Database: %s\n", cfg.DBPath)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if (*r).Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
