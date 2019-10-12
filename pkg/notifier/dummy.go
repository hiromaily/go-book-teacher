package notifier

import (
	"github.com/hiromaily/go-book-teacher/pkg/models"
)

// NewDummy is to create Dummy
func NewDummy() *Dummy {
	return &Dummy{}
}

// Dummy is Dummy object
type Dummy struct{}

func (d *Dummy) Send(ths []models.TeacherInfo) error {
	return nil
}
