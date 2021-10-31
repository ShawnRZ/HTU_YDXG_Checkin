package main

import (
	"HTU_YDXG_Checkin/internal/mail"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
)

var config Config

func (u *User) login() error {
	var err error
	req, err := http.NewRequest("GET", "https://htu.banjimofang.com/student", nil)
	if err != nil {
		return err
	}
	req.AddCookie(&http.Cookie{
		Name: "yxktml", Value: u.Cookie.Yxktml,
	})
	req.AddCookie(&http.Cookie{
		Name: u.Cookie.RememberStudentKey, Value: u.Cookie.RememberStudentValue,
	})
	resp, err := u.Client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	u.CheckinUrl = "https://" + resp.Request.URL.Host + resp.Request.URL.Path + "/profiles/6099"
	if !strings.Contains(string(body), "师大移动学工用户中心") {
		return errors.New("登录失败")
	}
	reg, err := regexp.Compile(`uname:'.*?'`)
	if err != nil {
		return err
	}
	name := reg.Find(body)
	name = name[7 : len(name)-1]
	u.Name = string(name)
	return nil
}

func (u *User) getFormID() error {
	var err error
	req, err := http.NewRequest("GET", u.CheckinUrl, nil)
	if err != nil {
		return err
	}
	req.AddCookie(&http.Cookie{
		Name: "yxktml", Value: u.Cookie.Yxktml,
	})
	req.AddCookie(&http.Cookie{
		Name: u.Cookie.RememberStudentKey, Value: u.Cookie.RememberStudentValue,
	})
	resp, err := u.Client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	reg, err := regexp.Compile(`<input type="hidden" name="form_id" value="\d*?">`)
	if err != nil {
		return err
	}
	formID := reg.Find(body)
	formID = formID[43 : len(formID)-2]
	u.FormID = string(formID)
	return nil
}

func (u User) checkin() error {
	var err error
	postBody := url.Values{}
	for k, v := range u.Postdata {
		postBody.Add(k, v)
	}
	postBody.Add("_bjmf_fields_s", `{"gps":["v"]}`)
	postBody.Add("form_id", u.FormID)
	fmt.Printf("postBody: %#v\n", postBody)
	req, err := http.NewRequest("POST", u.CheckinUrl, strings.NewReader(postBody.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{
		Name: "yxktml", Value: u.Cookie.Yxktml,
	})
	req.AddCookie(&http.Cookie{
		Name: u.Cookie.RememberStudentKey, Value: u.Cookie.RememberStudentValue,
	})
	resp, err := u.Client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if strings.Contains(string(body), `<div class="message">提示</div>`) {
		reg := regexp.MustCompile(`<div class="desc">.*?</div>`)
		msg := reg.Find(body)
		return errors.New(string(msg))
	}

	// fmt.Printf("body: %v\n", string(body))
	return nil
}

func (u *User) init() error {
	var err error
	u.Client = &http.Client{}
	err = u.login()
	if err != nil {
		return err
	}
	err = u.getFormID()
	if err != nil {
		return err
	}
	return nil
}

func start() {
	if _, err := toml.DecodeFile("./configs/config.toml", &config); err != nil {
		fmt.Println(err)
		panic(err)
	}
	mailBodyBytes, err := ioutil.ReadFile("./mail.html")
	mailBody := string(mailBodyBytes)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	mailBody = strings.Replace(mailBody, "%7B%7Bversion%7D%7D", "v0.4", -1)

	for i, v := range config.User {
		var err error
		err = v.init()
		if err != nil {
			fmt.Println("初始化失败：" + err.Error())
			continue
		}
		fmt.Printf("###第 %d 位用户###\n", i+1)
		fmt.Println("用户名：" + v.Name)
		mailBodyTmp := strings.Replace(mailBody, "{{name}}", v.Name, -1)
		now := time.Now()
		mailBodyTmp = strings.Replace(mailBodyTmp, "{{date}}", fmt.Sprintf("%02d-%02d-%02d %02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()), -1)
		err = v.checkin()
		if err != nil {
			fmt.Println("打卡失败：" + err.Error())
			mailBodyTmp = strings.Replace(mailBodyTmp, "{{msg}}", "打卡失败！"+err.Error(), -1)
			err = mail.SendEmail(v.Mail.Username, v.Mail.Username, v.Mail.Password, "HTU移动校工打卡推送", string(mailBodyTmp))
			if err != nil {
				fmt.Println("邮件发送失败", err)
			}
			continue
		}
		fmt.Println("打卡成功")
		mailBodyTmp = strings.Replace(mailBodyTmp, "{{msg}}", "打卡成功！！！", -1)
		err = mail.SendEmail(v.Mail.Username, v.Mail.Username, v.Mail.Password, "HTU移动校工打卡推送", string(mailBodyTmp))
		if err != nil {
			fmt.Println("邮件发送失败", err)
		}
	}
}
func main() {
	cloudfunction.Start(start)
}
