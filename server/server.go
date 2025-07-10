package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"thesilentcoder.com/m/health"
	"thesilentcoder.com/m/url"
)

func Start(ctx context.Context, config Config) error {
	repository := url.NewRepository()
	urlService := url.New(repository, config.Port, config.RedirectUrl, config.ApiPrefix, config.ApiVersion)
	healthService := health.New()
	services := []Service{urlService, healthService}

	httpHandler := mux.NewRouter()
	for _, service := range services {
		service.RegisterHandlers(httpHandler)
	}

	httpServer := &http.Server{Addr: config.Port, Handler: httpHandler}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	listenChan := make(chan error)
	go func() {
		fmt.Printf("Starting HTTP server on %s\n", config.Port)
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			// couldn't Start server
			listenChan <- err
		}
	}()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	select {
	case listenErr := <-listenChan:
		return fmt.Errorf("failed to Start HTTP server: %w", listenErr)
	case <-interruptChan:
		// Interrupt signal, need to shut down gracefully
		break
	case <-ctx.Done():
		// Start context cancelled, need to shut down gracefully
		break
	}
	log.Info().Msg("Shutting down...")
	if err := httpServer.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("failed to shut down HTTP server: %w", err)
	}
	// Graceful shutdown
	return nil
}

type Service interface {
	RegisterHandlers(mux *mux.Router)
}
