package utils

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"io"
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

// NewChaoJiYing 初始化，可以使用代理
func NewChaoJiYing(timeout time.Duration, httpsProxy string) *ChaoJiYing {
	client := &ChaoJiYing{Timeout: timeout, HttpsProxy: httpsProxy}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if httpsProxy != "" {
		proxyURL, err := url.Parse(httpsProxy)
		if err != nil {
			log.Println(err)
		} else {
			tr.Proxy = http.ProxyURL(proxyURL)
		}
	}

	client.HttpClient = &http.Client{Transport: tr, Timeout: timeout}
	return client
}

// GetScore 查询信息
func (client *ChaoJiYing) GetScore(urlString, user, pass string) ([]byte, error) {
	parameters := url.Values{}
	parameters.Add("user", user)
	parameters.Add("pass", pass)

	req, err := http.NewRequest("POST", urlString, strings.NewReader(parameters.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)")
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = client.HttpClient.Do(req)
		if err == nil {
			break
		}
		log.Printf("Request failed, retrying... (%d/3)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// 文件转码base64字符串
func getEncodedBase64(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded, nil
}

// GetPicVal 发出请求获得json结果
func (client *ChaoJiYing) GetPicVal(user, pass, softid, codetype, len_min, filename string) ([]byte, error) {
	urlString := "http://upload.chaojiying.net/Upload/Processing.php"

	encodedFile, err := getEncodedBase64(filename)
	if err != nil {
		return nil, err
	}

	parameters := url.Values{}
	parameters.Add("user", user)
	parameters.Add("pass", pass)
	parameters.Add("softid", softid)
	parameters.Add("codetype", codetype)
	parameters.Add("len_min", len_min)
	parameters.Add("file_base64", encodedFile)

	req, err := http.NewRequest("POST", urlString, strings.NewReader(parameters.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0)")
	req.Header.Set("Connection", "Keep-Alive")

	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = client.HttpClient.Do(req)
		if err == nil {
			break
		}
		log.Printf("Request failed, retrying... (%d/3)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
