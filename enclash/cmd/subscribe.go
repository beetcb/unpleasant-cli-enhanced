package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const CONF_TUN_WINDOWS = `
tun:
  enable: true
  stack: gvisor
  dns-hijack:
    - 198.18.0.2:53 
  auto-route: true   
  auto-detect-interface: true 
dns:
  enable: true
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16 
  nameserver:
    - 114.114.114.114 
    - 8.8.8.8
`

type ClashConfig map[string]interface{}

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Use:   "sub",
	Short: "Download yaml configuration from clash subscription provider",
	Long:  "Download yaml configuration from clash subscription provider. By default, enclash will use CLASH_SUB_URL environment variable as subscription URL",
	Run: func(cmd *cobra.Command, args []string) {
		tunEnabled, err := cmd.Flags().GetBool("tun")
		if err != nil {
			panic(err)
		}

		// Get remote clash config
		var clashConfOrig ClashConfig
		confRaw, err := grabSubscriptionConfig()
		if err != nil {
			panic(err)
		}
		if err = yaml.Unmarshal(confRaw, &clashConfOrig); err != nil {
			panic(err)
		}

		// Add TUN support
		if tunEnabled {
			var clashTUNConf ClashConfig
			if err := yaml.Unmarshal([]byte(CONF_TUN_WINDOWS), &clashTUNConf); err != nil {
				panic(err)
			}
			for k, v := range clashTUNConf {
				clashConfOrig[k] = v
			}
		}

		// Output yaml configuration
		conf, err := yaml.Marshal(clashConfOrig)
		if os.WriteFile("config.yaml", conf, 0666); err != nil {
			panic(err)
		}

		fmt.Println("Write clash configuration file to `config.yaml`")
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
	subscribeCmd.Flags().BoolP("tun", "t", false, "turn on TUN mode(Windows only, other os WIP)")
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
