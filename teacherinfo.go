package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	//"path/filepath"
)

type teachrInfo struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

type siteInfo struct {
	Url      string       `json:"url"`
	Teachers []teachrInfo `json:"teachers"`
}

func LoadJsonFile() *siteInfo {
	// Loading jsonfile

	dir := path.Dir(os.Args[0])
	fmt.Println(dir)

	//jsonPath, _ := filepath.Abs("./settings.json")
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/settings.json", dir))

	var siteInfo siteInfo
	err = json.Unmarshal(file, &siteInfo)

	if err != nil {
		fmt.Println("Format Error: ", err)
	}
	//fmt.Println(siteInfo)
	fmt.Println(siteInfo.Url)
	fmt.Println(siteInfo.Teachers[0].Id)
	fmt.Println(siteInfo.Teachers[0].Name)
	fmt.Println(siteInfo.Teachers[0].Country)

	return &siteInfo
}
