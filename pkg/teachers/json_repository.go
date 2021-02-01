package teachers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

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
	// j.logger.Debug("target_teacher", zap.Any("teachers", repo))
	return repo.Teachers, nil
}

// GetEnvJSONPath returns json path from environment variable
func GetEnvJSONPath() string {
	path := os.Getenv("GO_BOOK_JSON")
	if strings.Contains(path, "${GOPATH}") {
		gopath := os.Getenv("GOPATH")
		path = strings.Replace(path, "${GOPATH}", gopath, 1)
	}
	return path
}
