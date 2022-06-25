package telegram

import (
	"errors"
	"github.com/igilrgyrg/english"
	"github.com/igilrgyrg/english/clients/telegram"
	"github.com/igilrgyrg/english/e"
	"github.com/igilrgyrg/english/storage"
)

var (
	ErrUnknownType  = errors.New("unknown event type")
	ErrTelegramMeta = errors.New("error get telegram meta from event")
)

type TelegramProcessor struct {
	client  *telegram.Client
	offset  int
	storage storage.Storage
}

type TelegramMeta struct {
	ChatID   int
	Username string
}

func New(client *telegram.Client, storage storage.Storage) *TelegramProcessor {
	return &TelegramProcessor{client: client, offset: 0, storage: storage}
}

func (t TelegramProcessor) Process(e english.Event) error {
	switch e.Type {
	case english.Message:
		return t.processMessage(e)
	default:
		return ErrUnknownType
	}
}

func (t *TelegramProcessor) Fetch(limit int) ([]english.Event, error) {
	updates, err := t.client.Updates(t.offset, limit)
	if err != nil {
		return nil, e.WrapError("can`t get updates", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]english.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, ToEvent(u))
	}

	t.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (t TelegramProcessor) processMessage(event english.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.WrapError("can`t process message", err)
	}

	if err := t.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.WrapError("can`t process message", err)
	}

	return nil
}

func meta(event english.Event) (*TelegramMeta, error) {
	res, ok := event.Meta.(TelegramMeta)
	if !ok {
		return nil, ErrTelegramMeta
	}
	return &res, nil
}

func ToEvent(update telegram.Update) english.Event {
	typeEvent := fetchType(update)

	res := english.Event{
		Type: typeEvent,
		Text: fetchText(update),
	}

	if typeEvent == english.Message {
		res.Meta = TelegramMeta{ChatID: update.Message.Chat.ID, Username: update.Message.From.Username}
	}

	return res
}

func fetchType(update telegram.Update) english.Type {
	if update.Message == nil {
		return english.Unknown
	}

	return english.Message
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}

	return update.Message.Text
}
