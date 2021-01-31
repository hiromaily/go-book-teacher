package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-book-teacher/pkg/encryption"
	"github.com/hiromaily/go-book-teacher/pkg/storage"
)

// NewConfig returns *Root config
func NewConfig(fileName string) (*Root, error) {
	conf, err := loadConfig(fileName)
	if err != nil {
		return nil, err
	}

	if err = conf.decrypt(); err != nil {
		return nil, err
	}
	return conf, err
}

// GetEnvConfPath returns toml file path from environment variable `$GO-BOOK_CONF`
func GetEnvConfPath() string {
	path := os.Getenv("GO_BOOK_CONF")
	if strings.Contains(path, "${GOPATH}") {
		gopath := os.Getenv("GOPATH")
		path = strings.Replace(path, "${GOPATH}", gopath, 1)
	}
	return path
}

// load config file
func loadConfig(fileName string) (*Root, error) {
	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read file: %s", fileName)
	}

	var root Root
	_, err = toml.Decode(string(d), &root)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to parse: %s", fileName)
	}

	// check validation of config
	if err = root.validate(); err != nil {
		return nil, err
	}

	return &root, nil
}

func (r *Root) validate() error {
	validate := validator.New()

	excepted := make([]string, 0)
	if r.Storage.Mode == storage.TextMode {
		excepted = append(excepted, []string{"Storage.Redis", "Storage.Redis.URL"}...)
	} else {
		excepted = append(excepted, "Text")
	}
	if !r.Notification.Console.Enabled {
		excepted = append(excepted, "Console")
	}
	if !r.Notification.Slack.Enabled {
		excepted = append(excepted, []string{"Notification.Slack", "Notification.Key"}...)
	}

	return validate.StructExcept(r, excepted...)
}

// decrypt decrypts encrypted values in config file
func (r *Root) decrypt() error {
	crypt, err := encryption.NewCryptWithEnv()
	if err != nil {
		return err
	}

	if r.Storage.Redis.Encrypted {
		target := r.Storage.Redis
		target.URL, _ = crypt.DecryptBase64(target.URL)
	}

	if r.Notification.Slack.Encrypted {
		target := r.Notification.Slack
		target.Key, _ = crypt.DecryptBase64(target.Key)
	}
	return nil
}
