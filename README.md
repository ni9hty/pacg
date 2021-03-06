# Pacg - The Proxy(chain) Auto Config Generator

Related to the [Proxychain](https://github.com/rofl0r/proxychains-ng) package, i was too lazy to search for valid proxys and put them into my config. It's a golang script which check several conditions to be ready to work, grep proxys (2 by default) from [gimmeproxy API](https://gimmeproxy.com/) and write a custom proxychain config. 

Proxychains available config settings *dynamic* (order like in the chain but dead proxys skipped) and *random* (in combo with chain_len setting) are obsolete.
So "strict_chain" is what we want, our own working, individual chain.

```bash
./pacg
ProxyChain auto config generator.
Usage of ./pacg:
  -crawler
    	switch to crawler mode to generate own proxydb from urls file
  -dns
    	Generate config with proxy dns option, no leak for DNS data - (default false)
  -gimmeproxy
    	fetch gimmeproxy.com for proxys, check them and write a proxychain config
  -n int
    	how much proxys do you want to use (default 2)
  -q	Generate config with quiet mode setting (no output from the library) - (default false)
```
 
Example:
```bash
./pacg -gimmeproxy
ProxyChain auto config generator.
Checking enviroment ..
[+]  /usr/bin/proxychains4 [found]
[+]  current user gizn is the owner of /etc/proxychains.conf
[+]  current user gizn can edit the config.
[+]  We are starting from xx.xx.xx.xx in Germany
Checking 2 proxy(s) ..
[ 1 ] xxx.xxx.120.177:51822 open, time_tcp= 119.232756ms in RU - Russia
[ 2 ] xxx.xxx.108.58:41916 open, time_tcp= 72.95672ms in RU - Russia

[+] Found 2 Proxys.
[+]  /etc/proxychains.conf successfully written.
```

### already Done
- [x] Parameter how many jumps in your chain do you want
- [x] Check local permission conditions
- [x] no double ip checks
- [x] implemented [geoip](https://github.com/rainycape/geoip) function for generating smart routes from your start country (db already in this repo and home/start ip already catched)
- [x] Proxy availability, if a proxy timeout occur, get a new one
- [x] Proxy latency check via tcp (icmp option removed)

### Brainstorm & notes
- [ ] implement parameter through how many countrys you want to chain
- [ ] implement function to generate smart and fast routes through your target via countrys in an order with lowest latency 
- [ ] implement a function to re-check already generated config and by whish replace single proxys, like re-check availability/latency for currently inserted proxys
- [ ] In combination with the crawler, create db of suitable proxys as a service like Tor

Feel free to comment or PR if you like.
