package storage

import (
	"errors"
	"github.com/igilrgyrg/english/model"
)

var ErrNoSavedWord = errors.New("no saved file")

type Storage interface {
	Save(w *model.Word) error
	Remove(w *model.Word) error
	IsExists(w *model.Word) (bool, error)
	PickRandom(username string) (*model.Word, error)
	UpdateRank(w *model.Word) error
}
