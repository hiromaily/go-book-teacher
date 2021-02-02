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
  -v         show version
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

func checkVersion() {
	if *isVersion {
		fmt.Printf("%s %s\n", "book-teacher", version)
		os.Exit(0)
	}
}

func getConfig() *config.Root {
	configPath := *tomlPath
	if configPath == "" {
		//book-teacher.toml
		expectedFileName := fmt.Sprintf("%s.toml", os.Args[0])
		log.Println("config file: ", expectedFileName)
		if _, err := os.Stat(expectedFileName); !os.IsNotExist(err) {
			configPath = expectedFileName
		}
	}
	if configPath == "" {
		configPath = config.GetEnvConfPath()
	}
	conf, err := config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}
	return conf
}

func getJSON() string {
	jsonPath := *jsPath
	if jsonPath == "" {
		//book-teacher.json
		expectedFileName := fmt.Sprintf("%s.json", os.Args[0])
		log.Println("json file: ", expectedFileName)
		if _, err := os.Stat(expectedFileName); !os.IsNotExist(err) {
			jsonPath = expectedFileName
		}
	}
	if jsonPath == "" {
		jsonPath = teachers.GetEnvJSONPath()
	}
	return jsonPath
}

func main() {
	parseFlag()
	checkVersion()

	conf := getConfig()
	jsonPath := getJSON()

	// registry
	regi := NewRegistry(conf)

	// booker
	booker := regi.NewBooker(jsonPath, *day)
	if err := booker.Start(); err != nil {
		log.Fatal(err)
	}
	booker.Clean()
}
