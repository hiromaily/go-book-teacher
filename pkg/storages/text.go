package storages

import (
	"bufio"
	"io/ioutil"
	"os"

	conf "github.com/hiromaily/go-book-teacher/pkg/config"
	lg "github.com/hiromaily/golibs/log"
)

// Text object
type TextRepo struct {
	filePath string
}

var txt = TextRepo{
	filePath: "/tmp/status.log",
}

// Setup is settings
func SetupText() {
	if conf.GetConf().StatusFile != "" {
		txt.filePath = conf.GetConf().StatusFile
	}
}

// Get is to get Text instance
func GetText() *TextRepo {
	return &txt
}

// Save is to save data to text
func (t TextRepo) Save(newData string) bool {
	lg.Debug("Using TxtFile")

	//open saved log
	fp, err := os.OpenFile(t.filePath, os.O_CREATE, 0664)
	defer fp.Close()

	if err == nil {

		scanner := bufio.NewScanner(fp)
		scanner.Scan()
		filedata := scanner.Text()

		if newData == filedata {
			return false
		}
	} else {
		panic(err.Error())
	}

	//save latest info
	content := []byte(newData)
	ioutil.WriteFile(t.filePath, content, 0664)
	return true
}

// Delete is to delete file
func (t TextRepo) Delete() error {
	//func (t Text) Delete(txtPath string) {
	//os.Remove(txtPath)
	return os.Remove(t.filePath)
}
