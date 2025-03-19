package presenter

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/TomasCruz/event-processing/internal/eventstats/service"
)

func (sPres StatsPresenter) materialized(w http.ResponseWriter, r *http.Request) {
	svc := service.StatsSvc{
		Config: sPres.Config,
		DB:     sPres.DB,
	}

	md, err := svc.Materialized()
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	serialized, err := json.Marshal(md)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	io.Writer.Write(w, serialized)
}
