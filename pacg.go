// test.go
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/bclicn/color"
	"github.com/bitly/go-simplejson"
	"github.com/rainycape/geoip"
	"menteslibres.net/gosexy/to"
)

func check_enviroment() string {
	fmt.Println("Checking enviroment ..")
	//proxychain executeable
	var output_pc string
	var output_pc1 string

	cmd_check_pc := exec.Command("locate", "proxychain")
	output_check_pc, err := cmd_check_pc.CombinedOutput()
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "locate command failed: ", err)
	}
	if strings.Contains(to.String(output_check_pc), "proxychains4") == true {
		find_string := strings.Index(to.String(output_check_pc), "/usr/bin/proxychains4")
		output_pc = to.String(output_check_pc)[find_string : find_string+21]
		fmt.Println(color.LightGreen("[+] "), output_pc, "[found]")
	} else {
		fmt.Println(color.LightRed("[-] "), "No proxychain executeable found, unable to autoupdate proxy config entries.\nPlease follow instructions on https://github.com/rofl0r/proxychains-ng")
	}

	//check file permissions
	var _, err1 = os.Stat("/etc/proxychains.conf")
	if os.IsNotExist(err1) {
		fmt.Println(color.LightRed("[-] "), "/etc/proxychains.conf not present, please create one with write permissions for the current user.")
		os.Exit(1)
	}
	cmd_check_ls := exec.Command("ls", "-la", "/etc/proxychains.conf")
	output_check_ls, err := cmd_check_ls.CombinedOutput()
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "ls command failed: ", err)
		os.Exit(1)
	}
	user, err := user.Current()
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "Can't get current username. ", err)
		os.Exit(1)
	}
	U := user.Username
	User := strings.Replace(U, " ", "", -1)

	user_perm_splitted := strings.Split(to.String(output_check_ls), " ")
	config_user_owner_perm := user_perm_splitted[3]
	config_user_write_perm := user_perm_splitted[0]
	find_string1 := strings.Index(to.String(output_check_pc), "/etc/proxychains.conf")
	output_pc1 = to.String(output_check_pc)[find_string1 : find_string1+21]

	if strings.Contains(to.String(User), config_user_owner_perm) {
		fmt.Println(color.LightGreen("[+] "), "current user", User, "is the owner of /etc/proxychains.conf")
		if strings.HasPrefix(config_user_write_perm, "-rw") {
			fmt.Println(color.LightGreen("[+] "), "current user", User, "can edit the config.")
		} else {
			fmt.Println(color.LightRed("[-] "), "current user", User, "can NOT edit the config.")
			fmt.Println("Please adjust file permissions : chmod +w", output_pc1)
			os.Exit(1)
		}
	} else {
		fmt.Println(color.LightRed("[-] "), "current user", User, "is not the owner.")
		chown := fmt.Sprint(User, ":", User)
		fmt.Println("Please adjust file permissions : chown", chown, output_pc1)
	}

	myip := myip()
	country := geoip_request(myip)
	fmt.Println(color.LightGreen("[+] "), "We are starting from", myip, "in", country)

	return output_pc
}

func geoip_request(ip string) string {
	//geoipdb := "GeoIP.dat"
	db, err := geoip.Open("GeoLite2-Country.mmdb")
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "Couldn't open GeoLite2-Country.mmdb file. ", err)
	}

	res, err := db.Lookup(ip)
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "Couldn't lookup ", ip, " in geoip db ", err)
	}
	return to.String(res.Country.Name)
}

func myip() string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", "http://ipinfo.io/ip", nil)
	resp, err := client.Do(req) //do the request
	if err != nil {
		myip_err := fmt.Sprint("Can't fetch myip data.", err)
		fmt.Println(color.LightRed("[-] "), myip_err)
	}
	//defer resp.Body.Close()
	responseData, _ := ioutil.ReadAll(resp.Body) // fetch body
	resp_string := strings.TrimSuffix(to.String(responseData), "\n")
	return resp_string
}

