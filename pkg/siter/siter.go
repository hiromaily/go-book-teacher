package siter

import (
	"github.com/hiromaily/go-book-teacher/pkg/models"
	"github.com/hiromaily/go-book-teacher/pkg/siter/dmmer"
)

type Siter interface {
	FetchInitialData() error
	InitializeSavedTeachers()
	HandleTeachers()
	GetSavedTeachers() []models.TeacherInfo
}

func NewSiter(jsonPath string) Siter {
	//for now, only DMM is available
	return dmmer.NewDMM(jsonPath)
}
