package leancloud

import (
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultSiteURL = "https://leancloud.cn/1.1"

type Config struct {
	AppId       string
	AppKey      string
	MasterKey   string
	UsingMaster bool
	SiteURL     string
}

type Cloud struct {
	Cfg                Config
	HeaderSessionToken string
	HeaderProduction   string
	BeforeRequest      func(*http.Request) *http.Request
}

func (cloud *Cloud) makeSign() string {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	sign := ""
	if cloud.Cfg.UsingMaster {
		sign = fmt.Sprintf("%x", md5.Sum([]byte(timestamp+cloud.Cfg.MasterKey)))
		return fmt.Sprintf("%s,%s,%s", sign, timestamp, "master")
	} else {
		sign = fmt.Sprintf("%x", md5.Sum([]byte(timestamp+cloud.Cfg.AppKey)))
		return fmt.Sprintf("%s,%s", sign, timestamp)
	}
}

func (cloud *Cloud) makeURL(parts ...string) (url string) {
	var path string
	if len(parts) == 0 {
		log.Panicln("can not make url", parts)
	} else if len(parts) == 1 {
		path = parts[0]
	} else {
		path = strings.Join(parts, "/")
	}
	if cloud.Cfg.SiteURL == "" {
		url = fmt.Sprintf("%s/%s", defaultSiteURL, path)
	} else {
		url = fmt.Sprintf("%s/%s", cloud.Cfg.SiteURL, path)
	}
	return
}

func (cloud *Cloud) makeURLPrefix(prefix string, parts ...string) string {
	tmp := make([]string, 0)
	tmp = append(tmp, prefix)
	tmp = append(tmp, parts...)
	return cloud.makeURL(tmp...)
}

func (cloud *Cloud) DoRequest(url, method, body string) (*Result, error) {
	r, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Avoscloud-Application-Id", cloud.Cfg.AppId)
	r.Header.Set("X-AVOSCloud-Request-Sign", cloud.makeSign())
	if cloud.HeaderSessionToken != "" {
		r.Header.Set("X-AVOSCloud-Session-Token", cloud.HeaderSessionToken)
	}
	if cloud.HeaderProduction != "" {
		r.Header.Set("X-AVOSCloud-Application-Production", cloud.HeaderProduction)
	}
	if cloud.BeforeRequest != nil {
		r = cloud.BeforeRequest(r)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	sbody, _ := ioutil.ReadAll(res.Body)
	var location string
	if u, err := res.Location(); err == nil {
		location = u.String()
	}
	ret := &Result{res.StatusCode, location, string(sbody)}
	if !ret.CheckStatusCode() {
		return ret, errors.New(ret.Response)
	}
	return ret, nil
}

func (cloud *Cloud) Get(url string, param url.Values) (*Result, error) {
	withQuery := fmt.Sprintf("%s?%s", url, param.Encode())
	return cloud.DoRequest(withQuery, "GET", "")
}

func (cloud *Cloud) Put(url, body string) (*Result, error) {
	return cloud.DoRequest(url, "PUT", body)
}

func (cloud *Cloud) Delete(url string) (*Result, error) {
	return cloud.DoRequest(url, "DELETE", "")
}

func (cloud *Cloud) Post(url, body string) (*Result, error) {
	return cloud.DoRequest(url, "POST", body)
}
