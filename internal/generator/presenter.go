package generator

import (
	"context"
	"log"
	"time"

	"github.com/TomasCruz/event-processing/internal/ports"
	"google.golang.org/protobuf/proto"
)

type generatorPresenter struct {
	topic string
	db    ports.DB
	p     ports.AsyncMsgProducer
}

func (g generatorPresenter) run() {
	ctx, cancel := context.WithTimeout(context.Background(), 17*time.Second)
	defer func() { cancel(); g.p.Close() }()

	eventCh := Generate(ctx)
	for event := range eventCh {
		ev := ports.EventToPB(event)
		b, err := proto.Marshal(ev)
		if err != nil {
			log.Println("event Marshal failed", err)
			continue
		}

		cEvent := ports.PBToEvent(ev)
		log.Printf("%v\n", cEvent)

		// architectural note: in production scenario, the following would be extracted
		// to a service layer function (which would live in generator/service)
		// That structure is made in eventstats for demonstration
		// service then has all the dependencies injected as interfaces, making it testable
		if err := g.p.SendAsyncMsg(b); err != nil {
			log.Printf("failed to produce %s event: %v\n", g.topic, err)
		}
	}
}
