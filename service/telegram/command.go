package telegram

import (
	"errors"
	"fmt"
	"github.com/igilrgyrg/english/e"
	"github.com/igilrgyrg/english/storage"
	"log"
	"strings"
)

const (
	StartCmd = "/start"
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	NewCmd   = "/new"
)

func (t TelegramProcessor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command <%s> from <%s>", text, username)

	switch text {
	case StartCmd:
		return t.client.SendMessage(chatID, fmt.Sprintf("Hello %s!\nWelcome to boost english words bot", username))
	case HelpCmd:
		return t.client.SendMessage(chatID, fmt.Sprintf("I can save and keep you english words to lear in future"))
	case RndCmd:
		return t.sendRandom(chatID, username)
	default:
		return t.client.SendMessage(chatID, msgUnknownCommand)
	}
}

func (t TelegramProcessor) sendRandom(chatID int, username string) (err error) {
	defer func() {
		err = e.WrapError("can`t do command: sendRandom", err)
	}()

	word, err := t.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedWord) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedWord) {
		return t.client.SendMessage(chatID, msgNoSavedWords)
	}

	if err := t.client.SendMessage(chatID, word.En); err != nil {
		return err
	}

	return t.storage.UpdateRank(word)
}
