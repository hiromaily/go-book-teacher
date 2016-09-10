package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	conf "github.com/hiromaily/go-book-teacher/config"
	th "github.com/hiromaily/go-book-teacher/teacher"
	enc "github.com/hiromaily/golibs/cipher/encryption"
	rd "github.com/hiromaily/golibs/db/redis"
	hrk "github.com/hiromaily/golibs/heroku"
	lg "github.com/hiromaily/golibs/log"
	ml "github.com/hiromaily/golibs/mails"
	"github.com/hiromaily/golibs/signal"
	tm "github.com/hiromaily/golibs/times"
	"github.com/hiromaily/golibs/tmpl"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

// MaxGoRoutine is number of goroutine running at the same time
const MaxGoRoutine uint16 = 20

var (
	jsPath   = flag.String("j", "", "Json file path")
	tomlPath = flag.String("t", "", "Toml file path")
	interval = flag.Int64("i", 120, "Interval for scraping")
)

var usage = `Usage: %s [options...]
Options:
  -j     Json file path
  -t     Toml file path
  -i     Interval for scraping
`

var (
	mi *ml.MailInfo

	tmplMails = `
The following tachers are available now!
{{range .Teachers}}
{{$.Url}}teacher/index/{{.Id}}/ [{{.Name}} / {{.Country}}]
{{end}}
Enjoy!`

	redisKey      = "bookteacher:save"
	savedFilePath = "/tmp/status.log"

	//judge ment
	herokuFlg = os.Getenv("HEROKU_FLG")
	mailFlg   = false
	redisFlg  = false
)

//ENVIRONMENT VARIABLE
//HEROKU_FLG
//ENC_KEY
//ENC_IV

// cipher settings
func cipherSetup() {
	size := 16
	key := os.Getenv("ENC_KEY")
	iv := os.Getenv("ENC_IV")

	if key == "" || iv == "" {
		panic("set Environment Variable: ENC_KEY, ENC_IV")
	}

	enc.NewCrypt(size, key, iv)
}

// setting for sending mail
func settingMail() {
	//get environment variable
	subject := "[ENGLISH LESSON] It's Available."
	body := ""
	//mails
	smt := conf.GetConf().Mail.SMTP
	m := conf.GetConf().Mail

	smtp := ml.Smtp{Address: smt.Address, Pass: smt.Pass,
		Server: smt.Server, Port: smt.Port}

	mi = &ml.MailInfo{ToAddress: []string{m.MailTo}, FromAddress: m.MailFrom,
		Subject: subject, Body: body, Smtp: smtp}
}

func settingSavedFile() {
	if conf.GetConf().StatusFile != "" {
		savedFilePath = conf.GetConf().StatusFile
	}
}

// send mail
func sendMail(ths []th.Info) {
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
func openBrowser(ths []th.Info) {
	for _, t := range ths {
		//during test, it doesn't work.
		//open browser error: exec: "open": executable file not found in $PATH

		//out, err := exec.Command("open /Applications/Google\\ Chrome.app", fmt.Sprintf("http://eikaiwa.dmm.com/teacher/index/%d/", id)).Output()
		err := exec.Command("open", fmt.Sprintf("http://eikaiwa.dmm.com/teacher/index/%d/", t.ID)).Start()
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

	if err != nil {
		lg.Errorf("redis error is %s\n", err)
	}
	lg.Debugf("new value is %s, old value is %s\n", newData, val)

	if err != nil || newData != val {
		//save
		c.Do("SET", redisKey, newData)
		return true
	}
	return false
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
	fp, err := os.OpenFile(savedFilePath, os.O_CREATE, 0664)
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
	ioutil.WriteFile(savedFilePath, content, 0664)
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
func saveStatusLog(ths []th.Info) bool {

	//create string from ids slice
	var sum int
	for _, t := range ths {
		sum += t.ID
	}
	newData := strconv.Itoa(sum)

	//redis
	if redisFlg && rd.GetRedisInstance() != nil {
		//redis
		return checkRedis(newData)
	}

	//open saved log
	return checkFile(newData)
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
			if mailFlg {
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
		go func(t th.Info) {
			defer func() {
				<-chanSemaphore
				wg.Done()
			}()
			//concurrent func
			t.GetHTML(si.URL)
		}(teacher)
	}
	wg.Wait()
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

	//cipher
	cipherSetup()

	//config
	if *tomlPath != "" {
		conf.SetTomlPath(*tomlPath)
	}
	conf.New("")
	conf.Cipher()
}

func setupMain() {
	//saved file
	settingSavedFile()

	//flg
	if conf.GetConf().Redis.URL != "" {
		redisFlg = true

		//Redis
		settingRedis()
	}
	if conf.GetConf().Mail.MailTo != "" {
		mailFlg = true

		//Mail Check
		settingMail()
	}

	//th.SetPrintOn(true)
}

//initialize Redis
func settingRedis() {
	redisURL := conf.GetConf().Redis.URL
	host, pass, port, err := hrk.GetRedisInfo(redisURL)
	if err != nil {
		return
	}
	rd.New(host, uint16(port), pass)
	rd.GetRedisInstance().Connection(0)
}

func checkHeroku() error {
	//heroku mode
	if herokuFlg == "1" && (!mailFlg || !redisFlg) {
		return fmt.Errorf("%s", "mail settings is required for HEROKU")
	}
	return nil
}

//return value is whether executed scraping or not.
func execMain(testFlg uint8) bool {
	lg.Info("getting teacher's information")

	var si *th.SiteInfo
	//m = new(sync.Mutex)

	if *jsPath != "" {
		//json
		si = th.LoadJSONFile(*jsPath)
	} else {
		//use build teacher data
		si = th.GetDefinedData()
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

		time.Sleep(time.Duration(*interval) * time.Second)
	}
}

// Main
func main() {
	setupMain()

	err := checkHeroku()
	if err != nil {
		panic(err)
	}

	execMain(0)
}
