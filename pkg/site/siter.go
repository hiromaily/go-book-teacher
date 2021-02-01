package site

import (
	"github.com/hiromaily/go-book-teacher/pkg/models"
)

// Siter is Siter interface
type Siter interface {
	FetchInitialData() error
	// InitializeSavedTeachers()
	FindTeachers(day int) []models.TeacherInfo
}

// SiteType is site type
type SiteType string

const (
	// SiteTypeDMM is DMM site
	SiteTypeDMM SiteType = "dmm"
)

// String converts SiteType to string
func (s SiteType) String() string {
	return string(s)
}
