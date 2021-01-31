package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hiromaily/go-book-teacher/pkg/config"
	enc "github.com/hiromaily/golibs/cipher/encryption"
	lg "github.com/hiromaily/golibs/log"
	"github.com/hiromaily/golibs/signal"
)

// ENVIRONMENT VARIABLE
// ENC_KEY
// ENC_IV

// -d daemon mode

var (
	jsPath          = flag.String("j", "", "Json file path")
	tomlPath        = flag.String("t", "", "Toml file path")
	interval        = flag.Int("i", 0, "Interval for scraping (xxx second)") // if value is 0, it scrapes only once
	isEncryptedConf = flag.Bool("crypto", false, "if true, config file is handled as encrypted value")
	day             = flag.Int("d", 0, "0: all day, 1:today, 2: tommorw")
)

var usage = `Usage: %s [options...]
Options:
  -j      Json file path for teacher information
  -t      Toml file path for config
  -i      Interval for scraping, if 0 it scrapes only once
  -d      Day for teacher schedule list
  -crypto true is that conf file is handled as encrypted value
`

// init() can not be used because it affects main_test.go as well.
func init() {
}

func parseFlag() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
	}

	// command-line
	flag.Parse()
}

// Main
func main() {
	parseFlag()

	// log
	lg.InitializeLog(lg.DebugStatus, lg.TimeShortFile, "[BookingTeacher]", "", "hiromaily")

	// cipher
	if *isEncryptedConf {
		_, err := enc.NewCryptWithEnv()
		if err != nil {
			panic(err)
		}
	}

	// config
	conf, err := config.NewConfig(*tomlPath, *isEncryptedConf)
	if err != nil {
		panic(err)
	}

	// signal (Debug)
	go signal.StartSignal()

	regi := NewRegistry(conf)
	booker := regi.NewBooker(*jsPath, *day, *interval)
	if err := booker.Start(); err != nil {
		lg.Error(err)
	}
	//
	booker.Cleanup()
}
