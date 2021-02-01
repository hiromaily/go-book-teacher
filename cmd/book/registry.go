package main

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-book-teacher/pkg/booker"
	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/logger"
	"github.com/hiromaily/go-book-teacher/pkg/notifier"
	"github.com/hiromaily/go-book-teacher/pkg/save"
	"github.com/hiromaily/go-book-teacher/pkg/site"
	"github.com/hiromaily/go-book-teacher/pkg/site/dmm"
	"github.com/hiromaily/go-book-teacher/pkg/teachers"
)

// Registry interface
type Registry interface {
	NewBooker(string, int) booker.Booker
}

type registry struct {
	conf   *config.Root
	logger *zap.Logger
	saver  save.Saver
}

// NewRegistry is to register regstry interface
func NewRegistry(conf *config.Root) Registry {
	return &registry{conf: conf}
}

// NewBooker is to register for booker interface
func (r *registry) NewBooker(jsonFile string, day int) booker.Booker {
	return booker.NewBooker(
		r.newSaver(),
		r.newNotifier(),
		r.newSiter(jsonFile, day),
		r.newLogger(),
		r.conf.Interval,
	)
}

func (r *registry) newLogger() *zap.Logger {
	if r.logger == nil {
		r.logger = logger.NewZapLogger(r.conf.Logger)
	}
	return r.logger
}

func (r *registry) newSaver() save.Saver {
	if r.saver == nil {
		var err error
		switch r.conf.Save.Mode {
		case save.RedisMode:
			r.newLogger().Debug("storager: redis")
			r.saver, err = save.NewRedisSaver(r.newLogger(), r.conf.Save.Redis.URL)
		case save.TextMode:
			r.newLogger().Debug("storager: text")
			r.saver = save.NewTextSaver(r.newLogger(), r.conf.Save.Text.Path)
		case save.DummyMode:
			r.newLogger().Debug("storager: dummy")
			r.saver = save.NewDummySaver(r.newLogger())
		default:
			panic(errors.New("save mode is not found"))
		}
		if err != nil {
			panic(err)
		}
	}
	return r.saver
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

func (r *registry) newSiter(jsonFile string, day int) site.Siter {
	switch r.conf.Site.Type {
	case site.SiteTypeDMM:
		r.newLogger().Debug("site: dmm")
		return dmm.NewDMM(
			r.newLogger(),
			r.newTeacherFetcher(jsonFile),
			r.conf.Site.URL,
			day,
		)
	}
	panic(errors.Errorf("invalid site type: %s", r.conf.Site.Type))
}

func (r *registry) newTeacherFetcher(jsonFile string) teachers.Teacher {
	if jsonFile != "" {
		r.newLogger().Debug("target teachers: json")
		return teachers.NewJSONTeacher(
			r.newLogger(),
			jsonFile,
		)
	}
	r.newLogger().Debug("target teachers: dummy")
	return teachers.NewDummyTeacher(r.newLogger())
}
