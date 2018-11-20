package update

import (
	"os"
	"io"
	"io/ioutil"
	"fmt"
	"regexp"
	"path/filepath"
	"strings"
	"crypto/sha1"
	"encoding/hex"
	"archive/zip"
	"net/http"
	"net/url"
	"errors"
)

var (
	ErrNetwork = errors.New("network err")
	ErrUnzip = errors.New("unzip err")
	ErrDownloadLimited = errors.New("download limited")
)

func BuildURL(Token, fileType, language string, compress bool) *url.URL {
	downloadURL := &url.URL{
		Scheme: "https",
		Host: "user.ipip.net",
		Path: "download.php",
	}

	params := url.Values{}

	if len(Token) > 0 {
		params.Add("token", Token)
	}
	if len(fileType) > 0 {
		params.Add("type", fileType)
	}
	if fileType == "ipdb" {
		compress = false
	}
	if compress {
		params.Add("ext", "zip")
	}
	if language == "CN" || language == "EN" {
		params.Add("lang", language)
	}

	downloadURL.RawQuery = params.Encode()

	return downloadURL
}

func unzip(dst, src string) error {

	reader, err := zip.OpenReader(src)
	if err != nil {
		return ErrUnzip
	}

	defer reader.Close()

	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return ErrUnzip
		}

		defer rc.Close()

		w, err := os.Create(dst)
		if err != nil {
			return ErrUnzip
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return ErrUnzip
		}
	}

	return nil
}

func Download(api, dirPath, fileName string) error {

	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ErrNetwork
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		return ErrDownloadLimited
	}
	if resp.StatusCode != 200 {
		return ErrNetwork
	}

	tmpFile, err := ioutil.TempFile(dirPath, "ipip-")
	if err != nil {
		return err
	}

	tmp := tmpFile.Name()
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmp)
		return err
	}
	tmpFile.Close()
	all, err := ioutil.ReadFile(tmp)
	if err != nil{
		return err
	}

	fileHash := fmt.Sprintf("sha1-%s", sha1EncodeToString(all))
	respTag := resp.Header.Get("ETag")
	if !strings.EqualFold(respTag, fileHash) {
		os.Remove(tmp)
		return fmt.Errorf("%s [%s]-[%s]", "sha1 error", fileHash, respTag)
	}

	var newName string

	if len(fileName) == 0 {
		attachment := resp.Header.Get("Content-Disposition")
		g := regexp.MustCompile(`filename="(.+?)"`).FindAllStringSubmatch(attachment, -1)
		if len(g) == 0 {
			return fmt.Errorf("%s", "download http response header error")
		}
		newName = filepath.Join(dirPath, g[0][1])
	} else {
		newName = filepath.Join(dirPath, fileName)
	}

	if strings.HasSuffix(newName, ".zip") {
		err = unzip(newName[0:len(newName) - 4], tmp)
		os.Remove(tmp)
		return err
	}

	if err = os.Rename(tmp, newName); err != nil {
		os.Remove(tmp)
		return err
	}

	return nil
}

func sha1EncodeToString(all []byte) string {

	s := sha1.New()
	s.Write(all)

	return hex.EncodeToString(s.Sum(nil))
}