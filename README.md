XMemcache
---------
XMemcache is a load-balancer for memcache , was wrote in golang.

Installation
------------
```sh
 git clone https://github.com/Maslow/xmemcache.git
```

Introduction
------------
- config  , configuration of servers' ip / port 
- node , manage the server nodes (doctor, consistent-hash alg)
- hash , provide hash value for consistent-hash alg
- protocal, the TCP based binary protocal packet of memcache 
- main.go , start a server & run
