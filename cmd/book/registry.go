package main

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-book-teacher/pkg/booker"
	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/logger"
	"github.com/hiromaily/go-book-teacher/pkg/notifier"
	"github.com/hiromaily/go-book-teacher/pkg/site"
	"github.com/hiromaily/go-book-teacher/pkg/site/dmmer"
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
		r.newStorager(),
		r.newNotifier(),
		r.newSiter(jsonPath),
		r.newLogger(),
		day,
		r.conf.Interval,
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
		switch r.conf.Storage.Mode {
		case storage.RedisMode:
			r.newLogger().Debug("storager: redis")
			r.storage, err = storage.NewRedis(r.newLogger(), r.conf.Storage.Redis.URL)
		case storage.TextMode:
			r.newLogger().Debug("storager: text")
			r.storage = storage.NewText(r.newLogger(), r.conf.Storage.Text.Path)
		case storage.DummyMode:
			r.newLogger().Debug("storager: dummy")
			r.storage = storage.NewDummy(r.newLogger())
		default:
			panic(errors.New("storage mode is not found"))
		}
		if err != nil {
			panic(err)
		}
	}
	return r.storage
}

func (r *registry) newNotifier() notifier.Notifier {
	switch r.conf.Notification.Mode {
	case notifier.ConsoleMode:
		r.newLogger().Debug("notifier: console")
		return notifier.NewConsole(r.newLogger())
	case notifier.SlackMode:
		r.newLogger().Debug("notifier: slack")
		return notifier.NewSlack(
			r.newLogger(),
			r.conf.Notification.Slack.Key,
			r.conf.Site.URL,
		)
	}
	panic(errors.Errorf("invalid notification mode: %s", r.conf.Notification.Mode))
}

func (r *registry) newSiter(jsonPath string) site.Siter {
	switch r.conf.Site.Type {
	case site.SiteTypeDMM:
		r.newLogger().Debug("site: dmm")
		return dmmer.NewDMM(
			r.newLogger(),
			jsonPath,
			r.conf.Site.URL,
		)
	}
	panic(errors.Errorf("invalid site type: %s", r.conf.Site.Type))
}
