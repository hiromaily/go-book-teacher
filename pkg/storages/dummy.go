package storages

import (
	lg "github.com/hiromaily/golibs/log"
)

// DummyRepo is DummyRepo object
type DummyRepo struct {
	mode string
}

// NewDummy is return DummyRepo object
func NewDummy() *DummyRepo {
	return &DummyRepo{mode: "dummy"}
}

// Save is to do nothing
func (d *DummyRepo) Save(newData string) (bool, error) {
	lg.Debugf("Save by %s", d.mode)
	return true, nil
}

// Delete is to do nothing
func (d *DummyRepo) Delete() error {
	return nil
}

// Close is to do nothing
func (d *DummyRepo) Close() {}
