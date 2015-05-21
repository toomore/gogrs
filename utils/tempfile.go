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

func (hc HttpCache) Get(url string) {
	err := os.Mkdir(hc.Dir, 0700)
	if os.IsExist(err) {
		log.Println(err)
	}
	resp, _ := http.Get(url)
	filehash := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	t, err := os.Create(filepath.Join(hc.Dir, filehash))
	defer t.Close()
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	t.Write(content)
}
