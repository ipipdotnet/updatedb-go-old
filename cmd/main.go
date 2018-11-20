package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ipipdotnet/updatedb-go"
	"github.com/spf13/pflag"
)

var (
	token    string
	fileType string
	compress bool
	language string

	dirPath string
)

func init() {
	pflag.StringVar(&token, "token", "", "--token=XXX")
	pflag.StringVar(&fileType, "type", "", "--type=ipdb|txtx")
	pflag.BoolVar(&compress, "compress", true, "--compress")
	pflag.StringVar(&language, "lang", "", "-lang=EN|CN")

	pflag.StringVar(&dirPath, "dir", "", "-dir=/web")
	pflag.Parse()

	if pflag.NArg() == 0 {
		pflag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {

	if len(token) != 40 {
		fmt.Println()
		fmt.Println("Token error")
		fmt.Println()
		os.Exit(0)
	}

	if fileType != "ipdb" && fileType != "txtx" && fileType != "txt" {
		fmt.Println()
		fmt.Println("file type no support")
		fmt.Println()
		os.Exit(0)
	}

	info, err := os.Stat(dirPath)
	if err != nil {
		fmt.Println()
		fmt.Println(err)
		fmt.Println()
		os.Exit(0)
	}
	if !info.IsDir() {
		fmt.Println()
		fmt.Println(dirPath, "It's not a directory")
		fmt.Println()
		os.Exit(0)
	}

	if err := os.MkdirAll(dirPath, 0666); err != nil {
		fmt.Println(err)
		fmt.Println(dirPath, "is not writeable")
		fmt.Println()
		os.Exit(0)
	}

	api := updatedb.BuildURL(token, fileType, language, compress)

	retry := 3

RETRY:
	err = updatedb.Download(api.String(), dirPath, "")
	if err == updatedb.ErrNetwork {
		if retry > 0 {
			retry--
			time.Sleep(time.Minute)
			goto RETRY
		}
	} else if err != nil {
		fmt.Println()
		fmt.Println("download failed")
		fmt.Println(err)
		fmt.Println()
	} else {
		fmt.Println()
		fmt.Println("download ok")
		fmt.Println()
	}
}
