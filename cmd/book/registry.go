package main

import (
	"github.com/hiromaily/go-book-teacher/pkg/booker"
	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/notifier"
	"github.com/hiromaily/go-book-teacher/pkg/siter"
	storages "github.com/hiromaily/go-book-teacher/pkg/storager"
)

// Registry is for registry interface
type Registry interface {
	NewBooker(string, int, int) booker.Booker
}

type registry struct {
	conf *config.Config
	// storager   *storages.Storager
}

// NewRegistry is to register regstry interface
func NewRegistry(conf *config.Config) Registry {
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

func (r *registry) newStorager() storages.Storager {
	storager, err := storages.NewStorager(r.conf)
	if err != nil {
		panic(err)
	}
	return storager
}

func (r *registry) newNotifier() notifier.Notifier {
	return notifier.NewNotifier(r.conf)
}

func (r *registry) newSiter(jsonPath string) siter.Siter {
	return siter.NewSiter(jsonPath, r.conf.Site)
}
