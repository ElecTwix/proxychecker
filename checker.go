package proxychecker

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var globalip string = ""

func checkall(proxies *[]string, site string, timeout time.Duration, InsecureSkipVerify *bool) (wproxies []string, err error) {

	ip, err := GetNormalip()

	if err != nil {
		return nil, err
	}

	globalip = *ip

	wproxies = make([]string, 0)

	var wg sync.WaitGroup
	channel := make(chan string, len(*proxies))

	for _, v := range *proxies {
		wg.Add(1)
		go checkproxy(v, &site, &timeout, channel, &wg, InsecureSkipVerify)
	}

	wg.Wait()

	close(channel)

	for v := range channel {
		wproxies = append(wproxies, v)
	}

	return wproxies, nil
}

func checkproxy(proxy string, site *string, timeout *time.Duration, channel chan string, wg *sync.WaitGroup, insecure *bool) {

	defer wg.Done()
	t := &http.Transport{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   proxy,
		}),
	}

	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: *insecure}

	client := &http.Client{
		Transport: t,
		Timeout:   *timeout,
	}

	resp, err := client.Get(*site)

	if err != nil {
		return
	} else {
		if resp.StatusCode != http.StatusOK {
			return
		}
	}

	if checkip(resp.Body) {
		channel <- proxy
	}

	return
}

func checkip(respbody io.ReadCloser) bool {
	ip, err := getip(respbody)
	if err != nil {
		return false
	}

	if string(ip) != globalip {
		return true
	}
	return false

}

func getip(respbody io.ReadCloser) (string, error) {
	body, err := ioutil.ReadAll(respbody)
	if err != nil {
		return "", err
	}

	respbody.Close()

	return string(body), nil
}

func GetNormalip() (localip *string, err error) {
	resp, err := http.Get("http://ifconfig.me/ip")

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Status Code is not 200")
	}

	ip, err := getip(resp.Body)
	if err != nil {
		return nil, err
	}

	return &ip, nil
}
