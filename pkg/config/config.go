package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	enc "github.com/hiromaily/golibs/cipher/encryption"
	u "github.com/hiromaily/golibs/utils"
)

var (
	tomlFileName = "./data/toml/settings.toml"
	conf         *Config
)

// Config is root of toml config
type Config struct {
	//Environment int `toml:"environment"`
	Redis *RedisConfig
	Text  *TextConfig
	Slack *SlackConfig
	Mail  *MailConfig
}

// RedisConfig is for redis server
type RedisConfig struct {
	Encrypted bool   `toml:"encrypted"`
	URL       string `toml:"url"`
}

type TextConfig struct {
	Path string `toml:"path"`
}

// SlackConfig is for slack
type SlackConfig struct {
	Encrypted bool   `toml:"encrypted"`
	Key       string `toml:"key"`
}

// MailConfig is for mail
type MailConfig struct {
	Encrypted bool        `toml:"encrypted"`
	MailTo    string      `toml:"mail_to"`
	MailFrom  string      `toml:"mail_from"`
	SMTP      *SMTPConfig `toml:"smtp"`
}

// SMTPConfig is for smtp server of mail
type SMTPConfig struct {
	Address string `toml:"address"`
	Pass    string `toml:"pass"`
	Server  string `toml:"server"`
	Port    int    `toml:"port"`
}

var checkTomlKeys = [][]string{
	//{"environment"},
	{"redis", "encrypted"},
	{"redis", "url"},
	{"text", "path"},
	{"slack", "encrypted"},
	{"slack", "key"},
	{"mail", "encrypted"},
	{"mail", "mail_to"},
	{"mail", "mail_from"},
	{"mail", "smtp", "address"},
	{"mail", "smtp", "pass"},
	{"mail", "smtp", "server"},
	{"mail", "smtp", "port"},
}

func init() {
	tomlFileName = os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-book-teacher/data/toml/settings.toml"
}

// load configfile
func loadConfig(path string) (*Config, error) {
	if path != "" {
		tomlFileName = path
	}

	d, err := ioutil.ReadFile(tomlFileName)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read %s", tomlFileName)
	}

	var config Config
	md, err := toml.Decode(string(d), &config)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to parse %s: %v", tomlFileName, md)
	}

	//check validation of config
	err = config.validateConfig(&md)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// New is for creating config instance
func New(file string, cipherFlg bool) error {
	var err error
	conf, err = loadConfig(file)
	if err != nil {
		return err
	}

	if cipherFlg {
		conf.Cipher()
	}
	return nil
}

// GetConf is to get config instance by singleton architecture
func GetConf() *Config {
	var err error
	if conf == nil {
		conf, err = loadConfig("")
	}
	if err != nil {
		panic(err)
	}

	return conf
}

//check validation of config
func (c *Config) validateConfig(md *toml.MetaData) error {
	//for protection when debugging on non production environment
	var errStrings []string
	var format string

	for _, keys := range checkTomlKeys {
		if !md.IsDefined(keys...) {
			switch len(keys) {
			case 1:
				format = "%s"
			case 2:
				format = "[%s] %s"
			case 3:
				format = "[%s.%s] %s"
			//case 4:
			//	format = "[%s.%s.%s] %s"
			default:
				//invalid check string
				return errors.New("error validation should be checked")
			}
			keysIfc := u.SliceStrToInterface(keys)
			errStrings = append(errStrings, fmt.Sprintf(format, keysIfc...))
		}
	}

	// Error
	if len(errStrings) != 0 {
		return errors.Errorf("There are lacks of keys : %#v", errStrings)
	}
	return nil
}

// Cipher is to decrypt encrypted value of toml file
func (c *Config) Cipher() {
	crypt := enc.GetCrypt()

	if c.Redis.Encrypted {
		cf := conf.Redis
		cf.URL, _ = crypt.DecryptBase64(cf.URL)
	}

	if c.Slack.Encrypted {
		cf := conf.Slack
		cf.Key, _ = crypt.DecryptBase64(cf.Key)
	}

	if c.Mail.Encrypted {
		cf := conf.Mail
		cf.MailTo, _ = crypt.DecryptBase64(cf.MailTo)
		cf.MailFrom, _ = crypt.DecryptBase64(cf.MailFrom)

		cf2 := conf.Mail.SMTP
		cf2.Address, _ = crypt.DecryptBase64(cf2.Address)
		cf2.Pass, _ = crypt.DecryptBase64(cf2.Pass)
		cf2.Server, _ = crypt.DecryptBase64(cf2.Server)
	}
}

func (c *Config) ValidateSlack() bool {
	if conf.Slack == nil || conf.Slack.Key == "" {
		return false
	}
	return true
}

func (c *Config) ValidateMail() bool {
	if conf.Mail == nil {
		return false
	}
	if conf.Mail.MailFrom == "" || conf.Mail.MailTo == "" {
		return false
	}
	if conf.Mail.SMTP == nil {
		return false
	}
	smtp := conf.Mail.SMTP
	if smtp.Address == "" || smtp.Port == 0 || smtp.Server == "" {
		return false
	}
	return true
}

func (c *Config) ValidateRedis() bool {
	if conf.Redis == nil || conf.Redis.URL == "" {
		return false
	}
	return true
}

func (c *Config) ValidateText() bool {
	if conf.Text == nil || conf.Text.Path == "" {
		return false
	}
	return true
}
