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

func NewCheckin(host, email, passwd string) *Checkin {
	if host == "" || email == "" || passwd == "" {
		log.Fatal("Missing argument")
	}

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
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Login API %s", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	loginResp := &response{}
	if err = json.Unmarshal(data, loginResp); err != nil {
		return nil, err
	}
	if loginResp.Ret != 1 {
		return nil, fmt.Errorf("Login failed")
	}

	log.Println("Login successful")
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
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	checkinResp := &response{}
	if err = json.Unmarshal(data, checkinResp); err != nil {
		return err
	}

	log.Println(checkinResp.Msg)
	return nil
}
