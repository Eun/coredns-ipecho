package ipecho

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/coredns/caddy"
)

func init() {
	caddy.RegisterPlugin("ipecho", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	c.Next()
	config, err := newConfigFromDispenser(c.Dispenser)
	if err != nil {
		return plugin.Error("ipecho", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return ipecho{Next: next, Config: config}
	})

	return nil
}
