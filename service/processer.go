package service

import "github.com/igilrgyrg/english"

type Processor interface {
	Process(e english.Event) error
}

type Fetcher interface {
	Fetch(limit int) ([]english.Event, error)
}
