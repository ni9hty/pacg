// crawler.go
package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/bclicn/color"
	"menteslibres.net/gosexy/to"
)

func read_url_list() []string {
	var urls_greped []string
	fmt.Println("Crawler mode activated, reading URLs from urls file..")
	urls, err := ioutil.ReadFile("urls")
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "urls file not found. ", err)
		os.Exit(1)
	} else {
		urls_listed := strings.Split(to.String(urls), "\n")
		count := 0
		for i := 0; i < len(urls_listed); i++ {
			if strings.Contains(urls_listed[i], "http") {
				count++
				urls_greped = append(urls_greped, urls_listed[i])
			}
		}
		fmt.Println(color.BWhite(to.String(count)), "url(s) found.")
	}
	make_http_requests(urls_greped)
	return urls_greped
}

func make_http_requests(urls []string) {
	var url_content []string
	//TLS handling
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	for i := 0; i < len(urls); i++ {
		req, _ := http.NewRequest("GET", urls[i], nil)
		resp, err := client.Do(req) //do the request
		if err != nil {
			gimme_err := fmt.Sprint("Can't fetch data.", err)
			url_content = append(url_content, gimme_err)
		} else {
			//defer resp.Body.Close()
			responseData, _ := ioutil.ReadAll(resp.Body) // fetch body
			url_content = append(url_content, to.String(responseData))
		}
	}
	//fmt.Println(url_content[1])

}
