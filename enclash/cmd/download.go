package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/beetcb/ghdl"
	h "github.com/beetcb/ghdl/helper"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download clash binary from github release",
	Long:  `Download clash binary from github release`,
	Run: func(cmd *cobra.Command, args []string) {
		ghRelease := ghdl.GHRelease{RepoPath: "Dreamacro/clash", TagName: "premium"}
		ghReleaseDl, err := ghRelease.GetGHReleases(false)
		if err != nil {
			h.Println(fmt.Sprintf("get gh releases failed: %s", err), h.PrintModeErr)
			os.Exit(1)
		}

		h.Println(fmt.Sprintf("start downloading %s", h.Sprint(filepath.Base(ghReleaseDl.Url), h.SprintOptions{PromptOff: true, PrintMode: h.PrintModeSuccess})), h.PrintModeInfo)
		if err := ghReleaseDl.DlTo("."); err != nil {
			h.Println(fmt.Sprintf("download failed: %s", err), h.PrintModeErr)
			os.Exit(1)
		}
		if err := ghReleaseDl.ExtractBinary(); err != nil {
			switch err {
			case ghdl.ErrNeedInstall:
				h.Println(fmt.Sprintf("%s. You can install it with the appropriate commands", err), h.PrintModeInfo)
				os.Exit(0)
			case ghdl.ErrNoBin:
				h.Println(fmt.Sprintf("%s. Try to specify binary name flag", err), h.PrintModeInfo)
				os.Exit(0)
			default:
				h.Println(fmt.Sprintf("extract failed: %s", err), h.PrintModeErr)
				os.Exit(1)
			}
		}
		h.Println(fmt.Sprintf("saved executable to %s", ghReleaseDl.BinaryName), h.PrintModeSuccess)
		if err := os.Chmod(ghReleaseDl.BinaryName, 0777); err != nil {
			h.Println(fmt.Sprintf("chmod failed: %s", err), h.PrintModeErr)
		}

	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
