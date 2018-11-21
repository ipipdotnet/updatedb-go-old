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
	pflag.StringVar(&fileType, "type", "ipdb", "--type=ipdb|txtx")
	pflag.BoolVar(&compress, "compress", true, "--compress")
	pflag.StringVar(&language, "lang", "CN", "-lang=EN|CN")

	pflag.StringVar(&dirPath, "dir", "", "-dir=/tmp")
	pflag.Parse()

	if pflag.NFlag() == 0 {
		example := "\nExample: \n\t./updatedb --dir=/tmp --type=ipdb --token=XXX\n"
		fmt.Fprintln(os.Stderr, example)
		pflag.Usage()
		os.Exit(1)
	}
}

func main() {

	if len(token) != 40 {
		fmt.Println("Token error")
		fmt.Println()
		os.Exit(1)
	}

	if fileType != "ipdb" && fileType != "txtx" && fileType != "txt" {
		fmt.Println("file type no support")
		fmt.Println()
		os.Exit(1)
	}

	info, err := os.Stat(dirPath)
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Println(dirPath, "It's not a directory")
		fmt.Println()
		os.Exit(1)
	}

	if err := os.MkdirAll(dirPath, 0666); err != nil {
		fmt.Println(err)
		fmt.Println(dirPath, "is not writeable")
		fmt.Println()
		os.Exit(1)
	}
	retry := 3
	api := updatedb.BuildURL(token, fileType, language, compress)
RETRY:
	fn, err := updatedb.Download(api.String(), dirPath, "")
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
		os.Exit(1)
	} else {
		fmt.Println(fn, "\tdownload ok")
		fmt.Println()
		os.Exit(0)
	}
}
