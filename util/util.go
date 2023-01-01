package util

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

func MergeMap(m1 map[string][]byte, m2 map[string][]byte) {
	for k, v := range m2 {
		m1[k] = v
	}
}

func Zip(zipFileName string, files map[string][]byte) {
	f, err := os.Create(zipFileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	zipWriter := zip.NewWriter(f)
	for name, file := range files {
		zw, err := zipWriter.Create(name)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(zw, bytes.NewReader(file)); err != nil {
			log.Fatal(err)
		}
	}
	zipWriter.Close()
}

func HttpDownloadZip(w http.ResponseWriter, files map[string][]byte) {
	zipWriter := zip.NewWriter(w)
	for name, file := range files {
		zw, err := zipWriter.Create(name)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(zw, bytes.NewReader(file)); err != nil {
			log.Fatal(err)
		}
	}
	zipWriter.Close()
}

func MD5(content []byte) string {
	hash := md5.New()
	hash.Write(content)
	md5Bytes := hash.Sum(nil)
	return hex.EncodeToString(md5Bytes)
}

// 从图片src中解析出图片的扩展名
func ParseImageExtFromSrc(src string) string {
	reg := regexp.MustCompile(`(wx_fmt=)([a-zA-Z]+)(&?)`)
	matches := reg.FindStringSubmatch(src)
	if len(matches) < 3 {
		return ""
	}
	return matches[2]
}

// 判断路径是否存在
func PathIsExists(path string) (os.FileInfo, bool) {
	f, err := os.Stat(path)
	return f, err == nil || os.IsExist(err)
}
