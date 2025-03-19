package eventlistener

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"log"

	"github.com/TomasCruz/event-processing/internal/casino"
	"github.com/TomasCruz/event-processing/internal/config"
	"github.com/TomasCruz/event-processing/internal/ports"
	"google.golang.org/protobuf/proto"
)

type enricherSvc struct {
	config   config.Config
	db       ports.DB
	consumer ports.AsyncMsgConsumer
}

func (eSvc enricherSvc) run() {
	ctxProcess, cancelProcess := context.WithCancel(context.Background())
	ctxConsumer, cancelConsumer := context.WithCancel(context.Background())

	msgChan := make(chan []byte)

	go eSvc.process(ctxProcess, msgChan)
	go eSvc.consumer.Consume(ctxConsumer, msgChan)

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	eSvc.gracefulShutdown(cancelConsumer, cancelProcess)
}

func (eSvc enricherSvc) gracefulShutdown(cancelFunctions ...context.CancelFunc) {
	for _, cf := range cancelFunctions {
		cf()
	}

	// Kafka
	eSvc.consumer.Close()

	// DB
	eSvc.db.Close()
}

func (eSvc enricherSvc) process(ctxProcess context.Context, msgChan <-chan []byte) {
	wg := sync.WaitGroup{}
	defer func() { wg.Wait() }()

	for {
		var msg []byte
		var ok bool

		select {
		case <-ctxProcess.Done():
			return
		case msg, ok = <-msgChan:
			if !ok {
				return
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				var ev ports.Event
				err := proto.Unmarshal(msg, &ev)
				if err != nil {
					// ignore this event. Ideally should be sent to a dedicated error channel and processed asynchronously by support
					log.Printf("currency event Unmarshal failed: %v\n", err)
					return
				}

				// cutting corners here regarding publishing other events for enrichment of events;
				// doing it here instead. In a production situation these kind of enrichments would be addressed by separate components.
				// Also, event would be stored immediately on receipt, and would be updated after enrichment
				eSvc.enrichEvent(&ev)
				cEvent := ports.PBToEvent(&ev)

				// log enriched event in human readable form
				log.Printf("%s\n", humanReadable(cEvent))

				// store enriched event
				err = eSvc.db.CreateEvent(cEvent)
				if err != nil {
					log.Println("CreateEvent failed", err)
					return
				}
			}()
		}
	}
}

func humanReadable(ev casino.Event) string {
	b, err := json.Marshal(ev)
	if err != nil {
		return ""
	}

	return string(b)
}
