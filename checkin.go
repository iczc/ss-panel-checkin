package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	baseUrl = ""
)

func main() {
	email := ""
	passwd := ""
	cookies := login(email, passwd)
	//fmt.Print(cookies)
	checkin(cookies)
}

func login(email, passwd string) []*http.Cookie {
	type RespJson struct {
		Ret int    `json:"ret"`
		Msg string `json:"msg"`
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	postParam := url.Values{
		"email":  {email},
		"passwd": {passwd},
	}
	resp, err := client.PostForm(baseUrl+"auth/login", postParam)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var loginjson RespJson
	err = json.Unmarshal([]byte(body), &loginjson)
	if err != nil {
		log.Fatal(err)
	}
	if loginjson.Ret != 1 {
		fmt.Println("登陆失败，请检查密码配置")
		//return
	}
	fmt.Println("登陆成功")
	return resp.Cookies()

}

func checkin(cookies []*http.Cookie) {
	type RespJson struct {
		Ret int    `json:"ret"`
		Msg string `json:"msg"`
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar:       jar,
		Transport: tr,
	}
	req, _ := http.NewRequest("POST", baseUrl+"user/checkin", nil)
	for _, v := range cookies {
		req.AddCookie(v)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	data, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(data))
}
