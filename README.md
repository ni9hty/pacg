# Pacg - The Proxy(chain) Auto Config Generator

Related to the [Proxychain](https://github.com/rofl0r/proxychains-ng) package, i was too lazy to search for valid proxys and put them into my config. It's a golang script which check several conditions to be ready to work, grep some proxys from [gimmeproxy API](https://gimmeproxy.com/) and write a custom proxychain config. 

It's far away from to be done, just need a version control.

Commandline parameter will be increased.
```bash
./pacg --help             
Usage of ./pacg:
  -n int
    	how much proxys do you want to use (default 2)
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
map[ip:[77.244.42.** 121.52.141.***] country:[UA ID] port:[8080 8080] protocol:[http http]]
[-]  121.52.141.104  latency > 200ms
77.244.42.**:8080 open, time= 64.135386ms
121.52.141.***:8080 open, time= 283.265037ms
```

### Brainstorm & notes
- [x] Parameter how many jumps in your chain do you want
- [x] Check local permission conditions
- [x] if two request result the same proxy ip via the api, result in a new check 
- [x] implemented [goeip](https://github.com/rainycape/geoip) function for generating smart routes from your start country (db already in this repo and home/start ip already catched)
- [-] Proxy availability + latency check, if proxy timeout, get a new one (timeout already catched)
- [ ] implement parameter through how many countrys you want to chain, with latency checks (min/max) condition (if not fast enough, get a new one)
- [ ] implement other free proxy sources (like modules) to be more flexible

Feel free to comment or PR if you like.
