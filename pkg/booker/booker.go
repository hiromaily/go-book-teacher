package booker

import (
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/models"
	"github.com/hiromaily/go-book-teacher/pkg/notifier"
	"github.com/hiromaily/go-book-teacher/pkg/siter"
	storages "github.com/hiromaily/go-book-teacher/pkg/storager"
	lg "github.com/hiromaily/golibs/log"
)

// ----------------------------------------------------------------------------
// Booker interface
// ----------------------------------------------------------------------------

// Booker is interface
type Booker interface {
	Start() error
	Cleanup()
	Close()
}

// NewBooker is to return booker interface
func NewBooker(
	conf *config.Config,
	day int,
	interval int,
	storager storages.Storager,
	notifier notifier.Notifier,
	siter siter.Siter) Booker {
	return NewBook(conf, day, interval, storager, notifier, siter)
}

// ----------------------------------------------------------------------------
// Book
// ----------------------------------------------------------------------------

// Book is Book object
type Book struct {
	conf     *config.Config
	day      int
	interval int
	storager storages.Storager
	notifier notifier.Notifier
	siter    siter.Siter
	isLoop   bool
}

// NewBook is to return book object
func NewBook(
	conf *config.Config,
	day int,
	interval int,
	storager storages.Storager,
	notifier notifier.Notifier,
	siter siter.Siter) *Book {
	var isLoop bool
	if interval != 0 {
		isLoop = true
	}

	book := Book{
		conf:     conf,
		day:      day,
		interval: interval,
		storager: storager,
		notifier: notifier,
		siter:    siter,
		isLoop:   isLoop, // TODO: testmode, heroku env should be false
	}
	return &book
}

// Start is to start book execution
func (b *Book) Start() error {
	if err := b.siter.FetchInitialData(); err != nil {
		return errors.Wrap(err, "fail to call siter.FetchInitialData()")
	}

	for {
		// scraping
		teachers := b.siter.FindTeachers(b.day)

		// save
		b.saveAndNotify(teachers)

		// execute only once
		if !b.isLoop {
			b.storager.Close()
			return nil
		}

		time.Sleep(time.Duration(b.interval) * time.Second)
	}
}

// Cleanup is to clean up middleware object
func (b *Book) Cleanup() {
	b.storager.Delete()
}

// Close is to clean up middleware object
func (b *Book) Close() {
	b.storager.Close()
}

// saveAndNotify is to save and notify if something saved
func (b *Book) saveAndNotify(ths []models.TeacherInfo) {
	if len(ths) != 0 {
		// create string from ids slice
		var sum int
		for _, t := range ths {
			sum += t.ID
		}
		newData := strconv.Itoa(sum)

		// save
		isUpdated, err := b.storager.Save(newData)
		if err != nil {
			lg.Errorf("fail to save() %v", err)
		}

		if isUpdated {
			// notify
			b.notifier.Send(ths)
		}
	}
}

// ----------------------------------------------------------------------------
// DummyBook
// ----------------------------------------------------------------------------

// DummyBook is DummyBook object
type DummyBook struct{}

// NewDummyBook is to return NewDummyBook object
func NewDummyBook() *DummyBook {
	return &DummyBook{}
}

// Start is to do nothing
func (b *DummyBook) Start() error {
	return nil
}

// Cleanup is to do nothing
func (b *DummyBook) Cleanup() {}

// Close is to do nothing
func (b *DummyBook) Close() {}
