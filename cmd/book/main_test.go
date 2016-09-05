package main

import (
	"flag"
	lg "github.com/hiromaily/golibs/log"
	//r "github.com/hiromaily/golibs/runtimes"
	conf "github.com/hiromaily/go-book-teacher/config"
	th "github.com/hiromaily/go-book-teacher/teacher"
	"os"
	"testing"
)

var txtPath string = "./status.log"
var jsonPath string = "../../json/teachers.json"

//-----------------------------------------------------------------------------
// Test Framework
//-----------------------------------------------------------------------------
// Initialize
func init() {
	flag.Parse()

	//lg.InitializeLog(lg.DEBUG_STATUS, lg.LOG_OFF_COUNT, 0, "[GO-BOOK-TEACHER_TEST]", "/var/log/go/test.log")
	lg.InitializeLog(lg.INFO_STATUS, lg.LOG_OFF_COUNT, 0, "[GO-BOOK-TEACHER_TEST]", "/var/log/go/test.log")
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

	if m.MailTo == "" || m.MailFrom == "" || m.Smtp.Pass == "" ||
		m.Smtp.Server == "" || m.Smtp.Port == 0 || conf.GetConf().Redis.URL == "" {
		panic("paramter is wrong.")
	}
}

func clearData() {
	deleteTxt(txtPath)
	deleteRedisKey()

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

	clearData()
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

	clearData()
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

	clearData()
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

	clearData()
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

	clearData()
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

	clearData()
}