func gimmeproxy(count int) map[string][]string {
	//runtime map of all fetched proxys
	proxys := make(map[string][]string)
	out := simplejson.New()
	//if filtered json file exists, truncate the content
	var _, err1 = os.Stat("temp_proxy.json")
	if !os.IsNotExist(err1) {
		e, err := os.OpenFile("temp_proxy.json", os.O_RDWR, 644)
		if err != nil {
			fmt.Println(color.LightRed("[-] "), "Unable to truncate file temp_proxy.json: ", err)
		}
		e.Truncate(0)
		e.Close()
	}
	i := 0
request:
	//TLS handling
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, _ := http.NewRequest("GET", "https://gimmeproxy.com/api/getProxy?anonymityLevel=1&?protocol=socks5&maxCheckPeriod=300", nil)

	resp, err := client.Do(req) //do the request
	if err != nil {
		gimme_err := fmt.Sprint("Can't fetch data.", err)
		fmt.Println(gimme_err)
	}
	//defer resp.Body.Close()
	responseData, _ := ioutil.ReadAll(resp.Body) // fetch body
	//fmt.Println(to.String(responseData))

	js, err := simplejson.NewJson(responseData)
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "Unable to fetch json: ", err)
		fmt.Println("\n", to.String(resp.Body))
		os.Exit(1)
	}

	for i < count {

		ip := js.Get("ip").MustString()
		for _, no_double := range proxys {
			if strings.Contains(to.String(no_double), to.String(ip)) {
				fmt.Println(color.GGreen("[-]"), "got no new IP, next try in 2sec.", no_double, ip)
				time.Sleep(time.Second * 2)
				goto request
			}
		}
		out.Set("ip", ip)
		proxys["ip"] = append(proxys["ip"], to.String(ip))

		country := js.Get("country").MustString()
		out.Set("country", country)
		proxys["country"] = append(proxys["country"], to.String(country))

		port := js.Get("port").MustString()
		out.Set("port", port)
		proxys["port"] = append(proxys["port"], to.String(port))

		protocol := js.Get("protocol").MustString()
		out.Set("protocol", protocol)
		proxys["protocol"] = append(proxys["protocol"], to.String(protocol))

		i++
		create_filtered_json(out)
		goto request
	}
	//fmt.Println(proxys)
	return proxys
}

//for possible later usage
func create_filtered_json(output *simplejson.Json) {
	b, _ := output.EncodePretty()
	var _, err = os.Stat("tmp_proxys.json")

	if os.IsNotExist(err) {
		var file, err = os.Create("tmp_proxys.json")
		if err != nil {
			fmt.Println(color.LightRed("[-] "), "Unable to create file tmp_proxys.json: ", err)
			defer file.Close()
		}
	}
	f, err := os.OpenFile("tmp_proxys.json", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "Unable to to write tmp_proxys.json, ", err)
	}
	_, err = f.Write(b)
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "Unable to to write tmp_proxys.json, ", err)
	}
	f.Close()
}

func check_proxys(proxy_map map[string][]string, proxy_count int) map[string][]string {
	fmt.Println("Checking", len(proxy_map["ip"]), "proxy(s) ..")
	temp := make(map[string][]string, 1)
	//results := make(map[string][]string, len(proxy_map["ip"]))
	results := make(map[string][]string, proxy_count)
	nexttry := false

next_try:
	//check if port open
	i := 0

	country := ""
	for _, _ = range proxy_map["ip"] {
		country = ""
		con_string := fmt.Sprint(proxy_map["ip"][i], ":", proxy_map["port"][i])
		start := time.Now()
		_, err := net.DialTimeout("tcp", con_string, time.Second*5)
		ms := time.Since(start)
		if err != nil {
			fmt.Println(color.LightRed("[-] "), "Proxy not available ", err, "try next ..")
			nexttry = true
			if i == 0 {
				delete(results, "time") //if first proxy unavailable, delete time entry
			} else {
				del_last_one := len(results["time"]) - 1
				results["time"] = results["time"][:del_last_one]
			}

		} else {
			results["con_string"] = append(results["con_string"], con_string)
			results["ip"] = append(results["ip"], proxy_map["ip"][i])
			results["port"] = append(results["port"], proxy_map["port"][i])
			results["protocol"] = append(results["protocol"], proxy_map["protocol"][i])
			results["tld"] = append(results["tld"], proxy_map["country"][i])
			country = geoip_request(proxy_map["ip"][i])
			results["country"] = append(results["country"], country)
			if len(results["ip"]) > 0 {
				//fmt.Println("tcp ping test: ", results["ip"][i], " ", ms)
				results["tcp_ping"] = append(results["tcp_ping"], to.String(ms))
			}

		}
		i++
	}

	if len(results["con_string"]) < proxy_count && nexttry == true {
		delete(proxy_map, "ip")
		delete(proxy_map, "port")
		delete(proxy_map, "country")
		delete(proxy_map, "protocol")
		temp = make(map[string][]string, 1)
		temp = gimmeproxy(1)
		proxy_map["ip"] = append(proxy_map["ip"], temp["ip"][0])
		proxy_map["port"] = append(proxy_map["port"], temp["port"][0])
		proxy_map["country"] = append(proxy_map["country"], temp["country"][0])
		proxy_map["protocol"] = append(proxy_map["protocol"], temp["protocol"][0])
		goto next_try
	}

	for j := 0; j < len(results["con_string"]); j++ {
		fmt.Println("[", j+1, "]", results["con_string"][j], "open, time_tcp=", results["tcp_ping"][j], "in", results["tld"][j], "-", results["country"][j])
	}

	fmt.Print("\n", color.LightGreen("[+] "), "Found ", proxy_count, " Proxys.\n")
	return results
}

