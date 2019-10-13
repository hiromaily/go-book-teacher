package main

//
//import (
//	"flag"
//	"fmt"
//	"os"
//	"testing"
//
//	conf "github.com/hiromaily/go-book-teacher/pkg/config"
//	th "github.com/hiromaily/go-book-teacher/pkg/siter"
//	"github.com/hiromaily/go-book-teacher/pkg/storages"
//	lg "github.com/hiromaily/golibs/log"
//	r "github.com/hiromaily/golibs/runtimes"
//)
//
//var txtPath = "./status.log"
//var jsonPath = "../../testdata/json/teachers.json"
//
////-----------------------------------------------------------------------------
//// Test Framework
////-----------------------------------------------------------------------------
//// Initialize
//func init() {
//	flag.Parse()
//
//	//lg.InitializeLog(lg.DebugStatus, lg.LogOff, 99, "[GO-BOOK-TEACHER_TEST]", "/var/log/go/test.log")
//	lg.InitializeLog(lg.InfoStatus, lg.TimeShortFile, "[GO-BOOK-TEACHER_TEST]", "", "hiromaily")
//}
//
//func setup() {
//	//check parameter
//	checkParam()
//
//	//
//	setupMain()
//
//	//print off
//	th.SetPrintOn(false)
//}
//
//func teardown() {
//	if storages.GetRedis() != nil {
//		storages.GetRedis().RD.Close()
//	}
//}
//
//func TestMain(m *testing.M) {
//	setup()
//
//	code := m.Run()
//
//	teardown()
//
//	os.Exit(code)
//}
//
////-----------------------------------------------------------------------------
//// functions
////-----------------------------------------------------------------------------
//func checkParam() {
//	lg.Debugf("*tomlPath: %s", *tomlPath)
//	if *tomlPath == "" {
//		*tomlPath = os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-book-teacher/data/toml/mailon.toml"
//	}
//	conf.New(*tomlPath, true)
//
//	//In this case, skip mail test
//	//m := conf.GetConf().Mail
//	//if m.MailTo == "" || m.MailFrom == "" || m.SMTP.Pass == "" ||
//	//	m.SMTP.Server == "" || m.SMTP.Port == 0 || conf.GetConf().Redis.URL == "" {
//	//	panic("parameter is wrong.")
//	//}
//}
//
//func clear() {
//	if redisFlg {
//		clearData(storages.GetRedis())
//	} else {
//		clearData(storages.GetText())
//	}
//	//deleteTxt(txtPath)
//	//deleteRedisKey()
//
//	*jsPath = ""
//}
//
////-----------------------------------------------------------------------------
//// Test
////-----------------------------------------------------------------------------
////-----------------------------------------------------------------------------
//// Main
//// TODO:Execute all pattern
//// 1. on local and txt file and browser
//// 2. on local and redis and mail
//// 3. on local and txt file and load json and browser
//// 4. on heroku and txt file and mail
//// 5. on heroku and redis and browser
//// 6. on heroku and redis and mail
//// 7. on heroku and redis and browser
//// 8. when keeping saved file, how it work well.
////-----------------------------------------------------------------------------
//// 1. on local and txt file and browser
//func TestIntegrationOnLocalUsingTxtAndBrowser(t *testing.T) {
//	//t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))
//
//	//settings
//	herokuFlg = "0"
//	redisFlg = false
//	mailFlg = false
//	slackFlg = false
//
//	//#1
//	bRet := ExecMain(1)
//	if !bRet {
//		t.Error("failed something.")
//	}
//
//	//#2. using saved data.
//	bRet = ExecMain(1)
//	if !bRet {
//		t.Error("failed something.")
//	}
//
//	clear()
//}
//
//// 2. on local and redis and mail
//func TestIntegrationOnLocalUsingRedisAndMail(t *testing.T) {
//	if conf.GetConf().Mail.MailTo == "" {
//		t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))
//	}
//
//	//settings
//	herokuFlg = "0"
//	redisFlg = true
//	mailFlg = true
//	slackFlg = false
//
//	//#1
//	bRet := ExecMain(1)
//	if !bRet {
//		t.Error("failed something.")
//	}
//
//	//#2. using saved data.
//	bRet = ExecMain(1)
//	if !bRet {
//		t.Error("failed something.")
//	}
//
//	clear()
//}
//
//// 3. on local and redis and mail
//func TestIntegrationOnLocalUsingRedisAndSlack(t *testing.T) {
//	if conf.GetConf().Slack.Key == "" || conf.GetConf().Redis.URL == "" {
//		t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))
//	}
//
//	//settings
//	herokuFlg = "0"
//	redisFlg = true
//	mailFlg = false
//	slackFlg = true
//
//	//#1
//	bRet := ExecMain(1)
//	if !bRet {
//		t.Error("failed something.")
//	}
//
//	//#2. using saved data.
//	bRet = ExecMain(1)
//	if !bRet {
//		t.Error("failed something.")
//	}
//
//	clear()
//}
//
//// 4. on local and txt file and load json and browser
//func TestIntegrationOnLocalUsingTxtAndBrowserAndJson(t *testing.T) {
//	//t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))
//
//	//settings
//	herokuFlg = "0"
//	redisFlg = false
//	mailFlg = false
//	slackFlg = false
//
//	*jsPath = jsonPath
//
//	//#1
//	bRet := ExecMain(1)
//	if !bRet {
//		t.Error("failed something.")
//	}
//
//	clear()
//}
//
//// 5. on heroku and txt file and mail
//// It supposes not to work intentionally.
//func TestIntegrationOnHerokuUsingTxtAndMail(t *testing.T) {
//	if conf.GetConf().Mail.MailTo == "" {
//		t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))
//	}
//
//	//settings
//	herokuFlg = "1"
//	redisFlg = false
//	mailFlg = true
//	slackFlg = false
//
//	//#1
//	err := checkHeroku()
//	if err == nil {
//		t.Error("failed something.")
//	}
//
//	clear()
//}
//
//// 6. on heroku and txt file and mail
//// It supposes not to work intentionally.
//func TestIntegrationOnHerokuUsingTxtAndSlack(t *testing.T) {
//	if conf.GetConf().Slack.Key == "" {
//		t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))
//	}
//
//	//settings
//	herokuFlg = "1"
//	redisFlg = false
//	mailFlg = false
//	slackFlg = true
//
//	//#1
//	err := checkHeroku()
//	if err == nil {
//		t.Error("failed something.")
//	}
//
//	clear()
//}
//
//// 7. on heroku and redis and browser
//// It supposes not to work intentionally.
//func TestIntegrationOnHerokuUsingRedisAndBrowser(t *testing.T) {
//	if conf.GetConf().Redis.URL == "" {
//		t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))
//	}
//	//TODO:still failed
//
//	//settings
//	herokuFlg = "1"
//	redisFlg = true
//	mailFlg = false
//	slackFlg = false
//
//	//#1
//	err := checkHeroku()
//	if err == nil {
//		t.Error("failed something.")
//	}
//
//	clear()
//}
//
//// 8. on heroku and redis and mail
//func TestIntegrationOnHerokuUsingRedisAndMail(t *testing.T) {
//	if conf.GetConf().Redis.URL == "" || conf.GetConf().Mail.MailTo == "" {
//		t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))
//	}
//
//	//settings
//	herokuFlg = "1"
//	redisFlg = true
//	mailFlg = true
//	slackFlg = false
//
//	//#1
//	bRet := ExecMain(1)
//	if !bRet {
//		t.Error("failed something.")
//	}
//
//	clear()
//}
//
//// 9. on heroku and redis and mail
//func TestIntegrationOnHerokuUsingRedisAndSlack(t *testing.T) {
//	if conf.GetConf().Redis.URL == "" || conf.GetConf().Slack.Key == "" {
//		t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))
//	}
//
//	//settings
//	herokuFlg = "1"
//	redisFlg = true
//	mailFlg = false
//	slackFlg = true
//
//	//#1
//	bRet := ExecMain(1)
//	if !bRet {
//		t.Error("failed something.")
//	}
//
//	clear()
//}
