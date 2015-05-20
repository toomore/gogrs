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

func GetWithTemp(url string) {
	var dir = ".temp"
	err := os.Mkdir(dir, 0700)
	if os.IsExist(err) {
		log.Println(err)
	}
	resp, _ := http.Get(url)
	filehash := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	t, err := os.Create(filepath.Join(dir, filehash))
	defer t.Close()
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	t.Write(content)
}
