package main

import (
	"flag"
	conf "github.com/hiromaily/go-book-teacher/config"
	rd "github.com/hiromaily/go-book-teacher/redis"
	th "github.com/hiromaily/go-book-teacher/teacher"
	tt "github.com/hiromaily/go-book-teacher/text"
	lg "github.com/hiromaily/golibs/log"
	"os"
	"testing"
)

var txtPath = "./status.log"
var jsonPath = "../../json/teachers.json"

//-----------------------------------------------------------------------------
// Test Framework
//-----------------------------------------------------------------------------
// Initialize
func init() {
	flag.Parse()

	//lg.InitializeLog(lg.DebugStatus, lg.LogOff, 99, "[GO-BOOK-TEACHER_TEST]", "/var/log/go/test.log")
	lg.InitializeLog(lg.InfoStatus, lg.LogOff, 99, "[GO-BOOK-TEACHER_TEST]", "/var/log/go/test.log")
}

func setup() {
	//check parameter
	checkParam()

	//
	setupMain()

	//print off
	th.SetPrintOn(false)
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

//-----------------------------------------------------------------------------
// functions
//-----------------------------------------------------------------------------
func checkParam() {

	m := conf.GetConf().Mail

	if m.MailTo == "" || m.MailFrom == "" || m.SMTP.Pass == "" ||
		m.SMTP.Server == "" || m.SMTP.Port == 0 || conf.GetConf().Redis.URL == "" {
		panic("parameter is wrong.")
	}
}

func clear() {
	if redisFlg {
		clearData(rd.Get())
	} else {
		clearData(tt.Get())
	}
	//deleteTxt(txtPath)
	//deleteRedisKey()

	*jsPath = ""
}

//-----------------------------------------------------------------------------
// Test
//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
// Main
// TODO:Execute all pattern
// 1. on local and txt file and browser
// 2. on local and redis and mail
// 3. on local and txt file and load json and browser
// 4. on heroku and txt file and mail
// 5. on heroku and redis and browser
// 6. on heroku and redis and mail
// 7. on heroku and redis and browser
// 8. when keeping saved file, how it work well.
//-----------------------------------------------------------------------------
// 1. on local and txt file and browser
func TestIntegrationOnLocalUsingTxtAndBrowser(t *testing.T) {
	//t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))

	//settings
	herokuFlg = "0"
	redisFlg = false
	mailFlg = false

	//#1
	bRet := ExecMain(1)
	if !bRet {
		t.Error("failed something.")
	}

	//#2. using saved data.
	bRet = ExecMain(1)
	if !bRet {
		t.Error("failed something.")
	}

	clear()
}

// 2. on local and redis and mail
func TestIntegrationOnLocalUsingRedisAndMail(t *testing.T) {
	//t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))

	//settings
	herokuFlg = "0"
	redisFlg = true
	mailFlg = true

	//#1
	bRet := ExecMain(1)
	if !bRet {
		t.Error("failed something.")
	}

	//#2. using saved data.
	bRet = ExecMain(1)
	if !bRet {
		t.Error("failed something.")
	}

	clear()
}

// 3. on local and txt file and load json and browser
func TestIntegrationOnLocalUsingTxtAndBrowserAndJson(t *testing.T) {
	//t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))

	//settings
	herokuFlg = "0"
	redisFlg = false
	mailFlg = false

	*jsPath = jsonPath

	//#1
	bRet := ExecMain(1)
	if !bRet {
		t.Error("failed something.")
	}

	clear()
}

// 4. on heroku and txt file and mail
// It supposes not to work intentionally.
func TestIntegrationOnHerokuUsingTxtAndMail(t *testing.T) {
	//t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))
	//TODO:still failed

	//settings
	herokuFlg = "1"
	redisFlg = false
	mailFlg = true

	//#1
	err := checkHeroku()
	if err == nil {
		t.Error("failed something.")
	}

	clear()
}

// 5. on heroku and redis and browser
// It supposes not to work intentionally.
func TestIntegrationOnHerokuUsingRedisAndBrowser(t *testing.T) {
	//t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))
	//TODO:still failed

	//settings
	herokuFlg = "1"
	redisFlg = true
	mailFlg = false

	//#1
	err := checkHeroku()
	if err == nil {
		t.Error("failed something.")
	}

	clear()
}

// 6. on heroku and redis and mail
func TestIntegrationOnHerokuUsingRedisAndMail(t *testing.T) {
	//t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))

	//settings
	herokuFlg = "1"
	redisFlg = true
	mailFlg = true

	//#1
	bRet := ExecMain(1)
	if !bRet {
		t.Error("failed something.")
	}

	clear()
}
