package storages

import (
	lg "github.com/hiromaily/golibs/log"
)

// DummyRepo object
type DummyRepo struct{}

// NewDummy
func NewDummy() *DummyRepo {
	return &DummyRepo{}
}

// Save is to save data to text
func (d *DummyRepo) Save(newData string) (bool, error) {
	lg.Debug("Using Dummy")
	return true, nil
}

// Delete is to delete file
func (d *DummyRepo) Delete() error {
	return nil
}

func (d *DummyRepo) Close() {}