func generate_config(rein map[string][]string, quiet_mode bool, proxy_dns bool) {

	/*for j := 0; j < len(rein["ip"])-1; j++ {
		//fmt.Println(rein["con_string"][j], "open, time=", rein["time"][j], "in", rein["tld"][j], "-", rein["country"][j])
		fmt.Println(rein)
	}*/

	err := os.Truncate("/etc/proxychains.conf", 0)
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "Could not reset config file ", err)
	}

	//dynamic quiet mode setting
	var q string
	if quiet_mode == true {
		q = "quiet_mode"
	} else {
		q = "#quiet_mode"
	}

	//dynamic proxy_dns setting
	var dns string
	if proxy_dns == true {
		dns = "proxy_dns"
	} else {
		dns = "#proxy_dns"
	}

	default_conf_string := fmt.Sprint("strict_chain\n# Some timeouts in ms\ntcp_read_time_out 15000\ntcp_connect_time_out 8000\n", q, "\n", dns, "\n[ProxyList]\n")

	var proxys string
	for i := 0; i < len(rein["ip"]); i++ {
		proxys = fmt.Sprint(rein["protocol"][i], " ", rein["ip"][i], " ", rein["port"][i], "\n")
		default_conf_string += proxys
	}
	default_conf := []byte(default_conf_string)

	write_err_default_pacg_conf := ioutil.WriteFile("/etc/proxychains.conf", default_conf, 0644)
	if write_err_default_pacg_conf != nil {
		fmt.Println("Write /etc/proxychains.conf file ERROR: ", write_err_default_pacg_conf)
		fmt.Println("Pls run command as sudo !!")
		os.Exit(1)
	} else {
		fmt.Println(color.LightGreen("[+] "), "/etc/proxychains.conf successfully written.")
	}
}

func main() {

	howmuch := flag.Int("n", 2, "how much proxys do you want to use")
	quiet := flag.Bool("q", false, "Generate config with quiet mode setting (no output from the library) - (default false)")
	dns := flag.Bool("dns", false, "Generate config with proxy dns option, no leak for DNS data - (default false)")
	crawler := flag.Bool("crawler", false, "switch to crawler mode to generate own proxydb from urls file")
	gproxy := flag.Bool("gimmeproxy", false, "fetch gimmeproxy.com for proxys, check them and write a proxychain config")
	flag.Parse()

	fmt.Println("ProxyChain auto config generator.")
	proxys_map := make(map[string][]string, *howmuch)
	//if filtered json file exists, truncate the content
	var _, err1 = os.Stat("tmp_proxys.json")
	if !os.IsNotExist(err1) {
		e, err := os.OpenFile("tmp_proxys.json", os.O_RDWR, 644)
		if err != nil {
			fmt.Println(color.LightRed("[-] "), "Unable to truncate file tmp_proxys.json: ", err)
		}
		e.Truncate(0)
		e.Close()
	}

	if *crawler {
		read_url_list()
		os.Exit(0)
	}

	if *gproxy {
		check_enviroment()
		proxys_map = gimmeproxy(*howmuch)
		checked_proxys := check_proxys(proxys_map, *howmuch)
		generate_config(checked_proxys, *quiet, *dns)
	}

	if *crawler == false && *gproxy == false {
		flag.Usage()
		os.Exit(0)
	}

}
