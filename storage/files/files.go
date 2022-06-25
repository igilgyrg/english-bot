package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/igilrgyrg/english/e"
	"github.com/igilrgyrg/english/model"
	"github.com/igilrgyrg/english/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const defaultPerm = 0774

type FileStorage struct {
	basePath string
}

func New(path string) storage.Storage {
	return &FileStorage{basePath: path}
}

func (s *FileStorage) Save(w *model.Word) (err error) {
	defer func() { err = e.WrapError("can`t save word", err) }()

	filePath := filepath.Join(s.basePath, w.Username)

	if err = os.MkdirAll(filePath, defaultPerm); err != nil {
		return
	}

	fName, err := fileName(w)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err = gob.NewEncoder(file).Encode(w); err != nil {
		return err
	}

	return nil
}

func (s *FileStorage) Remove(w *model.Word) error {
	filename, err := fileName(w)
	if err != nil {
		return e.WrapError("can`t open file", err)
	}

	path := filepath.Join(s.basePath, w.Username, filename)

	if err := os.Remove(path); err != nil {
		errMsg := fmt.Sprintf("can`t delete file path <%s>", path)
		return e.WrapError(errMsg, err)
	}

	return nil
}

func (s *FileStorage) IsExists(w *model.Word) (bool, error) {
	filename, err := fileName(w)
	if err != nil {
		return false, e.WrapError("can`t open file", err)
	}

	path := filepath.Join(s.basePath, w.Username, filename)

	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.WrapError("error of exists file", err)
	}

	return true, nil
}

func (s *FileStorage) PickRandom(username string) (w *model.Word, err error) {
	defer func() { err = e.WrapError("can`t pick random word", err) }()

	path := filepath.Join(s.basePath, username)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedWord
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodeWord(filepath.Join(path, file.Name()))
}

func (s *FileStorage) UpdateRank(w *model.Word) error {
	return nil
}

func (s *FileStorage) decodeWord(filePath string) (w *model.Word, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.WrapError("can`t decode file", err)
	}
	defer func() { _ = f.Close() }()

	if err := gob.NewDecoder(f).Decode(w); err != nil {
		return nil, e.WrapError("can`t decode file to word", err)
	}

	return
}

func fileName(w *model.Word) (string, error) {
	return w.Hash()
}
