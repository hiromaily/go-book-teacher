package dummy

import (
	"github.com/hiromaily/go-book-teacher/pkg/models"
)

type DummySite struct{}

func NewDummySite() *DummySite {
	return &DummySite{}
}

func (d *DummySite) FetchInitialData() error {
	return nil
}

func (d *DummySite) InitializeSavedTeachers() {}

func (d *DummySite) HandleTeachers() {}

func (d *DummySite) GetSavedTeachers() []models.TeacherInfo {
	return nil
}
