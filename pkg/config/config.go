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
	// Concurrency int `toml:"concurrency"`
	Site    *SiteConfig
	Redis   *RedisConfig
	Text    *TextConfig
	Slack   *SlackConfig
	Browser *BrowserConfig
	Mail    *MailConfig
}

// SiteConfig is for site information
type SiteConfig struct {
	Type        string `toml:"type"`
	URL         string `toml:"url"`
	Concurrency int    `toml:"concurrency"`
}

// RedisConfig is storage by redis server
type RedisConfig struct {
	Encrypted bool   `toml:"encrypted"`
	URL       string `toml:"url"`
}

// TextConfig is storage by text
type TextConfig struct {
	Path string `toml:"path"`
}

// SlackConfig is notification by slack
type SlackConfig struct {
	Encrypted bool   `toml:"encrypted"`
	Key       string `toml:"key"`
}

// BrowserConfig is notification by browser
type BrowserConfig struct {
	Enabled bool `toml:"enabled"`
}

// MailConfig is notification by mail
type MailConfig struct {
	Encrypted bool        `toml:"encrypted"`
	MailTo    string      `toml:"mail_to"`
	MailFrom  string      `toml:"mail_from"`
	SMTP      *SMTPConfig `toml:"smtp"`
}

// SMTPConfig is configuraion of smtp server
type SMTPConfig struct {
	Address string `toml:"address"`
	Pass    string `toml:"pass"`
	Server  string `toml:"server"`
	Port    int    `toml:"port"`
}

var checkTomlKeys = [][]string{
	//{"concurrency"},
	{"site", "type"},
	{"site", "url"},
	{"site", "concurrency"},
	{"redis", "encrypted"},
	{"redis", "url"},
	{"text", "path"},
	{"slack", "encrypted"},
	{"slack", "key"},
	{"browser", "enabled"},
	{"mail", "encrypted"},
	{"mail", "mail_to"},
	{"mail", "mail_from"},
	{"mail", "smtp", "address"},
	{"mail", "smtp", "pass"},
	{"mail", "smtp", "server"},
	{"mail", "smtp", "port"},
}

func init() {
	tomlFileName = os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-book-teacher/config/toml/text-command.toml"
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

	// check validation of config
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

// validateConfig is to validate config settings
func (c *Config) validateConfig(md *toml.MetaData) error {
	// for protection when debugging on non production environment
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
			// case 4:
			//	format = "[%s.%s.%s] %s"
			default:
				// invalid check string
				return errors.New("error validation should be checked")
			}
			keysIfc := u.SliceStrToInterface(keys)
			errStrings = append(errStrings, fmt.Sprintf(format, keysIfc...))
		}
	}

	// error
	if len(errStrings) != 0 {
		return errors.Errorf("There are lacks of keys : %#v", errStrings)
	}
	return nil
}

// Cipher is to decrypt encrypted value of toml file
func (c *Config) Cipher() {
	crypt := enc.GetCrypt()

	if c.Redis.Encrypted {
		cf := c.Redis
		cf.URL, _ = crypt.DecryptBase64(cf.URL)
	}

	if c.Slack.Encrypted {
		cf := c.Slack
		cf.Key, _ = crypt.DecryptBase64(cf.Key)
	}

	if c.Mail.Encrypted {
		cf := c.Mail
		cf.MailTo, _ = crypt.DecryptBase64(cf.MailTo)
		cf.MailFrom, _ = crypt.DecryptBase64(cf.MailFrom)

		cf2 := c.Mail.SMTP
		cf2.Address, _ = crypt.DecryptBase64(cf2.Address)
		cf2.Pass, _ = crypt.DecryptBase64(cf2.Pass)
		cf2.Server, _ = crypt.DecryptBase64(cf2.Server)
	}
}

// ValidateSlack to validate slack is enabled or not
func (c *Config) ValidateSlack() bool {
	if c.Slack == nil || c.Slack.Key == "" {
		return false
	}
	return true
}

// ValidateBrowser to validate browser is enabled or not
func (c *Config) ValidateBrowser() bool {
	if c.Browser == nil || !c.Browser.Enabled {
		return false
	}
	return true
}

// ValidateMail to validate mail is enabled or not
func (c *Config) ValidateMail() bool {
	if c.Mail == nil {
		return false
	}
	if c.Mail.MailFrom == "" || c.Mail.MailTo == "" {
		return false
	}
	if c.Mail.SMTP == nil {
		return false
	}
	smtp := c.Mail.SMTP
	if smtp.Address == "" || smtp.Port == 0 || smtp.Server == "" {
		return false
	}
	return true
}

// ValidateRedis to validate redis is enabled or not
func (c *Config) ValidateRedis() bool {
	if c.Redis == nil || c.Redis.URL == "" {
		return false
	}
	return true
}

// ValidateText to validate text is enabled or not
func (c *Config) ValidateText() bool {
	if c.Text == nil || c.Text.Path == "" {
		return false
	}
	return true
}
