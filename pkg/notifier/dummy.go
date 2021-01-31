package notifier

import (
	"go.uber.org/zap"

	"github.com/hiromaily/go-book-teacher/pkg/models"
)

// dummy object
type dummy struct {
	mode   Mode
	logger *zap.Logger
}

// NewDummy returns Notifier interface
func NewDummy(logger *zap.Logger) Notifier {
	return &dummy{
		mode:   DummyMode,
		logger: logger,
	}
}

// Send is to do nothing
func (d *dummy) Notify(ths []models.TeacherInfo) error {
	d.logger.Debug("notify", zap.String("mode", d.mode.String()))

	return nil
}
