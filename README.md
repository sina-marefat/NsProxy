# **DNS Proxy**

## What is The DNSProxy?

A DNS proxy, also known as a DNS forwarder or a DNS resolver, is an intermediary server that sits between client devices and DNS servers. Its primary function is to handle DNS queries and forward them to appropriate DNS servers for resolution.

When a client device sends a DNS query to a DNS proxy, the proxy examines the request and acts as a middleman between the client and the DNS server. The DNS proxy checks its cache to see if it already has the resolved IP address for the requested domain. If the information is available in its cache, the proxy can respond directly to the client without sending a query to an external DNS server. This caching mechanism helps to improve the response time and reduce the load on DNS servers.



## ***Introduction***

"A UDP connection is established on port `53`. Whenever a request is made on this connection, the data is buffered, and then passed to the parser module. The parser module handles DNS requests, which are identified by a 12-byte header. Using this module, we separate the nameserver from the query, and check the query type (e.g., AAAA, A, etc.). If the query type is incorrect, a DNS error is returned. Otherwise, we check all the keys using the name server. If a key with the specified name server exists and is already cached, we retrieve it. Otherwise, a UDP connection is established with the DNS server. We send the request and receive the response, which is then stored in the UDP connection proxy to be displayed in the command line (cmd)."

![DNS Resolution](https://github.com/sina-marefat/NsProxy/blob/main/images/dns.png)

"As it is clear from the diagram, with the increase in the percentage of proxy use, the utilization time ratio decreases."

> The traffic volume saved by this proxy, considering a daily count of 2000 user request, is maximum of
> `(2000 * 0.5)* 100 /1024 = 97.65 KB`.



## ***Requirment***

To run this project, you will need:

- Golang
- Redis



## ***Installation***

1. Using Docker

```sh
 docker-compose up
```

2. Build from scratch

3. downloads the modules needed to build and test a package

```sh
 go mod download
```

and then using binary

```sh
 go bulid
 ./nsproxy proxy
```


## ***Development***

Open your favorite Terminal and run these commands.

- **First**: run the redis server

- **Second**: Set the Redis host and port in config

- **Third**: DNS Proxy server is listening on localhost `(127.0.0.1:53)`

```sh
 go run main.go proxy
```

- **Forth**: Requesting the domain and getting the desired IP from the proxy server in cmd or anything.

```sh
 nslookup `[Domain name]` 127.0.0.1
```

For example:

> input => nslookup snapp.ir 127.0.0.1

> output

![Ip Address](https://github.com/sina-marefat/NsProxy/blob/main/images/output.png)



## ***Configuration***

The `config.json` file will load the program's script settings,
including cache-expiration, external-dns-server, server port, server host, redis port and redis host in json format when the program is run.

- `cache-expiration-time` : The duration which a domain and it's response set in cache.
- `external-dns-server` : A list of external DNS server IPs to request and update in the cache
- `redis-port` : redis up on the port 6379
- `server-port` : UDP connnection on the port 53
