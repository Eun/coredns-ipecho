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

const (
	defaultTTL = 2629800
)

func newConfigFromDispenser(c caddyfile.Dispenser) (*config, error) {
	cfg := config{
		TTL: defaultTTL,
	}

	for c.NextBlock() {
		if strings.EqualFold(c.Val(), "domain") {
			parseDomainPart(c, &cfg)
		} else if strings.EqualFold(c.Val(), "ttl") {
			parseTTLPart(c, &cfg)
		} else if strings.EqualFold(c.Val(), "debug") {
			cfg.Debug = true
		}
	}
	if cfg.Debug {
		log.Println("[ipecho] Debug Mode is on")
		log.Printf("[ipecho] Parsed %d Domains: %s\n", len(cfg.Domains), strings.Join(cfg.Domains, ", "))
		log.Printf("[ipecho] TTL is %d", cfg.TTL)
	}
	if len(cfg.Domains) == 0 {
		return nil, fmt.Errorf("there is no domain to handle")
	}
	return &cfg, nil
}

func parseDomainPart(c caddyfile.Dispenser, cfg *config) {
	if !c.NextArg() {
		return
	}
	domain := strings.ToLower(strings.Trim(c.Val(), "."))
	if !govalidator.IsDNSName(domain) {
		return nil, fmt.Errorf("'%s' is not a valid domain name", domain)
	}
	domain += "."

	exists := false
	for i := range cfg.Domains {
		if cfg.Domains[i] == domain {
			exists = true
			break
		}
	}

	if !exists {
		cfg.Domains = append(cfg.Domains, domain)
	}
}

func parseTTLPart(c caddyfile.Dispenser, cfg *config) {
	if !c.NextArg() {
		return
	}
	//noling: gomnd // parse ttl as uint32 with base 10
	ttl, err := strconv.ParseUint(c.Val(), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid TTL value: '%s'", c.Val())
	}
	cfg.TTL = uint32(ttl)
}