package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hiromaily/go-book-teacher/pkg/config"
	"github.com/hiromaily/go-book-teacher/pkg/teachers"
)

var (
	jsPath    = flag.String("json", "", "JSON file path")
	tomlPath  = flag.String("toml", "", "TOML file path")
	day       = flag.Int("day", 0, "0: all day, 1:today, 2: tomorrow")
	isVersion = flag.Bool("v", false, "version")
	// -d daemon mode
	version string
)

var usage = `Usage: %s [options...]
Options:
  -json      Json file path for teacher information
  -toml      Toml file path for config
  -day       range of schedule to get teacher's availability: 0: all day, 1:today, 2: tomorrow
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

	// version
	if *isVersion {
		fmt.Printf("%s %s\n", "book-teacher", version)
		os.Exit(0)
	}

	// config
	configPath := *tomlPath
	if configPath == "" {
		configPath = config.GetEnvConfPath()
	}
	conf, err := config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}

	// registry
	regi := NewRegistry(conf)

	// json
	jsonPath := *jsPath
	if jsonPath == "" {
		jsonPath = teachers.GetEnvJSONPath()
	}

	// booker
	booker := regi.NewBooker(jsonPath, *day)
	if err := booker.Start(); err != nil {
		log.Fatal(err)
	}
	//
	booker.Clean()
}
