package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

func MergeConf(source string, clashConf ClashConfig) {
	var clashWebConf ClashConfig
	if err := yaml.Unmarshal([]byte(source), &clashWebConf); err != nil {
		panic(err)
	}
	for k, v := range clashWebConf {
		clashConf[k] = v
	}
}

func grabSubscriptionConfig() ([]byte, error) {
	subUrl := os.Getenv("CLASH_SUB_URL")
	fmt.Println("Downloading clash configuration file from " + subUrl)
	if subUrl != "" {
		resp, err := http.Get(subUrl)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		return body, err
	}
	return nil, fmt.Errorf("CLASH_SUB_URL environment variable not set")
}

func prepareWebStatic() error {
	fmt.Println("Cloning Dreamacro/calsh-dashboard to pwd")
	repo, branch := "https://github.com/Dreamacro/clash-dashboard.git", "gh-pages"
	// fmt.Println("Cloning clash-dashboard from" + )
	// check if clash-dashboard folder exists
	cmd := exec.Command("git", "clone", repo, "-b", branch)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("CLASH_SUB_URL environment variable not set")
	}

	return nil
}
