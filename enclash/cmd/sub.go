package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

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
		webEnabled, err := cmd.Flags().GetBool("web")
		if err != nil {
			panic(err)
		}

		// Get remote clash config
		var clashConf ClashConfig
		confRaw, err := grabSubscriptionConfig()
		if err != nil {
			panic(err)
		}
		if err = yaml.Unmarshal(confRaw, &clashConf); err != nil {
			panic(err)
		}

		// Add TUN support
		if tunEnabled {
			MergeConf(CONF_TUN_DNS, clashConf)
		}

		// Serve web ui
		if webEnabled {
			if _, err := os.Stat(WEB_STATIC_FOLDER); os.IsNotExist(err) {
				err := prepareWebStatic()
				if err != nil {
					panic(err)
				}
			}
			MergeConf(CONF_WEB_UI, clashConf)
		}

		// Output yaml configuration
		conf, err := yaml.Marshal(clashConf)
		fmt.Println("Saving clash configuration file to `config.yaml`")
		if os.WriteFile("config.yaml", conf, 0666); err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
	subscribeCmd.Flags().BoolP("tun", "t", false, "turn on TUN mode(Windows only, other os WIP)")
	subscribeCmd.Flags().BoolP("web", "w", false, "prepare web static files, let calsh serve it")
}
