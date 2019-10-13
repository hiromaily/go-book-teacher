package siter

import (
	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/models"
	"github.com/hiromaily/go-book-teacher/pkg/siter/dmmer"
	"github.com/hiromaily/go-book-teacher/pkg/siter/dummy"
)

type Siter interface {
	FetchInitialData() error
	InitializeSavedTeachers()
	FindTeachers() []models.TeacherInfo
}

func NewSiter(jsonPath string, siteConf *config.SiteConfig) Siter {
	//for now, only DMM is available
	switch siteConf.Type {
	case SiteTypeDMM.String():
		return dmmer.NewDMM(jsonPath, siteConf.URL, siteConf.Concurrency)
	default:
		return dummy.NewDummySite()
	}
}

type SiteType string

const (
	SiteTypeDMM   SiteType = "DMM"
	SiteTypeDummy SiteType = "Dummy"
)

func (s SiteType) String() string {
	return string(s)
}
