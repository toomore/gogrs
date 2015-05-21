package utils

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// HTTPCache net/http 快取功能
type HTTPCache struct {
	Dir string
}

// NewHTTPCache New 一個 HTTPCache.
//
// dir 為暫存位置
func NewHTTPCache(dir string) *HTTPCache {
	err := os.Mkdir(dir, 0700)
	if os.IsNotExist(err) {
		dir = filepath.Join(os.TempDir(), ".gogrs")
		os.Mkdir(dir, 0700)
	}
	return &HTTPCache{Dir: dir}
}

// Get 透過 http.Get 取得檔案或從暫存中取得檔案
//
// rand 為是否支援網址帶入亂數值，url 需有 '%d' 格式。
func (hc HTTPCache) Get(url string, rand bool) ([]byte, error) {
	filehash := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	content, err := hc.readFile(filehash)
	if err != nil {
		return hc.saveFile(url, filehash, rand)
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

func (hc HTTPCache) saveFile(url, filehash string, rand bool) ([]byte, error) {
	if rand {
		url = fmt.Sprintf(url, RandInt())
	}
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	f, err := os.Create(filepath.Join(hc.Dir, filehash))
	defer f.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		f.Write(content)
	}
	return content, err
}
