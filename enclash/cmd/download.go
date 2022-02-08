package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/unpleasant-cli-enhanced/enclash/helper"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download clash binary from github release",
	Long:  `Download clash binary from github release`,
	Run: func(cmd *cobra.Command, args []string) {
		url, binaryName, err := helper.GetReleases(&helper.GitRepo{RepoPath: "Dreamacro/clash", TagName: "premium"})
		if err != nil {
			panic(err)
		}
		binaryNameFlag, err := cmd.Flags().GetString("bin")
		if err != nil {
			panic(err)
		}
		if binaryNameFlag != "" {
			binaryName = binaryNameFlag
		}
		dlAndDecompression(url, binaryName)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("bin", "b", "", "Specify bin file name")
}

func dlAndDecompression(url string, binaryName string) {
	b := binaryName + func() string {
		if runtime.GOOS == "windows" {
			return ".exe"
		} else {
			return ""
		}
	}()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	bytesReader, fileExt := bytes.NewReader(body), filepath.Ext(url)
	var decompressedBinary *[]byte
	switch fileExt {
	case ".zip":
		decompressedBinary, err = zipBinary(bytesReader, b)
		if err != nil {
			panic(err)
		}
	case ".gz":
		if strings.Contains(url, ".tar.gz") {
			decompressedBinary, err = targzBinary(bytesReader, b)
			if err != nil {
				panic(err)
			}
		} else {
			decompressedBinary, err = gzBinary(bytesReader, b)
			if err != nil {
				panic(err)
			}
		}
	case ".deb":
	case ".rpm":
	case ".apk":
		fileName := b + fileExt
		fmt.Printf("Detected deb/rpm/apk package, download directly to ./%s\nYou can install it with the appropriate commands\n", fileName)
		if err := os.WriteFile(fileName, body, 0777); err != nil {
			panic(err)
		}
	case "":
		decompressedBinary = &body
	default:
		panic("unsupported file format")
	}
	if err := os.WriteFile(b, *decompressedBinary, 0777); err != nil {
		panic(err)
	}
}

func zipBinary(r *bytes.Reader, b string) (*[]byte, error) {
	zipR, err := zip.NewReader(r, int64(r.Len()))
	if err != nil {
		return nil, err
	}

	for _, f := range zipR.File {
		if filepath.Base(f.Name) == b || len(zipR.File) == 1 {
			open, err := f.Open()
			if err != nil {
				return nil, err
			}
			ret, err := ioutil.ReadAll(open)
			if err != nil {
				return nil, err
			}
			return &ret, err
		}
	}
	return nil, fmt.Errorf("Binary file %v not found", b)
}

func gzBinary(r *bytes.Reader, b string) (*[]byte, error) {
	gzR, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzR.Close()
	ret, err := ioutil.ReadAll(gzR)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func targzBinary(r *bytes.Reader, b string) (*[]byte, error) {
	gzR, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gzR.Close()
	tarR := tar.NewReader(gzR)

	var file []byte
	for {
		header, err := tarR.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if (header.Typeflag != tar.TypeDir) && filepath.Base(header.Name) == b {
			file, err = ioutil.ReadAll(tarR)
			if err != nil {
				return nil, err
			}
		}
	}
	return &file, nil
}
