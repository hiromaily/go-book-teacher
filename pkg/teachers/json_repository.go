package teachers

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// TeacherRepository includes teacher's information
type TeacherRepository struct {
	Teachers []TeacherRepo `json:"teachers"`
}

type jsonTeacher struct {
	logger   *zap.Logger
	fileName string
}

// NewJSONTeacher returns Teacher
func NewJSONTeacher(logger *zap.Logger, fileName string) Teacher {
	return &jsonTeacher{
		logger:   logger,
		fileName: fileName,
	}
}

// Fetch fetches []TeacherRepo from JSON file
func (j *jsonTeacher) Fetch() ([]TeacherRepo, error) {
	repo := TeacherRepository{}
	file, err := ioutil.ReadFile(j.fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call ReadFile() %s", j.fileName)
	}
	err = json.Unmarshal(file, &repo)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to Unmarshal json binary: %s", j.fileName)
	}

	return repo.Teachers, nil
}
