package save

import (
	"go.uber.org/zap"
)

// dummySaver object
type dummySaver struct {
	mode   Mode
	logger *zap.Logger
}

// NewDummy returns Saver interface
func NewDummySaver(logger *zap.Logger) Saver {
	return &dummySaver{
		mode:   DummyMode,
		logger: logger,
	}
}

// Save saves nothing
func (d *dummySaver) Save(_ string) (bool, error) {
	d.logger.Debug("save", zap.String("mode", d.mode.String()))
	return true, nil
}

// Delete deletes nothing
func (d *dummySaver) Delete() error {
	return nil
}

// Close closes nothing
func (d *dummySaver) Close() {}
