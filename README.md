# Pacg - The Proxy Chain Auto Config Generator

Related to the [Proxychain](https://github.com/rofl0r/proxychains-ng) package, i was too lazy to search for valid proxys and put them into my config. 

It's not done, just need a version control.

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
map[ip:[77.244.42.55 121.52.141.104] country:[UA ID] port:[8080 8080] protocol:[http http]]
[-]  121.52.141.104  latency > 200ms
77.244.42.55:8080 open, time= 64.135386ms
121.52.141.104:8080 open, time= 283.265037ms
```
