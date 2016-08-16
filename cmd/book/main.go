package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	th "github.com/hiromaily/go-book-teacher/teacher"
	rd "github.com/hiromaily/golibs/db/redis"
	hrk "github.com/hiromaily/golibs/heroku"
	lg "github.com/hiromaily/golibs/log"
	ml "github.com/hiromaily/golibs/mails"
	"github.com/hiromaily/golibs/signal"
	tm "github.com/hiromaily/golibs/times"
	"github.com/hiromaily/golibs/tmpl"
	u "github.com/hiromaily/golibs/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var (
	jsPath = flag.String("f", "", "Json file path")
)

var usage = `Usage: %s [options...]
Options:
  -f     Json file path
`

var mi *ml.MailInfo

var tmplMails string = `
The following tachers are available now!
{{range .Teachers}}
{{$.Url}}teacher/index/{{.Id}}/ [{{.Name}} / {{.Country}}]
{{end}}
Enjoy!`

const MaxGoRoutine uint16 = 20
const OpenFileName string = "/tmp/status.log"

var redisKey string = "bookteacher:save"

// setting for sending mail
func settingMail() {
	//get environment variable
	subject := "[ENGLISH LESSON] It's Available."
	body := ""
	//mails
	smtp := ml.Smtp{Address: os.Getenv("SMTP_ADDRESS"), Pass: os.Getenv("SMTP_PASS"),
		Server: os.Getenv("SMTP_SERVER"), Port: u.Atoi(os.Getenv("SMTP_PORT"))}

	//cannot use mails.MailInfo literal (type *mails.MailInfo) as type mails.MailInfo in assignment
	mi = &ml.MailInfo{ToAddress: []string{os.Getenv("MAIL_TO_ADDRESS")}, FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
		Subject: subject, Body: body, Smtp: smtp}
}

// send mail
func sendMail(ths []th.TeacherInfo) {
	//make body
	si := th.CreateSiteInfo(ths)
	body, err := tmpl.StrTempParser(tmplMails, &si)
	if err != nil {
		lg.Debugf("mail couldn't be send caused by err : %s\n", err)
	} else {
		mi.Body = body
		mi.SendMail("10s")
	}
}

//open browser on PC
func openBrowser(ths []th.TeacherInfo) {
	for _, t := range ths {
		//during test, it doesn't work.
		//open browser error: exec: "open": executable file not found in $PATH

		//out, err := exec.Command("open /Applications/Google\\ Chrome.app", fmt.Sprintf("http://eikaiwa.dmm.com/teacher/index/%d/", id)).Output()
		err := exec.Command("open", fmt.Sprintf("http://eikaiwa.dmm.com/teacher/index/%d/", t.Id)).Start()
		if err != nil {
			//panic(fmt.Sprintf("open browser error: %v", err))
			lg.Errorf(fmt.Sprintf("open browser error: %v", err))
		}
	}
}

// check Redis
func checkRedis(newData string) bool {
	lg.Debug("Using Redis")
	//redis
	//key := "bookteacher:save"

	//close
	//TODO:when use close
	defer rd.GetRedisInstance().Close()

	c := rd.GetRedisInstance().Conn
	val, err := redis.String(c.Do("GET", redisKey))
	lg.Debugf("redis error is %s\n", err)
	lg.Debugf("new value is %s, old value is %s\n", newData, val)
	if err != nil || newData != val {
		//save
		c.Do("SET", redisKey, newData)
		return true
	} else {
		return false
	}
}

func deleteRedisKey() {
	c := rd.GetRedisInstance().Conn
	_, err := c.Do("DEL", redisKey)
	if err != nil {
		lg.Debug("delete key on redis is failed.")
	}
}

// check txt file
func checkFile(newData string) bool {
	lg.Debug("Using TxtFile")
	//open saved log
	filePath := os.Getenv("SAVE_LOG")
	if filePath == "" {
		filePath = OpenFileName
	}

	fp, err := os.OpenFile(filePath, os.O_CREATE, 0664)
	defer fp.Close()

	if err == nil {

		scanner := bufio.NewScanner(fp)
		scanner.Scan()
		filedata := scanner.Text()

		if newData == filedata {
			return false
		}
	} else {
		panic(err.Error())
	}

	//save latest info
	content := []byte(newData)
	ioutil.WriteFile(filePath, content, 0664)
	return true
}

