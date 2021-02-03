package booker

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-book-teacher/pkg/notifier"
	"github.com/hiromaily/go-book-teacher/pkg/save"
	"github.com/hiromaily/go-book-teacher/pkg/site"
	"github.com/hiromaily/go-book-teacher/pkg/teachers"
)

// ----------------------------------------------------------------------------
// Booker interface
// ----------------------------------------------------------------------------

// Booker interface
type Booker interface {
	Start() error
	Clean()
	Close()
}

// NewBooker returns Booker interface
func NewBooker(
	saver save.Saver,
	notifier notifier.Notifier,
	siter site.Siter,
	logger *zap.Logger,
	interval int,
) Booker {
	return NewBook(
		saver,
		notifier,
		siter,
		logger,
		interval,
	)
}

// ----------------------------------------------------------------------------
// Book
// ----------------------------------------------------------------------------

// Book object
type Book struct {
	saver    save.Saver
	notifier notifier.Notifier
	siter    site.Siter
	logger   *zap.Logger
	interval int
	isLoop   bool
}

// NewBook is to return book object
func NewBook(
	saver save.Saver,
	notifier notifier.Notifier,
	siter site.Siter,
	logger *zap.Logger,
	interval int,
) *Book {
	book := Book{
		saver:    saver,
		notifier: notifier,
		siter:    siter,
		logger:   logger,
		interval: interval,
		isLoop:   interval != 0, // Note: testmode, heroku env should be false
	}
	return &book
}

// Start starts book execution
func (b *Book) Start() error {
	b.logger.Debug("book Start()")
	defer b.saver.Close()

	// fetch initial teacher data
	b.logger.Debug("book siter.Fetch()")
	if err := b.siter.Fetch(); err != nil {
		return errors.Wrap(err, "fail to call siter.Fetch()")
	}

	for {
		// scraping
		b.logger.Debug("book siter.FindTeachers()")
		teachers := b.siter.FindTeachers()

		// save
		b.logger.Debug("book siter.save()")
		isUpdated, err := b.save(teachers)
		if err != nil {
			b.logger.Error("fail to call save()", zap.Error(err))
		}
		if isUpdated {
			// notify
			b.notifier.Notify(teachers)
		}

		// execute only once
		if !b.isLoop {
			return nil
		}

		b.logger.Debug(fmt.Sprintf("sleep %d second for next execution", b.interval))
		time.Sleep(time.Duration(b.interval) * time.Second)
	}
}

// Clean deletes save
func (b *Book) Clean() {
	b.saver.Delete()
}

// Close closes middleware
func (b *Book) Close() {
	b.saver.Close()
}

// save saves status on storage
func (b *Book) save(teachers []teachers.TeacherRepo) (bool, error) {
	if len(teachers) == 0 {
		return false, nil
	}

	// create string from ids slice
	var sum int
	for _, t := range teachers {
		sum += t.ID
	}
	newData := strconv.Itoa(sum)

	// save
	isUpdated, err := b.saver.Save(newData)
	if err != nil {
		return false, errors.Wrap(err, "fail to call Save()")
	}
	return isUpdated, nil
}
