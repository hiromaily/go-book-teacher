package notifier

import (
	"github.com/hiromaily/go-book-teacher/pkg/models"
	lg "github.com/hiromaily/golibs/log"
)

// NewDummy is to return Dummy object
func NewDummy() *Dummy {
	return &Dummy{mode: "dummy"}
}

// Dummy is Dummy object
type Dummy struct {
	mode string
}

// Send is to do nothing
func (d *Dummy) Send(ths []models.TeacherInfo) error {
	lg.Debugf("Send by %s", d.mode)

	return nil
}
