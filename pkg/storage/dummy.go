package storage

import (
	"go.uber.org/zap"
)

// dummyStorage object
type dummyStorage struct {
	mode   Mode
	logger *zap.Logger
}

// NewDummy returns Storager interface
func NewDummy(logger *zap.Logger) Storager {
	return &dummyStorage{
		mode:   DummyMode,
		logger: logger,
	}
}

// Save saves nothing
func (d *dummyStorage) Save(_ string) (bool, error) {
	d.logger.Debug("save", zap.String("mode", d.mode.String()))
	return true, nil
}

// Delete deletes nothing
func (d *dummyStorage) Delete() error {
	return nil
}

// Close closes nothing
func (d *dummyStorage) Close() {}
