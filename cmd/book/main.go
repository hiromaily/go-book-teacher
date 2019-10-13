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

//ENVIRONMENT VARIABLE
//ENC_KEY
//ENC_IV

var (
	jsPath          = flag.String("j", "", "Json file path")
	tomlPath        = flag.String("t", "", "Toml file path")
	interval        = flag.Int("i", 0, "Interval for scraping") //if value is 0, it scrapes only once
	isEncryptedConf = flag.Bool("crypto", false, "if true, config file is handled as encrypted value")
)

var usage = `Usage: %s [options...]
Options:
  -j      Json file path for teacher information
  -t      Toml file path for config
  -i      Interval for scraping, if 0 it scrapes only once
  -crypto true is that conf file is handled as encrypted value
`

func init() {
	//log
	lg.InitializeLog(lg.DebugStatus, lg.TimeShortFile, "[BookingTeacher]", "", "hiromaily")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, os.Args[0]))
	}

	//command-line
	flag.Parse()
}

// Main
func main() {
	//cipher
	if *isEncryptedConf {
		_, err := enc.NewCryptWithEnv()
		if err != nil {
			panic(err)
		}
	}

	//config
	if err := config.New(*tomlPath, *isEncryptedConf); err != nil {
		panic(err)
	}

	//signal (Debug)
	go signal.StartSignal()

	regi := NewRegistry(config.GetConf())
	booker := regi.NewBooker(*jsPath, *interval)
	if err := booker.Start(); err != nil {
		lg.Error(err)
	}
	//
	booker.Cleanup()
}
