package config

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	enc "github.com/hiromaily/golibs/cipher/encryption"
	u "github.com/hiromaily/golibs/utils"
	"io/ioutil"
	"os"
)

var (
	tomlFileName = "./config/settings.toml"
	conf         *Config
)

// Config is root of toml config
type Config struct {
	Environment int    `toml:"environment"`
	StatusFile  string `toml:"status_file"`
	Redis       *RedisConfig
	Mail        *MailConfig
}

// RedisConfig is for redis server
type RedisConfig struct {
	Encrypted bool   `toml:"encrypted"`
	URL       string `toml:"url"`
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
	{"environment"},
	{"redis", "encrypted"},
	{"redis", "url"},
	{"mail", "encrypted"},
	{"mail", "mail_to"},
	{"mail", "mail_from"},
	{"mail", "smtp", "address"},
	{"mail", "smtp", "pass"},
	{"mail", "smtp", "server"},
	{"mail", "smtp", "port"},
}

func init() {
	tomlFileName = os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-book-teacher/config/settings.toml"
}

//check validation of config
func validateConfig(conf *Config, md *toml.MetaData) error {
	//for protection when debugging on non production environment
	var errStrings []string

	format := "[%s]"
	inValid := false
	for _, keys := range checkTomlKeys {
		if !md.IsDefined(keys...) {
			switch len(keys) {
			case 1:
				format = "[%s]"
			case 2:
				format = "[%s] %s"
			case 3:
				format = "[%s.%s] %s"
			default:
				//invalid check string
				inValid = true
				break
			}
			keysIfc := u.SliceStrToInterface(keys)
			errStrings = append(errStrings, fmt.Sprintf(format, keysIfc...))
		}
	}

	// Error
	if inValid {
		return errors.New("Error: Check Text has wrong number of parameter")
	}
	if len(errStrings) != 0 {
		return fmt.Errorf("Error: There are lacks of keys : %#v \n", errStrings)
	}

	return nil
}

// load configfile
func loadConfig(path string) (*Config, error) {
	if path != "" {
		tomlFileName = path
	}

	d, err := ioutil.ReadFile(tomlFileName)
	if err != nil {
		return nil, fmt.Errorf(
			"Error reading %s: %s", tomlFileName, err)
	}

	var config Config
	md, err := toml.Decode(string(d), &config)
	if err != nil {
		return nil, fmt.Errorf(
			"Error parsing %s: %s(%v)", tomlFileName, err, md)
	}

	//check validation of config
	err = validateConfig(&config, &md)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// New is for creating config instance
func New(file string) {
	var err error
	conf, err = loadConfig(file)
	if err != nil {
		panic(err)
	}
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

// SetTomlPath is to set toml file path
func SetTomlPath(path string) {
	tomlFileName = path
}

// Cipher is to decrypt encrypted value of toml file
func Cipher() {
	crypt := enc.GetCrypt()

	if conf.Redis.Encrypted {
		c := conf.Redis
		c.URL, _ = crypt.DecryptBase64(c.URL)
	}

	if conf.Mail.Encrypted {
		c := conf.Mail
		c.MailTo, _ = crypt.DecryptBase64(c.MailTo)
		c.MailFrom, _ = crypt.DecryptBase64(c.MailFrom)

		c2 := conf.Mail.SMTP
		c2.Address, _ = crypt.DecryptBase64(c2.Address)
		c2.Pass, _ = crypt.DecryptBase64(c2.Pass)
		c2.Server, _ = crypt.DecryptBase64(c2.Server)
	}
}
