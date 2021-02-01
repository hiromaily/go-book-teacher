package site

import (
	"github.com/hiromaily/go-book-teacher/pkg/teachers"
)

// Siter interface
type Siter interface {
	Fetch() error
	FindTeachers() []teachers.TeacherRepo
}

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------

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
