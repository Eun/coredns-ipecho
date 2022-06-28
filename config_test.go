package ipecho

import (
	"testing"

	"github.com/coredns/caddy/caddyfile"
	"github.com/stretchr/testify/require"
	"github.com/tdewolff/buffer"
)

func TestNewConfigFromDispenser(t *testing.T) {
	t.Run("Valid Config", func(t *testing.T) {
		dispenser := caddyfile.NewDispenser("", buffer.NewReader([]byte(`
			{
				Domain example1.com
				TTL 60
				Domain example2.com
				Domain example1.com
				Debug
			}
		`)))
		config, err := newConfigFromDispenser(dispenser)
		require.NoError(t, err)
		require.NotNil(t, config)
		require.Equal(t, 2, len(config.Domains))
		require.Equal(t, "example1.com.", config.Domains[0])
		require.Equal(t, "example2.com.", config.Domains[1])
		require.Equal(t, uint32(60), config.TTL)
		require.Equal(t, true, config.Debug)
	})
	t.Run("Emtpy Config", func(t *testing.T) {
		dispenser := caddyfile.NewDispenser("", buffer.NewReader([]byte(`
			{
			}
		`)))
		config, err := newConfigFromDispenser(dispenser)
		require.Error(t, err)
		require.Nil(t, config)

		dispenser = caddyfile.NewDispenser("", buffer.NewReader([]byte(``)))
		config, err = newConfigFromDispenser(dispenser)
		require.Error(t, err)
		require.Nil(t, config)
	})

	t.Run("Default Values", func(t *testing.T) {
		dispenser := caddyfile.NewDispenser("", buffer.NewReader([]byte(`
			{
				Domain example1.com
			}
		`)))
		config, err := newConfigFromDispenser(dispenser)
		require.NoError(t, err)
		require.NotNil(t, config)
		require.Equal(t, 1, len(config.Domains))
		require.Equal(t, uint32(2629800), config.TTL)
		require.Equal(t, false, config.Debug)
	})
	t.Run("Invalid Domain", func(t *testing.T) {
		dispenser := caddyfile.NewDispenser("", buffer.NewReader([]byte(`
			{
				Domain 127.0.0.1
			}
		`)))
		config, err := newConfigFromDispenser(dispenser)
		require.Error(t, err)
		require.Nil(t, config)
	})
	t.Run("Invalid TTL", func(t *testing.T) {
		dispenser := caddyfile.NewDispenser("", buffer.NewReader([]byte(`
			{
				Domain example1.com
				TTL SixtySeconds
			}
		`)))
		config, err := newConfigFromDispenser(dispenser)
		require.Error(t, err)
		require.Nil(t, config)
	})
}
