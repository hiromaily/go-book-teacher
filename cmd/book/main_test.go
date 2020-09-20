package main

import (
	"os"
	"testing"

	"github.com/hiromaily/go-book-teacher/pkg/config"
	enc "github.com/hiromaily/golibs/cipher/encryption"
	lg "github.com/hiromaily/golibs/log"
)

var bookerTests = []struct {
	storage      int    // 1:text, 2:redis, anything else:dummy
	notification int    // 1:slack, 2:browser, 3:mail, anything else: console
	jsonPath     string // if blank, default value should be used
	explanation  string
	err          error
}{
	{1, 99, "", "text storage should work", nil},
	{2, 99, "", "redis storage should work", nil},
	{1, 99, "../../testdata/json/teachers.json", "initial data is created from json file", nil},
}

func setup() {
	lg.InitializeLog(lg.InfoStatus, lg.TimeShortFile, "[GO-BOOK-TEACHER_TEST]", "", "hiromaily")

	// Note: redis should be run in advance
	_, err := enc.NewCryptWithEnv()
	if err != nil {
		panic(err)
	}
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func TestIntegrationBooker(t *testing.T) {
	for _, tt := range bookerTests {
		// create config
		conf := createConfig(tt.storage, tt.notification)
		regi := NewRegistry(conf)

		// run
		booker := regi.NewBooker(tt.jsonPath, 0)
		if err := booker.Start(); err != nil {
			t.Errorf("fail: %s", tt.explanation)
		}
		booker.Cleanup()
	}
}

func createConfig(storage, notification int) *config.Config {
	conf := config.Config{}
	// site
	conf.Site = &config.SiteConfig{
		Type:        "DMM",
		URL:         "http://eikaiwa.dmm.com/",
		Concurrency: 10,
	}

	crypt := enc.GetCrypt()
	switch storage {
	case 1:
		// text
		conf.Text = &config.TextConfig{Path: "test.log"}
	case 2:
		// redis
		conf.Redis = &config.RedisConfig{URL: "redis://h:password@127.0.0.1:6379"}
	default:
		// dummy
	}

	switch notification {
	case 1:
		// slack
		key := "HP/9upIf+CwLGuDj2V0xfqulICwv1nHhQXy+S2TSEhFzYfEnt9zzWjVtoMT/8Rb7"
		conf.Slack.Key, _ = crypt.DecryptBase64(key)
	case 2:
		// browser
		conf.Browser.Enabled = true
	// case 3:
	//	//mail
	//	conf.Mail.Encrypted = true
	//	conf.Mail.MailFrom = "your@mail.com"
	default:
		// console
	}
	return &conf
}
