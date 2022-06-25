package consumer

import (
	"github.com/igilrgyrg/english"
	"github.com/igilrgyrg/english/service"
	"log"
	"time"
)

type Consumer interface {
	Start() error
}

type EventConsumer struct {
	processor service.Processor
	fetcher   service.Fetcher
	batchSize int
}

func New(fetcher service.Fetcher, processor service.Processor, batchSize int) Consumer {
	return &EventConsumer{fetcher: fetcher, processor: processor, batchSize: batchSize}
}

func (e EventConsumer) Start() error {
	for {
		gotEvents, err := e.fetcher.Fetch(e.batchSize)
		if err != nil {
			log.Printf("[ERROR] consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(time.Second * 1)
			continue
		}

		if err := e.handleEvents(gotEvents); err != nil {
			log.Println("can`t handle events: %w", err)
			continue
		}
	}
}

func (e *EventConsumer) handleEvents(events []english.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := e.processor.Process(event); err != nil {
			log.Printf("can`t handle event: %s", event.Text)

			continue
		}
	}
	return nil
}
