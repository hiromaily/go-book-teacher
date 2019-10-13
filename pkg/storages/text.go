package storages

import (
	"bufio"
	"io/ioutil"
	"os"

	lg "github.com/hiromaily/golibs/log"
)

// TextRepo is TextRepo object
type TextRepo struct {
	mode     string
	filePath string
}

// NewText is to return TextRepo object
func NewText(path string) *TextRepo {
	return &TextRepo{
		mode:     "text",
		filePath: path,
	}
}

// Save is to save data to text
func (t *TextRepo) Save(newData string) (bool, error) {
	lg.Debugf("Save by %s", t.mode)

	//open saved log
	fp, err := os.OpenFile(t.filePath, os.O_CREATE, 0664)
	if err != nil {
		return false, err
	}
	defer fp.Close()

	// scan
	scanner := bufio.NewScanner(fp)
	scanner.Scan()
	filedata := scanner.Text()
	if newData == filedata {
		return false, nil
	}

	//save latest info
	content := []byte(newData)
	ioutil.WriteFile(t.filePath, content, 0664)
	return true, nil
}

// Delete is to delete file
func (t *TextRepo) Delete() error {
	//os.Remove(txtPath)
	return os.Remove(t.filePath)
}

// Close is to do nothing
func (t *TextRepo) Close() {}
