package presenter

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TomasCruz/event-processing/internal/config"
	"github.com/TomasCruz/event-processing/internal/ports"
)

type StatsPresenter struct {
	Config config.Config
	DB     ports.DB
	Server *http.Server
}

func (sPres StatsPresenter) Run() {
	http.HandleFunc("/materialized", sPres.materialized)

	go func() {
		if err := sPres.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	sPres.gracefulShutdown()
}

func (sPres StatsPresenter) gracefulShutdown() {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := sPres.Server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	// DB
	sPres.DB.Close()
}
