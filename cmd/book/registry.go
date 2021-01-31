package main

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-book-teacher/pkg/booker"
	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/logger"
	"github.com/hiromaily/go-book-teacher/pkg/notifier"
	"github.com/hiromaily/go-book-teacher/pkg/site"
	"github.com/hiromaily/go-book-teacher/pkg/storage"
)

// Registry interface
type Registry interface {
	NewBooker(string, int) booker.Booker
}

type registry struct {
	conf    *config.Root
	logger  *zap.Logger
	storage storage.Storager
}

// NewRegistry is to register regstry interface
func NewRegistry(conf *config.Root) Registry {
	return &registry{conf: conf}
}

// NewBooker is to register for booker interface
func (r *registry) NewBooker(jsonPath string, day int) booker.Booker {
	return booker.NewBooker(
		day,
		r.conf.Interval,
		r.newStorager(),
		r.newNotifier(),
		r.newSiter(jsonPath),
	)
}

func (r *registry) newLogger() *zap.Logger {
	if r.logger == nil {
		r.logger = logger.NewZapLogger(r.conf.Logger)
	}
	return r.logger
}

func (r *registry) newStorager() storage.Storager {
	if r.storage == nil {
		var err error
		r.storage, err = storage.NewStorager(
			r.conf.Storage.Mode,
			r.newLogger(),
			r.conf.Storage.Redis.URL,
			r.conf.Storage.Text.Path,
		)
		if err != nil {
			panic(err)
		}
	}
	return r.storage
}

func (r *registry) newNotifier() notifier.Notifier {
	switch r.conf.Notification.Mode {
	case notifier.ConsoleMode:
		return notifier.NewConsole(r.newLogger())
	case notifier.SlackMode:
		return notifier.NewSlack(
			r.newLogger(),
			r.conf.Notification.Slack.Key,
			r.conf.Site.URL,
		)
	}
	panic(errors.Errorf("invalid notification mode: %s", r.conf.Notification.Mode))
}

func (r *registry) newSiter(jsonPath string) site.Siter {
	return site.NewSiter(jsonPath, r.conf.Site)
}
