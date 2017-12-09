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
	"github.com/sparrc/go-ping"
	"menteslibres.net/gosexy/to"
)

func check_enviroment() string {
	fmt.Println("Checking enviroment ..")
	//proxychain executeable
	var output_pc string

	cmd_check_pc := exec.Command("locate", "proxychain")
	output_check_pc, err := cmd_check_pc.CombinedOutput()
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "locate command failed: ", err)
	}
	if strings.Contains(to.String(output_check_pc), "proxychains4") == true {
		find_string := strings.Index(to.String(output_check_pc), "/usr/local/bin/proxychains4")
		output_pc = to.String(output_check_pc)[find_string : find_string+27]
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
	User := user.Username

	user_perm_splitted := strings.Split(to.String(output_check_ls), " ")
	config_user_owner_perm := user_perm_splitted[3]
	config_user_write_perm := user_perm_splitted[0]

	if strings.Contains(to.String(User), config_user_owner_perm) {
		fmt.Println(color.LightGreen("[+] "), "current user", User, "is the owner of /etc/proxychains.conf")
		if strings.HasPrefix(config_user_write_perm, "-rw") {
			fmt.Println(color.LightGreen("[+] "), "current user", User, "can edit the config.")
		} else {
			fmt.Println(color.LightRed("[-] "), "current user", User, "can NOT edit the config.")
			fmt.Println("Please adjust file permissions : chmod +w", output_pc)
			os.Exit(1)
		}
	} else {
		fmt.Println(color.LightRed("[-] "), "current user", User, "is not the owner.")
		fmt.Println("Please adjust file permissions : chown", User, ":", User, output_pc)
	}

	sys_ping := fmt.Sprint("/proc/sys/net/ipv4/ping_group_range")
	content, err := ioutil.ReadFile(sys_ping)
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "Icmp ping group config file not found. Check proxy(s) via ping might be only possible via sudo. ", err)
	}
	if strings.Contains(to.String(content), "0\t2147483647") {
		fmt.Println(color.LightGreen("[+] "), "Icmp ping group config file have the correct settings.")
	} else {
		fmt.Println(color.LightRed("[-] "), "Icmp ping group settings are wrong, ping are only possible via sudo.\nPlease adjust with: sudo sysctl -w net.ipv4.ping_group_range=\"0   2147483647\"")
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
	fmt.Println(proxys)
	//check_proxys(proxys)
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

func check_proxys(proxy_map map[string][]string) {

	results := make(map[string][]string, len(proxy_map["ip"]))
	//next_try := make(map[string][]string)

	//ping them
	for _, ip := range proxy_map["ip"] {
		p, err := ping.NewPinger(ip)
		if err != nil {
			fmt.Println("couldn't ping: ", err)
		}
		p.Count = 1
		p.SetPrivileged(false)
		p.Timeout = time.Second * 2
		p.OnRecv = func(pkt *ping.Packet) {
			//fmt.Println(pkt.Nbytes, "bytes from", pkt.IPAddr, "time=", pkt.Rtt)
			if pkt.Rtt > time.Millisecond*200 {
				fmt.Println(color.GGreen("[-] "), ip, " latency > 200ms")
			}
			results["time"] = append(results["time"], to.String(pkt.Rtt))
		}
		p.Run()
	}

	//check if port open
	i := 0
	for _, _ = range proxy_map["ip"] {
		con_string := fmt.Sprint(proxy_map["ip"][i], ":", proxy_map["port"][i])
		_, err := net.DialTimeout("tcp", con_string, time.Second*5)
		if err != nil {
			fmt.Println("Proxy not available ", err)
			//here implement new try until len(proxy_map) ends
		} else {
			//fmt.Println(con_string, " is open")
			results["con_string"] = append(results["con_string"], con_string)
		}
		i++
	}

	for j := 0; j < len(results["con_string"]); j++ {
		fmt.Println(results["con_string"][j], "open, time=", results["time"][j])
	}

}

func generate_config(proxychain string) {
	content, err := ioutil.ReadFile("/etc/proxychains.conf")
	if err != nil {
		fmt.Println(color.LightRed("[-] "), "Could not read config file ", err)
	}
	lastIndex_raute := strings.LastIndex(to.String(content), "#")
	fmt.Println("Raute: ", lastIndex_raute)

}

func main() {

	howmuch := flag.Int("n", 2, "how much proxys do you want to use")
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

	pc_exec := check_enviroment()
	proxys_map = gimmeproxy(*howmuch)
	check_proxys(proxys_map)
	generate_config(pc_exec)

}
