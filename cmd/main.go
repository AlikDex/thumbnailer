package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"app/internal/config"
	"app/internal/controllers"
	"app/pkg/middleware"

	"github.com/gorilla/mux"
)

var (
	cfg *config.Config
)

func init() {
	cfg = config.GetConfig()
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/{path:[^.]+\\.(?:jpg|jpeg|png|webp)}", controllers.ImageController).Methods("GET")

	r.Use(middleware.Recovery)

	log.SetOutput(os.Stdout)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
