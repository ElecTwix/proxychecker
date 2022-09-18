package proxychecker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestCheckAll(T *testing.T) {

	resp, err := http.Get("https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt")

	if err != nil {
		T.Fatal(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)

	strarr := strings.Split(string(body), "\n")

	site := "http://ifconfig.me/ip"
	dur := time.Second * 5
	blo := false
	arr2, err := checkall(&strarr, site, dur, &blo)

	if err != nil {
		T.Fatal(err.Error())
	}
	fmt.Println(len(arr2))
}
