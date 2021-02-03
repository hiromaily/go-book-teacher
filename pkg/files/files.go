package files

import (
	"fmt"
	"os"

	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/teachers"
)

// GetConfigPath returns toml file path
func GetConfigPath(tomlPath string) string {
	if tomlPath != "" && isExist(tomlPath) {
		return tomlPath
	}
	// book-teacher.toml
	expectedFileName := fmt.Sprintf("/usr/local/bin/%s.toml", os.Args[0])
	if isExist(expectedFileName) {
		return expectedFileName
	}
	envFile := config.GetEnvConfPath()
	if envFile != "" && isExist(envFile) {
		return envFile
	}
	return ""
}

// GetJSONPath returns JSON file path
func GetJSONPath(jsonPath string) string {
	if jsonPath != "" && isExist(jsonPath) {
		return jsonPath
	}
	// book-teacher.json
	expectedFileName := fmt.Sprintf("/usr/local/bin/%s.json", os.Args[0])
	if isExist(expectedFileName) {
		return expectedFileName
	}
	envFile := teachers.GetEnvJSONPath()
	if envFile != "" && isExist(envFile) {
		return envFile
	}
	return ""
}

func isExist(file string) bool {
	_, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return false // file is not existing
		}
		return false // error occurred somehow, e.g. permission error
	}
	return true
}
