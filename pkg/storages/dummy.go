package storages

import (
	lg "github.com/hiromaily/golibs/log"
)

// DummyRepo object
type DummyRepo struct {
	mode string
}

// NewDummy
func NewDummy() *DummyRepo {
	return &DummyRepo{mode: "dummy"}
}

// Save is to save data to text
func (d *DummyRepo) Save(newData string) (bool, error) {
	lg.Debugf("Save by %s", d.mode)
	return true, nil
}

// Delete is to delete file
func (d *DummyRepo) Delete() error {
	return nil
}

func (d *DummyRepo) Close() {}
