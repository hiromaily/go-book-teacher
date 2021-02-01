package teachers

import (
	"go.uber.org/zap"
)

type dummyTeacher struct {
	logger *zap.Logger
}

// NewDummyTeacher returns Teacher
func NewDummyTeacher(logger *zap.Logger) Teacher {
	return &dummyTeacher{
		logger: logger,
	}
}

// Fetch fetches dummy []TeacherRepo
func (d *dummyTeacher) Fetch() ([]TeacherRepo, error) {
	return []TeacherRepo{
		{ID: 6214, Name: "Aleksandra S", Country: "Serbia"},
		{ID: 4808, Name: "Joxyly", Country: "Serbia"},
		{ID: 23979, Name: "Lina Bianca", Country: "USA"},
		{ID: 25070, Name: "Celene", Country: "Australia"},
		{ID: 24721, Name: "Kenzie", Country: "USA"},
		{ID: 27828, Name: "Sanndy", Country: "UK"},
		{ID: 28302, Name: "Danni", Country: "South Africa"},
		{ID: 30216, Name: "Tamm", Country: "UK"},
		{ID: 25302, Name: "Nami", Country: "USA"},
		{ID: 32141, Name: "Colleen Marie", Country: "USA"},
	}, nil

}
