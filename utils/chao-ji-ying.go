package utils

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type ChaoJiYing struct {
	Timeout    time.Duration
	HttpsProxy string
	HttpClient *http.Client
}

// InitWithOptions 初始化，可以使用代理
func (client *ChaoJiYing) InitWithOptions() {
	//使用https，设置不验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	//设置代理
	if client.HttpsProxy != "" {
		proxyURL, err := url.Parse(client.HttpsProxy)
		if err != nil {
			log.Println(err)
		} else {
			tr.Proxy = http.ProxyURL(proxyURL)
		}
	}
	client.HttpClient = &http.Client{Transport: tr}
	client.HttpClient.Timeout = 1 * time.Minute
}

// GetScore 查询信息
func (client *ChaoJiYing) GetScore(urlString string, user string, pass string) []byte {
	var req *http.Request
	var resp *http.Response
	var err error
	var body []byte

	parameters := url.Values{}
	parameters.Add("user", user)
	parameters.Add("pass", pass)

	req, err = http.NewRequest("POST", urlString, strings.NewReader(parameters.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)")
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("content: %s\n", string(body))
	return body
}

// 文件转码base64字符串
func getEncodedBase64(filename string) string {
	f, _ := os.Open(filename)
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded
}

// 发出请求获得json结果
func (client *ChaoJiYing) GetPicVal(user string, pass string, softid string, codetype string,
	len_min string, filename string) []byte {
	var req *http.Request
	var resp *http.Response
	var err error
	var body []byte
	urlString := "http://upload.chaojiying.net/Upload/Processing.php"

	parameters := url.Values{}
	parameters.Add("user", user)
	parameters.Add("pass", pass)
	parameters.Add("softid", softid)
	//http://www.chaojiying.com/price.html
	parameters.Add("codetype", codetype)
	parameters.Add("len_min", len_min)
	parameters.Add("file_base64", getEncodedBase64(filename))

	req, err = http.NewRequest("POST", urlString, strings.NewReader(parameters.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)")
	req.Header.Set("Connection", "Keep-Alive")

	resp, err = client.HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("content: %s\n", string(body))
	return body
}
