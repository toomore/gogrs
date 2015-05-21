package utils

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// HTTPCache Dir 為暫存位置，Rand 為是否支援網址帶入亂數值
type HTTPCache struct {
	Dir  string
	Rand bool
}

// NewHTTPCache New 一個 HTTPCache
func NewHTTPCache(dir string, rand bool) *HTTPCache {
	os.Mkdir(dir, 0700)
	return &HTTPCache{Dir: dir, Rand: rand}
}

// Get 透過 http.Get 取得檔案或從暫存中取得檔案
func (hc HTTPCache) Get(url string) ([]byte, error) {
	filehash := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	content, err := hc.readFile(filehash)
	if err != nil {
		return hc.saveFile(url, filehash)
	}
	return content, nil
}

func (hc HTTPCache) readFile(filehash string) ([]byte, error) {
	f, err := os.Open(filepath.Join(hc.Dir, filehash))
	defer f.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func (hc HTTPCache) saveFile(url, filehash string) ([]byte, error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	f, err := os.Create(filepath.Join(hc.Dir, filehash))
	defer f.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		f.Write(content)
	}
	return content, err
}