func deleteTxt(txtPath string) {
	os.Remove(txtPath)
	//err := os.Remove(txtPath)
	//if err != nil {
	//	lg.Errorf("to delete file was failed: %s, error is %s\n", txtPath, err)
	//}
}

//save teacher status to log
func saveStatusLog(ths []th.TeacherInfo) bool {

	//create string from ids slice
	var sum int = 0
	for _, t := range ths {
		sum += t.Id
	}
	newData := strconv.Itoa(sum)

	//redis
	redisURL := os.Getenv("REDIS_URL")
	if redisURL != "" && rd.GetRedisInstance() != nil {
		//redis
		return checkRedis(newData)
	} else {
		//open saved log
		return checkFile(newData)
	}
}

//check saved data and run browser if needed
func checkSavedTeachers() {
	ths := th.GetsavedTeachers()
	//open browser
	if len(ths) != 0 {
		//save status
		openFlg := saveStatusLog(ths)
		fmt.Println(openFlg)
		if openFlg {
			if os.Getenv("MAIL_TO_ADDRESS") != "" {
				// for sending mail
				sendMail(ths)
			} else {
				//Browser Mode
				openBrowser(ths)

				//emit a sound
				_ = exec.Command("say", "Found").Start()
			}
		}
	}
}

//handle each teacher data
func handleTeachers(si *th.SiteInfo) {
	defer tm.Track(time.Now(), "handleTeachers()")

	wg := &sync.WaitGroup{}
	chanSemaphore := make(chan bool, MaxGoRoutine)

	//si.Teachers
	for _, teacher := range si.Teachers {
		wg.Add(1)
		chanSemaphore <- true

		//chanSemaphore <- true
		go func(t th.TeacherInfo) {
			defer func() {
				<-chanSemaphore
				wg.Done()
			}()
			//concurrent func
			t.GetHTML(si.Url)
		}(teacher)
	}
	wg.Wait()
}

//initialize Redis
func redisInit() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL != "" {
		host, pass, port, err := hrk.GetRedisInfo(redisURL)
		if err != nil {
			return
		}
		rd.New(host, uint16(port), pass)
		rd.GetRedisInstance().Connection(0)
	}
}

func init() {
	lg.InitializeLog(lg.DEBUG_STATUS, lg.LOG_OFF_COUNT, 0, "[BookingTeacher]", "/var/log/go/book.log")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, os.Args[0]))
	}

	//command-line
	flag.Parse()

	//signal (Debug)
	go signal.StartSignal()

	//Redis
	redisInit()
}

//return value is whether executed scraping or not.
func execMain(testFlg uint8) bool {
	lg.Info("getting teacher's information")

	//heroku mode
	herokuFlg := os.Getenv("HEROKU_FLG")
	if herokuFlg == "1" && (os.Getenv("MAIL_TO_ADDRESS") == "" || os.Getenv("REDIS_URL") == "") {
		lg.Info("environment is not met to run on HEROKU")
		return false
	}

	var si *th.SiteInfo
	//m = new(sync.Mutex)

	if *jsPath != "" {
		//json
		si = th.LoadJsonFile(*jsPath)
	} else {
		//use build teacher data
		si = th.GetDefinedData()
	}

	//Mail Check
	if os.Getenv("MAIL_TO_ADDRESS") != "" {
		settingMail()
	}

	for {
		//reset
		th.InitSavedTeachers()

		//scraping
		handleTeachers(si)

		//save
		checkSavedTeachers()

		//TODO:when integration test, send channel
		//execuite only once on heroku
		if herokuFlg == "1" || testFlg == 1 {
			return true
		}

		time.Sleep(120 * time.Second)
	}
}

// Main
func main() {
	execMain(0)
}
