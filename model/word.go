package model

import (
	"crypto/sha1"
	"fmt"
	"github.com/igilrgyrg/english/e"
	"io"
	"time"
)

type Word struct {
	En       string
	Ru       string
	Created  time.Time
	Rank     int8
	Username string
}

func (w Word) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, w.Username); err != nil {
		return "", e.WrapError("can`t hash username", err)
	}

	if _, err := io.WriteString(h, w.Created.String()); err != nil {
		return "", e.WrapError("can`t hash username", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
