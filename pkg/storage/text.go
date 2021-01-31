package storage

import (
	"bufio"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
)

// TextRepo text repository
type textRepo struct {
	mode     Mode
	logger   *zap.Logger
	filePath string
}

// NewText is to return TextRepo object
func NewText(logger *zap.Logger, path string) Storager {
	return &textRepo{
		mode:     TextMode,
		logger:   logger,
		filePath: path,
	}
}

// Save is to save data to text
func (t *textRepo) Save(newData string) (bool, error) {
	t.logger.Debug("save", zap.String("mode", t.mode.String()))

	// open saved log
	fp, err := os.OpenFile(t.filePath, os.O_CREATE, 0o664)
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

	// save latest info
	content := []byte(newData)
	ioutil.WriteFile(t.filePath, content, 0o664)
	return true, nil
}

// Delete deletes file
func (t *textRepo) Delete() error {
	return os.Remove(t.filePath)
}

// Close closes nothing
func (t *textRepo) Close() {}
