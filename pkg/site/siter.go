package site

import (
	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/models"
	"github.com/hiromaily/go-book-teacher/pkg/site/dmmer"
	"github.com/hiromaily/go-book-teacher/pkg/site/dummy"
)

// Siter is Siter interface
type Siter interface {
	FetchInitialData() error
	// InitializeSavedTeachers()
	FindTeachers(day int) []models.TeacherInfo
}

// NewSiter is to return Siter interface
func NewSiter(jsonPath string, siteConf *config.Site) Siter {
	switch siteConf.Type {
	case SiteTypeDMM.String():
		return dmmer.NewDMM(jsonPath, siteConf.URL)
	default:
		return dummy.NewDummySite()
	}
}

// SiteType is SiteType
type SiteType string

const (
	// SiteTypeDMM is DMM
	SiteTypeDMM SiteType = "DMM"
	// SiteTypeDummy is dummy
	SiteTypeDummy SiteType = "Dummy"
)

// String is to convert SiteType to string
func (s SiteType) String() string {
	return string(s)
}
