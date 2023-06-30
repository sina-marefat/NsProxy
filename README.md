***DNS Proxy***

# Introduction
"A UDP connection is established on port 53. Whenever a request is made on this connection, the data is buffered, and then passed to the parser module. The parser module handles DNS requests, which are identified by a 12-byte header. Using this module, we separate the nameserver from the query, and check the query type (e.g., AAAA, A, etc.).  If the query type is incorrect, a DNS error is returned. Otherwise, we check all the keys using the nameserver. If a key with the specified name exists and is already cached, we retrieve it. Otherwise, a UDP connection is established with the DNS server.  We send the request and receive the response, which is then stored in the UDP connection proxy to be displayed in the command line (cmd)."

![DNS Resolution](https://github.com/sina-marefat/NsProxy/blob/main/images/dns.jpg)

As it is clear from the diagram, with the increase in the percentage of proxy use, the utilization time ratio decreases

## What is The DNSProxy?
A DNS proxy, also known as a DNS forwarder or a DNS resolver, is an intermediary server that sits between client devices and DNS servers. Its primary function is to handle DNS queries and forward them to appropriate DNS servers for resolution.

When a client device sends a DNS query to a DNS proxy, the proxy examines the request and acts as a middleman between the client and the DNS server. The DNS proxy checks its cache to see if it already has the resolved IP address for the requested domain. If the information is available in its cache, the proxy can respond directly to the client without sending a query to an external DNS server. This caching mechanism helps to improve the response time and reduce the load on DNS servers.

## Requirment
To run this project, you will need:
- Golang
- Redis

## Installation

First: Downloads the modules needed to build and test a package
```sh
 go mod download
```

Second: Using Docker
```sh
 docker-compose up
```
and then using binary
```sh
 go bulid
 ./nsproxy proxy
```



