package files

import (
	"fmt"
	"os"

	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/teachers"
)

const installPath = "/usr/local/etc/book-teacher/"

// GetConfigPath returns toml file path
func GetConfigPath(tomlPath string) string {
	if tomlPath != "" && isExist(tomlPath) {
		return tomlPath
	}
	// book-teacher.toml
	if installed := getInstallPath("toml"); isExist(installed) {
		return installed
	}
	//targetDir, err := getBinDir()
	//if err == nil {
	//	expectedFileName := fmt.Sprintf("%s%s.toml", targetDir, os.Args[0])
	//	if isExist(expectedFileName) {
	//		return expectedFileName
	//	}
	//}

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
	if installed := getInstallPath("json"); isExist(installed) {
		return installed
	}
	//targetDir, err := getBinDir()
	//if err == nil {
	//	expectedFileName := fmt.Sprintf("%s%s.json", targetDir, os.Args[0])
	//	if isExist(expectedFileName) {
	//		return expectedFileName
	//	}
	//}

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

func getInstallPath(ext string) string {
	return fmt.Sprintf("%s%s.%s", installPath, os.Args[0], ext)
}

//func getBinDir() (string, error) {
//	out, err := exec.Command("which", []string{os.Args[0]}...).Output()
//	if err != nil {
//		return "", err
//	}
//	// FIXME: windows newline is \r\n
//	return strings.TrimRight(string(out), fmt.Sprintf("%s\n", os.Args[0])), nil
//}
