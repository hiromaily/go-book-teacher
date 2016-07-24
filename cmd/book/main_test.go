package main

import (
	"flag"
	//"fmt"
	lg "github.com/hiromaily/golibs/log"
	//r "github.com/hiromaily/golibs/runtimes"
	u "github.com/hiromaily/golibs/utils"
	"os"
	"testing"
)

var (
	benchFlg    = flag.Int("bc", 0, "Normal Test or Bench Test")
	mailToAdd   = flag.String("toadd", "", "MAIL_TO_ADDRESS")
	mailFromAdd = flag.String("fradd", "", "MAIL_FROM_ADDRESS")
	smtpPass    = flag.String("smpass", "", "SMTP_PASS")
	smtpServer  = flag.String("smsvr", "", "SMTP_SERVER")
	smtpPort    = flag.Int("smport", 0, "SMTP_PORT")
	redisURL    = flag.String("rdsvr", "", "REDIS_URL")
)

var txtPath string = "./status.log"
var jsonPath string = "../../settings.json"

func setup() {
	lg.InitializeLog(lg.DEBUG_STATUS, lg.LOG_OFF_COUNT, 0, "[GO-BOOK-TEACHER_TEST]", "/var/log/go/test.log")

	if *benchFlg == 0 {
	}

	//check parameter
	checkParam()

	//redis
	os.Setenv("REDIS_URL", *redisURL)
	RedisInit()
	//os.Clearenv()
	clearEnvOnThisTest()
}

func teardown() {
	if *benchFlg == 0 {
	}
}

// Initialize
func TestMain(m *testing.M) {
	flag.Parse()

	//TODO: According to argument, it switch to user or not.
	//TODO: For bench or not bench
	setup()

	code := m.Run()

	teardown()

	// 終了
	os.Exit(code)
}

func checkParam() {
	if *mailToAdd == "" || *mailFromAdd == "" || *smtpPass == "" ||
		*smtpServer == "" || *smtpPort == 0 || *redisURL == "" {
		panic("paramter is wrong.")
	}
}

func setupMail(t *testing.T) {
	if *mailToAdd == "" || *mailFromAdd == "" || *smtpPass == "" || *smtpServer == "" || *smtpPort == 0 {
		t.Fatal("parameter is wrong.")
	}

	//Add environmental variable from command-line parameter
	os.Setenv("MAIL_TO_ADDRESS", *mailToAdd)     //mail
	os.Setenv("MAIL_FROM_ADDRESS", *mailFromAdd) //mail
	os.Setenv("SMTP_ADDRESS", *mailFromAdd)      //mail
	os.Setenv("SMTP_PASS", *smtpPass)            //mail
	os.Setenv("SMTP_SERVER", *smtpServer)        //mail
	os.Setenv("SMTP_PORT", u.Itoa(*smtpPort))    //mail
}

func clearEnvOnThisTest() {
	os.Unsetenv("MAIL_TO_ADDRESS")
	os.Unsetenv("MAIL_FROM_ADDRESS")
	os.Unsetenv("SMTP_ADDRESS")
	os.Unsetenv("SMTP_PASS")
	os.Unsetenv("SMTP_SERVER")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("HEROKU_FLG")
	os.Unsetenv("SAVE_LOG")
	os.Unsetenv("REDIS_URL")
}

func clearData() {
	DeleteTxt(txtPath)
	DeleteRedisKey()

	*jsPath = ""

	//This is a little risk because all environment variables are removed.
	//os.Clearenv()
	clearEnvOnThisTest()
}

//-----------------------------------------------------------------------------
// Main
// TODO:全パターンをテストできるように
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

	//Add environmental variable from command-line parameter
	os.Setenv("HEROKU_FLG", "0")   //local
	os.Setenv("SAVE_LOG", txtPath) //txt log
	//os.Setenv("MAIL_TO_ADDRESS", "") //browser

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

	//Add environmental variable from command-line parameter
	setupMail(t)
	os.Setenv("HEROKU_FLG", "0")      //local
	os.Setenv("REDIS_URL", *redisURL) //Redis log

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

	//Add environmental variable from command-line parameter
	os.Setenv("HEROKU_FLG", "0")   //local
	os.Setenv("SAVE_LOG", txtPath) //txt log
	//os.Setenv("MAIL_TO_ADDRESS", "") //browser

	//json //set parameter dinamically
	//os.Args = append(os.Args, "-f")
	//os.Args = append(os.Args, "./settings.json")
	*jsPath = jsonPath //this path is track back from testfile

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

	//Add environmental variable from command-line parameter
	setupMail(t)
	os.Setenv("HEROKU_FLG", "1")   //local
	os.Setenv("SAVE_LOG", txtPath) //txt log

	bRet := ExecMain(1)
	if bRet {
		t.Error("failed something.")
	}

	clearData()
}

// 5. on heroku and redis and browser
// It supposes not to work intentionally.
func TestIntegrationOnHerokuUsingRedisAndBrowser(t *testing.T) {
	//t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))

	//Add environmental variable from command-line parameter
	os.Setenv("HEROKU_FLG", "1")      //local
	os.Setenv("REDIS_URL", *redisURL) //Redis log

	bRet := ExecMain(1)
	if bRet {
		t.Error("failed something.")
	}

	clearData()
}

// 6. on heroku and redis and mail
func TestIntegrationOnHerokuUsingRedisAndMail(t *testing.T) {
	//t.Skip(fmt.Sprintf("skipping %s", r.CurrentFunc(1)))

	//Add environmental variable from command-line parameter
	setupMail(t)
	os.Setenv("HEROKU_FLG", "1")      //local
	os.Setenv("REDIS_URL", *redisURL) //Redis log

	bRet := ExecMain(1)
	if !bRet {
		t.Error("failed something.")
	}

	clearData()
}
