package main

import (
	"github.com/pkg/errors"

	"github.com/hiromaily/go-book-teacher/pkg/booker"
	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/notifier"
	"github.com/hiromaily/go-book-teacher/pkg/site"
	"github.com/hiromaily/go-book-teacher/pkg/storage"
)

// Registry is for registry interface
type Registry interface {
	NewBooker(string, int, int) booker.Booker
}

type registry struct {
	conf *config.Root
	// storage   *storages.Storager
}

// NewRegistry is to register regstry interface
func NewRegistry(conf *config.Root) Registry {
	return &registry{conf: conf}
}

// NewBooker is to register for booker interface
func (r *registry) NewBooker(jsonPath string, day, interval int) booker.Booker {
	return booker.NewBooker(
		r.conf,
		day,
		interval,
		r.newStorager(),
		r.newNotifier(),
		r.newSiter(jsonPath),
	)
}

func (r *registry) newStorager() storage.Storager {
	storager, err := storage.NewStorager(
		r.conf.Storage.Mode,
		r.conf.Storage.Redis.URL,
		r.conf.Storage.Text.Path,
	)
	if err != nil {
		panic(err)
	}
	return storager
}

func (r *registry) newNotifier() notifier.Notifier {
	switch r.conf.Notification.Mode {
	case notifier.ConsoleMode:
		return notifier.NewConsole()
	case notifier.SlackMode:
		return notifier.NewSlack(r.conf.Notification.Slack.Key, r.conf.Site.URL)
	}

	panic(errors.New("invalid notification mode"))
}

func (r *registry) newSiter(jsonPath string) site.Siter {
	return site.NewSiter(jsonPath, r.conf.Site)
}
