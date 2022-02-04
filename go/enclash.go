package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

type ClashConfigCustom map[string]interface{}

func main() {
	arg := os.Args[1]
	data, err := os.ReadFile(arg)
	check(err)

	var clashMerge ClashConfigCustom
	err = yaml.Unmarshal(data, &clashMerge)
	check(err)

	// Get remote clash config
	var clashConfOrig ClashConfigCustom
	confRaw, err := grabSubscriptionConfig()
	check(err)
	err = yaml.Unmarshal(confRaw, &clashConfOrig)
	check(err)

	for k, v := range clashMerge {
		clashConfOrig[k] = v
	}

	// Output yaml configuration
	conf, err := yaml.Marshal(clashConfOrig)
	err = os.WriteFile("config.yaml", conf, 1)
	check(err)
	fmt.Println("Write merged file to `config.yaml`")
}

func grabSubscriptionConfig() ([]byte, error) {
	subUrl := os.Getenv("CLASH_SUB_URL")
	fmt.Println(subUrl)
	if subUrl != "" {
		resp, err := http.Get(subUrl)
		check(err)
		body, err := ioutil.ReadAll(resp.Body)
		check(err)
		return body, err
	}
	return nil, fmt.Errorf("CLASH_SUB_URL not set")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
