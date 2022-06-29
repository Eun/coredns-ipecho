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
		var err error
		if strings.EqualFold(c.Val(), "domain") {
			err = parseDomainPart(c, &cfg)
		} else if strings.EqualFold(c.Val(), "ttl") {
			err = parseTTLPart(c, &cfg)
		} else if strings.EqualFold(c.Val(), "debug") {
			err = parseDebugPart(c, &cfg)
		}
		if err != nil {
			return nil, err
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

func parseDomainPart(c caddyfile.Dispenser, cfg *config) error {
	if !c.NextArg() {
		return nil
	}
	domain := strings.ToLower(strings.Trim(c.Val(), "."))
	if !govalidator.IsDNSName(domain) {
		return fmt.Errorf("'%s' is not a valid domain name", domain)
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
	return nil
}

func parseTTLPart(c caddyfile.Dispenser, cfg *config) error {
	if !c.NextArg() {
		return nil
	}
	//nolint: gomnd // parse ttl as uint32 with base 10
	ttl, err := strconv.ParseUint(c.Val(), 10, 32)
	if err != nil {
		return fmt.Errorf("invalid TTL value: '%s'", c.Val())
	}
	cfg.TTL = uint32(ttl)
	return nil
}

func parseDebugPart(c caddyfile.Dispenser, cfg *config) error {
	cfg.Debug = true
	return nil
}
