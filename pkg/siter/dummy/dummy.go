package dummy

import (
	"github.com/hiromaily/go-book-teacher/pkg/models"
)

// DummySite is DummySite object
type DummySite struct{}

// NewDummySite is to return DummySite object
func NewDummySite() *DummySite {
	return &DummySite{}
}

// FetchInitialData is to do nothing
func (d *DummySite) FetchInitialData() error {
	return nil
}

// InitializeSavedTeachers is to do nothing
func (d *DummySite) InitializeSavedTeachers() {}

// FindTeachers is to do nothing
func (d *DummySite) FindTeachers() []models.TeacherInfo {
	return nil
}
