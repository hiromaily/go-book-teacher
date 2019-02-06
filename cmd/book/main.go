package main

import (
	"flag"
	"fmt"
	conf "github.com/hiromaily/go-book-teacher/pkg/config"
	ioo "github.com/hiromaily/go-book-teacher/pkg/io"
	ml "github.com/hiromaily/go-book-teacher/pkg/mail"
	rd "github.com/hiromaily/go-book-teacher/pkg/redis"
	sl "github.com/hiromaily/go-book-teacher/pkg/slack"
	th "github.com/hiromaily/go-book-teacher/pkg/teacher"
	tt "github.com/hiromaily/go-book-teacher/pkg/text"
	enc "github.com/hiromaily/golibs/cipher/encryption"
	lg "github.com/hiromaily/golibs/log"
	"github.com/hiromaily/golibs/signal"
	tm "github.com/hiromaily/golibs/time"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

//ENVIRONMENT VARIABLE
//HEROKU_FLG
//ENC_KEY
//ENC_IV

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
	//judge ment
	herokuFlg = os.Getenv("HEROKU_FLG")
	mailFlg   = false
	slackFlg  = false
	redisFlg  = false
)

func checkHeroku() error {
	//heroku mode
	if herokuFlg == "1" && ((!mailFlg && !slackFlg) || !redisFlg) {
		return fmt.Errorf("%s", "mail settings is required for HEROKU")
	}
	return nil
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

func clearData(s ioo.Deleter) {
	err := s.Delete()
	lg.Error(err)
}

//check saved data and run browser if needed
func checkSavedTeachers() {
	var openFlg bool
	ths := th.GetsavedTeachers()
	//open browser
	if len(ths) != 0 {
		//save status
		if redisFlg {
			openFlg = saveStatusLog(rd.Get(), ths)
		} else {
			openFlg = saveStatusLog(tt.Get(), ths)
		}
		fmt.Println(openFlg)
		if openFlg {
			if mailFlg {
				// for sending mail
				ml.Send(ths)
			} else if slackFlg {
				// for sending slack
				err := sl.Send(ths)
				if err != nil {
					lg.Error(err)
				}
			} else {
				//Browser Mode
				openBrowser(ths)

				//emit a sound
				_ = exec.Command("say", "Found").Start()
			}
		}
	}
}

//save teacher status to log
// check how to use interface
func saveStatusLog(s ioo.Saver, ths []th.Info) bool {

	//create string from ids slice
	var sum int
	for _, t := range ths {
		sum += t.ID
	}
	newData := strconv.Itoa(sum)

	return s.Save(newData)
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
	lg.InitializeLog(lg.DebugStatus, lg.TimeShortFile,  "[BookingTeacher]", "/var/log/go/book.log", "hiromaily")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, os.Args[0]))
	}

	//command-line
	flag.Parse()

	//signal (Debug)
	go signal.StartSignal()

	//cipher
	_, err := enc.NewCryptWithEnv()
	if err != nil {
		panic(err)
	}

	//config
	conf.New(*tomlPath, true)
	//conf.Cipher()
}

//1.setupMain()
func setupMain() {

	//flg
	if conf.GetConf().Redis.URL != "" {
		//Redis
		_, err := rd.Setup()
		if err == nil {
			redisFlg = true
		}
	}
	//saved file
	tt.Setup()

	if conf.GetConf().Mail.MailTo != "" {
		mailFlg = true

		//Mail Check
		ml.Setup()
	}

	if conf.GetConf().Slack.Key != "" {
		slackFlg = true
	}

	//th.SetPrintOn(true)
}

//2.execMain() return value is whether executed scraping or not.
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

		if herokuFlg == "1" {
			if rd.Get() != nil {
				rd.Get().RD.Close()
			}
			return true
		} else if testFlg == 1 {

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
