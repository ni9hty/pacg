# Pacg - The Proxy(chain) Auto Config Generator

Related to the [Proxychain](https://github.com/rofl0r/proxychains-ng) package, i was too lazy to search for valid proxys and put them into my config. It's a golang script which check several conditions to be ready to work, grep proxys (2 by default) from [gimmeproxy API](https://gimmeproxy.com/) and write a custom proxychain config. 

Proxychains available config settings *dynamic* (order like in the chain but dead proxys skipped) and *random* (in combo with chain_len setting) are obsolete.
So "strict_chain" is what we want, our own working, individual chain.

```bash
./pacg --help            
Usage of ./pacg:
  -crawler
    	switch to crawler mode to generate own proxydb from urls file
  -dns
    	Generate config with proxy dns option, no leak for DNS data - (default false)
  -n int
    	how much proxys do you want to use (default 2)
  -q	Generate config with quiet mode setting (no output from the library) - (default false)
```
 
Example:
```bash
./pacg       
ProxyChain auto config generator.
Checking enviroment ..
[+]  /usr/local/bin/proxychains4 [found]
[+]  current user * is the owner of /etc/proxychains.conf
[+]  current user * can edit the config.
[+]  Icmp ping group config file have the correct settings.
[+]  We are starting from xx.xxx.xxx.xx in Germany
Checking 2 proxy(s) ..
[-]  121.52.141.104  latency > 200ms
77.244.42.***:43524 open, time= 64.135386ms in FR - France
121.52.141.***:8080 open, time= 283.265037ms in BR - Brazil
[+]  /etc/proxychains.conf successfully written.
```

### already Done
- [x] Parameter how many jumps in your chain do you want
- [x] Check local permission conditions
- [x] if a request result the same proxy ip already checked via the api, result in a new check 
- [x] implemented [geoip](https://github.com/rainycape/geoip) function for generating smart routes from your start country (db already in this repo and home/start ip already catched)
- [x] implement other free proxy sources with a crawler mode, slapword html-tables_to_json updateable regular/manually from an extra file with links (started)
- [x] Proxy availability, if a proxy timeout occur, get a new one
- [x] Proxy latency check via tcp, if icmp is blocked

### Brainstorm & notes
- [ ] implement parameter through how many countrys you want to chain
- [ ] implement function to generate smart and fast routes through your target via countrys in an order with lowest latency (need some conceptual planing, start country already catched)
- [ ] implement a function to re-check already generated config and by whish replace single proxys, like re-check availability/latency for currently inserted proxys

Feel free to comment or PR if you like.
