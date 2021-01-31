package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-book-teacher/pkg/encryption"
	"github.com/hiromaily/go-book-teacher/pkg/storage"
)

// NewConfig returns *Root config
func NewConfig(fileName string, isEncrypted bool) (*Root, error) {
	conf, err := loadConfig(fileName)
	if err != nil {
		return nil, err
	}

	if isEncrypted {
		crypt, err := encryption.NewCryptWithEnv()
		if err != nil {
			return nil, err
		}
		conf.decrypt(crypt)
	}
	return conf, err
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
	// return validate.Struct(r)

	excepted := make([]string, 0)
	if r.Storage.Mode == storage.TextMode {
		excepted = append(excepted, "Redis")
	} else {
		excepted = append(excepted, "Text")
	}
	if !r.Notification.Console.Enabled {
		excepted = append(excepted, "CLI")
	}
	if !r.Notification.Browser.Enabled {
		excepted = append(excepted, "Browser")
	}
	if !r.Notification.Slack.Enabled {
		excepted = append(excepted, "Slack")
	}
	if !r.Notification.Mail.Enabled {
		excepted = append(excepted, []string{"Mail", "SMTP"}...)
	}

	return validate.StructExcept(r, excepted...)
}

// decrypt decrypts encrypted values in config file
func (r *Root) decrypt(crypt encryption.Crypt) {
	if r.Storage.Redis.Encrypted {
		target := r.Storage.Redis
		target.URL, _ = crypt.DecryptBase64(target.URL)
	}

	if r.Notification.Slack.Encrypted {
		target := r.Notification.Slack
		target.Key, _ = crypt.DecryptBase64(target.Key)
	}

	if r.Notification.Mail.Encrypted {
		target := r.Notification.Mail
		target.MailTo, _ = crypt.DecryptBase64(target.MailTo)
		target.MailFrom, _ = crypt.DecryptBase64(target.MailFrom)

		target.SMTP.Address, _ = crypt.DecryptBase64(target.SMTP.Address)
		target.SMTP.Pass, _ = crypt.DecryptBase64(target.SMTP.Pass)
		target.SMTP.Server, _ = crypt.DecryptBase64(target.SMTP.Server)
	}
}
