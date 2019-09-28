// crawler.go
package main

import (
	//"crypto/tls"
	"fmt"
	"io/ioutil"

	//"net/http"
	"os"
	"strings"

	"github.com/bclicn/color"
	"github.com/headzoo/surf/agent"
	"gopkg.in/headzoo/surf.v1"
	"menteslibres.net/gosexy/to"
)

func read_url_list() []string {
	var urls_greped []string
	fmt.Println("Entering crawler mode, reading URLs from urls file..")
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

	read_content()
	//make_http_requests(urls_greped)
	return urls_greped
}

func read_content() string {
	body, err := ioutil.ReadFile("bodys")
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "content file not found. ", err)
		os.Exit(1)
	}
	decode_premproxy(to.String(body))
	return to.String(body)
}

func decode_premproxy(content string) {

	//cut the proxy table
	table_open := strings.Index(content, "<tbody>")
	table_close := strings.Index(content, "</tbody>")
	table := content[table_open:table_close]
	//split it per line
	table_splitted := strings.Split(table, "\n")

	i := 0
	proxys := make(map[string][]string, len(table_splitted))

	for i = 0; i < len(table_splitted); i++ {
		if strings.Index(table_splitted[i], "value=\"") > 0 {
			ip_index := strings.Index(table_splitted[i], "value=\"")
			port_index := strings.Index(table_splitted[i], "|")
			if port_index == -1 {
				break
			}

			proxys["ip"] = append(proxys["ip"], table_splitted[i][ip_index+7:port_index])
			proxys["port_encrypted"] = append(proxys["port_encrypted"], table_splitted[i][port_index+1:port_index+6])
			//fmt.Println(proxys["ip"][i], ":", proxys["port_encrypted"][i])
		}
	}
	fmt.Println(color.BWhite(to.String(i-1)), "IP's found.")

	encoding := make(map[string][]string)
	encoding["codes"] = append(encoding["codes"], "r1dff")
	encoding["port"] = append(encoding["port"], "4145")
	encoding["codes"] = append(encoding["codes"], "r1ba7")
	encoding["port"] = append(encoding["port"], "4153")
	encoding["codes"] = append(encoding["codes"], "r55ac")
	encoding["port"] = append(encoding["port"], "8291")
	encoding["codes"] = append(encoding["codes"], "r336e")
	encoding["port"] = append(encoding["port"], "9999")
	encoding["codes"] = append(encoding["codes"], "r926a")
	encoding["port"] = append(encoding["port"], "1080")
	encoding["codes"] = append(encoding["codes"], "rf1e7")
	encoding["port"] = append(encoding["port"], "8888")
	encoding["codes"] = append(encoding["codes"], "r8e22")
	encoding["port"] = append(encoding["port"], "9201")
	encoding["codes"] = append(encoding["codes"], "r7e60")
	encoding["port"] = append(encoding["port"], "3629")
	encoding["codes"] = append(encoding["codes"], "rf4b6")
	encoding["port"] = append(encoding["port"], "6363")
	encoding["codes"] = append(encoding["codes"], "r80f5")
	encoding["port"] = append(encoding["port"], "8082")
	encoding["codes"] = append(encoding["codes"], "ra789")
	encoding["port"] = append(encoding["port"], "9249")
	encoding["codes"] = append(encoding["codes"], "r17e8")
	encoding["port"] = append(encoding["port"], "5893")
	encoding["codes"] = append(encoding["codes"], "r0544")
	encoding["port"] = append(encoding["port"], "1157")
	encoding["codes"] = append(encoding["codes"], "r8324")
	encoding["port"] = append(encoding["port"], "8010")

	j := 0
	k := 0

	for j = 0; j < i-1; j++ {
		for k = 0; k < len(encoding["codes"]); k++ {
			//fmt.Println("probe encoding: Port:", j, proxys["port_encrypted"][j], "code:", k, encoding["codes"][k])
			if strings.Contains(proxys["port_encrypted"][j], encoding["codes"][k]) == true {
				proxys["port_decrypted"] = append(proxys["port_decrypted"], encoding["port"][k])
				break
			} else if len(encoding["codes"]) == k {
				fmt.Println(color.BWhite(to.String(proxys["port_encrypted"][j])), "Not found.")
				proxys["port_decrypted"] = append(proxys["port_decrypted"], "Not found.")
			}

		}
		k = 0
		fmt.Println(proxys["ip"][j], ":", proxys["port_decrypted"][j])
	}

	if len(proxys["port_encrypted"]) == len(proxys["port_decrypted"]) {

		fmt.Println("All", color.BWhite(to.String(i-1)), "codes found.")
	}

}

func make_http_requests(urls []string) {
	var url_content []string
	//TLS handling
	/*tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	for i := 0; i < len(urls); i++ {
		req, _ := http.NewRequest("GET", urls[i], nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36")
		resp, err := client.Do(req) //do the request
		if err != nil {
			gimme_err := fmt.Sprint("Can't fetch data.", err)
			url_content = append(url_content, gimme_err)
		} else {
			//defer resp.Body.Close()
			responseData, _ := ioutil.ReadAll(resp.Body) // fetch body
			url_content = append(url_content, to.String(responseData))
		}
	}*/

	bow := surf.NewBrowser()
	bow.SetUserAgent(agent.Chrome())

	for i := 0; i < len(urls); i++ {

		err := bow.Open(urls[i])
		if err != nil {
			gimme_err := fmt.Sprint("Can't fetch data.", err)
			url_content = append(url_content, gimme_err)
		} else {
			url_content = append(url_content, to.String(bow.Body()))

		}
	}
	fmt.Println(url_content)
	fmt.Println(color.BWhite(to.String(len(url_content))), "url(s) crawled.")

}
