package eventstats

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TomasCruz/event-processing/internal/config"
	"github.com/TomasCruz/event-processing/internal/database"
	"github.com/TomasCruz/event-processing/internal/eventstats/presenter"
)

func Run() {
	// populate configuration
	c, err := config.ConfigFromEnvVars()
	if err != nil {
		log.Fatal("failed to read environment variables", err)
	}

	// init DB
	db, err := database.New(c)
	if err != nil {
		log.Fatal(err, "failed to initialize database")
	}

	// run
	presenter := presenter.StatsPresenter{
		Config: c,
		DB:     db,
		Server: &http.Server{Addr: fmt.Sprintf(":%s", c.Port)},
	}
	presenter.Run()
}
