package app

import (
	"context"
	"fmt"
	"github.com/Hermes-Bird/betera-test-task/app/config"
	"github.com/Hermes-Bird/betera-test-task/app/handlers"
	"github.com/Hermes-Bird/betera-test-task/app/repos"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func Start() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.GetConfig()

	apodRepo := repos.NewApodPostgresRepo(cfg)
	imageRepo := repos.NewImageDropboxRepo(cfg)

	apodHandler := handlers.NewApodHandler(apodRepo, imageRepo)

	mux := http.NewServeMux()
	mux.Handle("/apods", apodHandler)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: mux,
	}

	go func() {
		log.Printf("Server listening on port %s", cfg.ServerPort)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("error while listen and serve:", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down server")

	sdCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	server.Shutdown(sdCtx)
	apodRepo.Close()
}
