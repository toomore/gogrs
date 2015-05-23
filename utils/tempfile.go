package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	iconv "github.com/djimenez/iconv-go"
)

// TempFolderName 快取資料夾名稱
const TempFolderName = ".gogrscache"

// HTTPCache net/http 快取功能
type HTTPCache struct {
	Dir            string
	iconvConverter func([]byte) []byte
}

// NewHTTPCache New 一個 HTTPCache.
//
// dir 為暫存位置，fromEncoding 來源檔案的編碼，一律轉換為 utf8
func NewHTTPCache(dir string, fromEncoding string) *HTTPCache {
	err := os.Mkdir(dir, 0700)
	if os.IsNotExist(err) {
		dir = filepath.Join(os.TempDir(), TempFolderName)
		os.Mkdir(dir, 0700)
	}
	return &HTTPCache{Dir: dir, iconvConverter: renderIconvConverter(fromEncoding)}
}

// Get 透過 http.Get 取得檔案或從暫存中取得檔案
//
// rand 為是否支援網址帶入亂數值，url 需有 '%d' 格式。
func (hc HTTPCache) Get(url string, rand bool) ([]byte, error) {
	filehash := fmt.Sprintf("%x", md5.Sum([]byte(url)))
	content, err := hc.readFile(filehash)
	if err != nil {
		return hc.saveFile(url, filehash, rand, nil)
	}
	return content, nil
}

// PostForm 透過 http.PostForm 取得檔案或從暫存中取得檔案
func (hc HTTPCache) PostForm(url string, data url.Values) ([]byte, error) {
	hash := md5.New()
	io.WriteString(hash, url)
	io.WriteString(hash, data.Encode())

	filehash := fmt.Sprintf("%x", hash.Sum(nil))
	content, err := hc.readFile(filehash)
	if err != nil {
		return hc.saveFile(url, filehash, false, data)
	}
	return content, nil
}

// readFile 從快取資料裡面取得
func (hc HTTPCache) readFile(filehash string) ([]byte, error) {
	f, err := os.Open(filepath.Join(hc.Dir, filehash))
	defer f.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

var transport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   0,
		KeepAlive: 0,
	}).Dial,
	TLSHandshakeTimeout: 10 * time.Second,
}

var httpClient = &http.Client{Transport: transport}

// saveFile 從網路取得資料後放入快取資料夾
func (hc HTTPCache) saveFile(url, filehash string, rand bool, data url.Values) ([]byte, error) {
	if rand {
		url = fmt.Sprintf(url, RandInt())
	}
	var resp *http.Response
	var req *http.Request
	var err error
	if len(data) == 0 {
		// http.Get
		req, _ = http.NewRequest("GET", url, nil)
		req.Header.Set("Connection", "close")
	} else {
		// http.PostForm
		req, _ = http.NewRequest("POST", url, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Connection", "close")
	}

	resp, _ = httpClient.Do(req)
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)

	f, _ := os.Create(filepath.Join(hc.Dir, filehash))
	defer f.Close()

	out := hc.iconvConverter(content)
	f.Write(out)

	return out, err
}

// renderIconvConverter wrapper function for iconv converter.
func renderIconvConverter(fromEncoding string) func([]byte) []byte {
	if fromEncoding == "utf8" || fromEncoding == "utf-8" {
		return func(str []byte) []byte {
			return str
		}
	}
	return func(content []byte) []byte {
		converter, _ := iconv.NewConverter(fromEncoding, "utf-8")
		var out []byte
		out = make([]byte, len(content)*2)
		_, outLen, _ := converter.Convert(content, out)
		return out[:outLen]
	}
}
