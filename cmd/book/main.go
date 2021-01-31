package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/signal"
)

var (
	jsPath   = flag.String("json", "", "JSON file path")
	tomlPath = flag.String("toml", "", "TOML file path")
	day      = flag.Int("day", 0, "0: all day, 1:today, 2: tomorrow")
	// -d daemon mode
)

var usage = `Usage: %s [options...]
Options:
  -json      Json file path for teacher information
  -toml      Toml file path for config
  -day       Day for teacher schedule list
`

// init() can not be used because it affects main_test.go as well.
func init() {
}

func parseFlag() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
	}
	flag.Parse()
}

func main() {
	parseFlag()

	// config
	configPath := *tomlPath
	if configPath == "" {
		configPath = config.GetEnvConfPath()
	}
	conf, err := config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}

	// signal
	go signal.StartSignal()

	regi := NewRegistry(conf)
	booker := regi.NewBooker(*jsPath, *day)
	if err := booker.Start(); err != nil {
		log.Fatal(err)
	}
	//
	booker.Cleanup()
}
