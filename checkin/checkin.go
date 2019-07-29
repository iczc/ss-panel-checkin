package checkin

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

type response struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

type Checkin struct {
	host   string
	email  string
	passwd string
}

func NewCheckin(host string, email string, passwd string) *Checkin {
	return &Checkin{
		host:   host,
		email:  email,
		passwd: passwd,
	}
}

func (c *Checkin) login() ([]*http.Cookie, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	postParam := url.Values{
		"email":  {c.email},
		"passwd": {c.passwd},
	}
	resp, err := client.PostForm(c.host+"/auth/login", postParam)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var loginjson response
	// loginjson := &response{}
	if err = json.Unmarshal([]byte(body), &loginjson); err != nil {
		return nil, err
	}
	if loginjson.Ret != 1 {
		return nil, fmt.Errorf("登陆失败，请检查密码配置")
	}
	log.Println("登陆成功")
	return resp.Cookies(), nil

}

func (c *Checkin) Handle() error {
	cookies, err := c.login()
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	client := http.Client{
		Jar:       jar,
		Transport: tr,
	}
	req, _ := http.NewRequest("POST", c.host+"/user/checkin", nil)
	for _, v := range cookies {
		req.AddCookie(v)
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Println(string(data))
	return nil
}
