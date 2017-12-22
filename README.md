# Pacg - The Proxy(chain) Auto Config Generator

Related to the [Proxychain](https://github.com/rofl0r/proxychains-ng) package, i was too lazy to search for valid proxys and put them into my config. It's a golang script which check several conditions to be ready to work, grep some proxys (2 by default) from [gimmeproxy API](https://gimmeproxy.com/) and write a custom proxychain config. 

Commandline parameter will be increased.
```bash
./pacg --help             
Usage of ./pacg:
  -n int
    	how much proxys do you want to use (default 2)
```
Default gimmeproxy curl = 
```bash
https://gimmeproxy.com/api/getProxy?anonymityLevel=1&?protocol=socks5&maxCheckPeriod=300
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
[-]  121.52.141.104  latency > 200ms
77.244.42.***:43524 open, time= 64.135386ms in FR
121.52.141.***:8080 open, time= 283.265037ms in BR
[+]  /etc/proxychains.conf successfully written.
```

### already Done
- [x] Parameter how many jumps in your chain do you want
- [x] Check local permission conditions
- [x] if a request result the same proxy ip already checked via the api, result in a new check 
- [x] implemented [goeip](https://github.com/rainycape/geoip) function for generating smart routes from your start country (db already in this repo and home/start ip already catched)
- [ ] Proxy availability + latency check, if proxy timeout or to slow, get a new one (both conditions already catched)

### Brainstorm & notes
- [ ] Proxy availability + latency check, if a proxy timeout occur, get a new one (both conditions already catched)
- [ ] implement parameter through how many countrys you want to chain, with latency checks (min/max) condition (if not fast enough, get a new one)
- [ ] implement other free proxy sources (like modules) to be more flexible

Feel free to comment or PR if you like.
