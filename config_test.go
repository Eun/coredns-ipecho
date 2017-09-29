package ipecho

import "testing"
import "github.com/mholt/caddy/caddyfile"
import "github.com/tdewolff/buffer"
import "github.com/stretchr/testify/require"

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
		config, err := NewConfigFromDispenser(dispenser)
		require.NoError(t, err)
		require.NotNil(t, config)
		require.Equal(t, 2, len(config.Domains))
		require.Equal(t, "example1.com.", config.Domains[0])
		require.Equal(t, "example2.com.", config.Domains[1])
		require.Equal(t, uint32(60), config.TTL)
		require.Equal(t, true, config.Debug)
	})
	t.Run("Default Values", func(t *testing.T) {
		dispenser := caddyfile.NewDispenser("", buffer.NewReader([]byte(`
			{
			}
		`)))
		config, err := NewConfigFromDispenser(dispenser)
		require.NoError(t, err)
		require.NotNil(t, config)
		require.Equal(t, 0, len(config.Domains))
		require.Equal(t, uint32(2629800), config.TTL)
		require.Equal(t, false, config.Debug)

		dispenser = caddyfile.NewDispenser("", buffer.NewReader([]byte(``)))
		config, err = NewConfigFromDispenser(dispenser)
		require.NoError(t, err)
		require.NotNil(t, config)
		require.Equal(t, 0, len(config.Domains))
		require.Equal(t, uint32(2629800), config.TTL)
		require.Equal(t, false, config.Debug)
	})
	t.Run("Invalid TTL", func(t *testing.T) {
		dispenser := caddyfile.NewDispenser("", buffer.NewReader([]byte(`
			{
				TTL SixtySeconds
			}
		`)))
		config, err := NewConfigFromDispenser(dispenser)
		require.Error(t, err)
		require.Nil(t, config)
	})
	t.Run("Invalid Domain", func(t *testing.T) {
		dispenser := caddyfile.NewDispenser("", buffer.NewReader([]byte(`
			{
				Domain 127.0.0.1
			}
		`)))
		config, err := NewConfigFromDispenser(dispenser)
		require.Error(t, err)
		require.Nil(t, config)
	})
}
