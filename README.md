# ipecho 
[![Actions Status](https://github.com/Eun/coredns-ipecho/workflows/push/badge.svg)](https://github.com/Eun/coredns-ipecho/actions)
[![go-report](https://goreportcard.com/badge/github.com/Eun/coredns-ipecho)](https://goreportcard.com/report/github.com/Eun/coredns-ipecho)
---
*ipecho* is an [coredns](https://github.com/coredns/coredns/) plugin, it parses the IP out of a subdomain and echos it back as an record.

## Example
```
A IN 127.0.0.1.example.com. -> A: 127.0.0.1
AAAA IN ::1.example.com. -> AAAA: ::1
```

## Syntax
```
ipecho {
    domain example1.com
    domain example2.com
    ttl 2629800
}
```

* **domain** adds the domain that should be handled
* **ttl** defines the ttl that should be used in the response
* **debug** enables debug logging
