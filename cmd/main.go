package main

import (
	"log"
	"net/http"

	// _ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"app/internal/config"
	"app/internal/controllers"
	"app/pkg/middleware"

	"github.com/gorilla/mux"
)

var cfg *config.Config

func init() {
	cfg = config.LoadConfig()
}

func main() {
	debug.SetMemoryLimit(500 * 1024 * 1024)

	r := mux.NewRouter()

	r.HandleFunc("/image/{path:(?:.*)}", controllers.ThumbController).Methods("GET")

	r.Use(middleware.Recovery)

	log.SetOutput(os.Stdout)
	log.Printf("Server starting on port %s\n", cfg.Server.Port)

	go func() {
		for {
			time.Sleep(1 * time.Minute) // Ожидание 1 минута
			runtime.GC()                // Принудительная сборка мусора
			debug.FreeOSMemory()        // Явно освобождаем память ОС
		}
	}()

	/*go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()*/

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
