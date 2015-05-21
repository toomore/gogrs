package utils

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type HttpCache struct {
	Dir  string
	Rand bool
}

func NewHttpCache(dir string, rand bool) *HttpCache {
	return &HttpCache{Dir: dir, Rand: rand}
}

func (hc HttpCache) Get(url string) ([]byte, error) {
	err := os.Mkdir(hc.Dir, 0700)
	if os.IsExist(err) {
		log.Println(err)
	}
	filehash := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	content, err := hc.readFile(filehash)
	if err != nil {
		return hc.saveFile(url, filehash)
	}
	return content, nil
}

func (hc HttpCache) readFile(filehash string) ([]byte, error) {
	f, err := os.Open(filepath.Join(hc.Dir, filehash))
	defer f.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func (hc HttpCache) saveFile(url, filehash string) ([]byte, error) {
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
