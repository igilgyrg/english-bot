package main

import (
	"flag"
	telegramclient "github.com/igilrgyrg/english/clients/telegram"
	"github.com/igilrgyrg/english/consumer"
	telegramservice "github.com/igilrgyrg/english/service/telegram"
	"github.com/igilrgyrg/english/storage/files"
	"log"
)

const (
	tgHost      = "api.telegram.org"
	pathStorage = "storage"
	batchSize   = 100
)

func main() {
	token := mustToken()

	tgClient := telegramclient.New(tgHost, token)

	processor := telegramservice.New(tgClient, files.New(pathStorage))

	consumer := consumer.New(processor, processor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("Error of starting handler telegram bot")
	}
}

func mustToken() string {
	token := flag.String("token", "", "bot token for access to telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is required")
	}

	return *token
}
