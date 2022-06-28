package ipecho

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/coredns/caddy/caddyfile"
)

type config struct {
	// Domains defines the Domains we will react to
	Domains []string
	// TTL to use for response
	TTL uint32
	// Debug mode
	Debug bool
}

func newConfigFromDispenser(c caddyfile.Dispenser) (*config, error) {
	config := config{
		TTL: 2629800,
	}

	for c.NextBlock() {
		if strings.EqualFold(c.Val(), "domain") {
			if c.NextArg() {
				domain := strings.ToLower(strings.Trim(c.Val(), "."))
				if !govalidator.IsDNSName(domain) {
					return nil, fmt.Errorf("'%s' is not a valid domain name", domain)
				}
				domain = domain + "."

				exists := false
				for i := range config.Domains {
					if config.Domains[i] == domain {
						exists = true
						break
					}
				}

				if exists == false {
					config.Domains = append(config.Domains, domain)
				}
			}
		} else if strings.EqualFold(c.Val(), "ttl") {
			if c.NextArg() {
				ttl, err := strconv.ParseUint(c.Val(), 10, 32)
				if err != nil {
					return nil, fmt.Errorf("Invalid TTL value: '%s'", c.Val())
				}
				config.TTL = uint32(ttl)
			}
		} else if strings.EqualFold(c.Val(), "debug") {
			config.Debug = true
		}
	}
	if config.Debug {
		log.Println("[ipecho] Debug Mode is on")
		log.Printf("[ipecho] Parsed %d Domains: %s\n", len(config.Domains), strings.Join(config.Domains, ", "))
		log.Printf("[ipecho] TTL is %d", config.TTL)
	}
	if len(config.Domains) <= 0 {
		return nil, fmt.Errorf("There is no domain to handle")
	}
	return &config, nil
}
