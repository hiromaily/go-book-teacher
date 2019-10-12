package storages

import (
	"github.com/hiromaily/go-book-teacher/pkg/config"
)

// Saver is to save
type Saver interface {
}

// Deleter is to Notice
type Deleter interface {
}

// SaveDeleter is to Save and Delete
type Storager interface {
	Save(string) (bool, error)
	Delete() error
	Close()
}

func NewStorager(conf *config.Config) (Storager, error) {
	if conf.ValidateRedis() {
		//redis mode
		rd, err := NewRedis(conf.Redis.URL)
		if err == nil {
			return rd, nil
		}
	}

	if conf.ValidateText() {
		//text mode
		return NewText(conf.Text.Path), nil
	}

	return NewDummy(), nil
}
