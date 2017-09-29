package ipecho

import (
	"context"
	"net"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

type dummyResponseWriter struct {
	localAddr  net.Addr
	remoteAddr net.Addr
	msgs       []*dns.Msg
	bytes      []byte
}

func (d *dummyResponseWriter) LocalAddr() net.Addr  { return d.localAddr }
func (d *dummyResponseWriter) RemoteAddr() net.Addr { return d.remoteAddr }
func (d *dummyResponseWriter) WriteMsg(m *dns.Msg) error {
	d.msgs = append(d.msgs, m)
	return nil
}
func (d *dummyResponseWriter) Write(b []byte) (int, error) {
	d.bytes = append(d.bytes, b...)
	return len(b), nil
}
func (*dummyResponseWriter) Close() error        { return nil }
func (*dummyResponseWriter) TsigStatus() error   { return nil }
func (*dummyResponseWriter) TsigTimersOnly(bool) {}
func (*dummyResponseWriter) Hijack()             {}

func (d *dummyResponseWriter) GetMsgs() []*dns.Msg { return d.msgs }
func (d *dummyResponseWriter) ClearMsgs()          { d.msgs = nil }

func (d *dummyResponseWriter) GetBytes() []byte { return d.bytes }
func (d *dummyResponseWriter) ClearBytes()      { d.bytes = nil }

func TestServeDNS(t *testing.T) {
	p := IPEcho{
		Config: &Config{
			Domains: []string{
				"example1.com.",
			},
			TTL:   60,
			Debug: true,
		},
	}

	t.Run("A", func(t *testing.T) {
		d := &dummyResponseWriter{}
		p.ServeDNS(context.Background(), d, &dns.Msg{
			Question: []dns.Question{
				dns.Question{
					Name:   "127.0.0.1.example1.com.",
					Qclass: dns.ClassINET,
					Qtype:  dns.TypeA,
				},
			},
		})
		require.Equal(t, 1, len(d.GetMsgs()))
		require.Equal(t, 1, len(d.GetMsgs()[0].Answer))
		require.Equal(t, dns.Class(dns.ClassINET), dns.Class(d.GetMsgs()[0].Answer[0].Header().Class))
		require.Equal(t, dns.Type(dns.TypeA), dns.Type(d.GetMsgs()[0].Answer[0].Header().Rrtype))
		require.Equal(t, "127.0.0.1.example1.com.", d.GetMsgs()[0].Answer[0].Header().Name)
		require.Equal(t, net.ParseIP("127.0.0.1"), d.GetMsgs()[0].Answer[0].(*dns.A).A)
	})

	t.Run("AAAA", func(t *testing.T) {
		d := &dummyResponseWriter{}
		p.ServeDNS(context.Background(), d, &dns.Msg{
			Question: []dns.Question{
				dns.Question{
					Name:   "::1.example1.com.",
					Qclass: dns.ClassINET,
					Qtype:  dns.TypeAAAA,
				},
			},
		})
		require.Equal(t, 1, len(d.GetMsgs()))
		require.Equal(t, 1, len(d.GetMsgs()[0].Answer))
		require.Equal(t, dns.Class(dns.ClassINET), dns.Class(d.GetMsgs()[0].Answer[0].Header().Class))
		require.Equal(t, dns.Type(dns.TypeAAAA), dns.Type(d.GetMsgs()[0].Answer[0].Header().Rrtype))
		require.Equal(t, "::1.example1.com.", d.GetMsgs()[0].Answer[0].Header().Name)
		require.Equal(t, net.ParseIP("::1"), d.GetMsgs()[0].Answer[0].(*dns.AAAA).AAAA)
	})

	t.Run("Requested A but is AAAA", func(t *testing.T) {
		d := &dummyResponseWriter{}
		p.ServeDNS(context.Background(), d, &dns.Msg{
			Question: []dns.Question{
				dns.Question{
					Name:   "::1.example1.com.",
					Qclass: dns.ClassINET,
					Qtype:  dns.TypeA,
				},
			},
		})

		require.Equal(t, 1, len(d.GetMsgs()))
		require.Equal(t, 1, len(d.GetMsgs()[0].Answer))
		require.Equal(t, dns.Class(dns.ClassINET), dns.Class(d.GetMsgs()[0].Answer[0].Header().Class))
		require.Equal(t, dns.Type(dns.TypeAAAA), dns.Type(d.GetMsgs()[0].Answer[0].Header().Rrtype))
		require.Equal(t, "::1.example1.com.", d.GetMsgs()[0].Answer[0].Header().Name)
		require.Equal(t, net.ParseIP("::1"), d.GetMsgs()[0].Answer[0].(*dns.AAAA).AAAA)
	})

	t.Run("Requested AAAA but is A", func(t *testing.T) {
		d := &dummyResponseWriter{}
		p.ServeDNS(context.Background(), d, &dns.Msg{
			Question: []dns.Question{
				dns.Question{
					Name:   "127.0.0.1.example1.com.",
					Qclass: dns.ClassINET,
					Qtype:  dns.TypeAAAA,
				},
			},
		})
		require.Equal(t, 1, len(d.GetMsgs()))
		require.Equal(t, 1, len(d.GetMsgs()[0].Answer))
		require.Equal(t, dns.Class(dns.ClassINET), dns.Class(d.GetMsgs()[0].Answer[0].Header().Class))
		require.Equal(t, dns.Type(dns.TypeA), dns.Type(d.GetMsgs()[0].Answer[0].Header().Rrtype))
		require.Equal(t, "127.0.0.1.example1.com.", d.GetMsgs()[0].Answer[0].Header().Name)
		require.Equal(t, net.ParseIP("127.0.0.1"), d.GetMsgs()[0].Answer[0].(*dns.A).A)
	})

	t.Run("Invalid Subdomain", func(t *testing.T) {
		d := &dummyResponseWriter{}
		p.ServeDNS(context.Background(), d, &dns.Msg{
			Question: []dns.Question{
				dns.Question{
					Name:   "test.example1.com.",
					Qclass: dns.ClassINET,
					Qtype:  dns.TypeA,
				},
				dns.Question{
					Name:   "test.example1.com.",
					Qclass: dns.ClassINET,
					Qtype:  dns.TypeAAAA,
				},
			},
		})
		require.Equal(t, 0, len(d.GetMsgs()))
	})

	t.Run("No Subdomain", func(t *testing.T) {
		d := &dummyResponseWriter{}
		p.ServeDNS(context.Background(), d, &dns.Msg{
			Question: []dns.Question{
				dns.Question{
					Name:   "example1.com.",
					Qclass: dns.ClassINET,
					Qtype:  dns.TypeA,
				},
			},
		})
		require.Equal(t, 0, len(d.GetMsgs()))
	})

	t.Run("Emtpy Subdomain", func(t *testing.T) {
		d := &dummyResponseWriter{}
		p.ServeDNS(context.Background(), d, &dns.Msg{
			Question: []dns.Question{
				dns.Question{
					Name:   ".example1.com.",
					Qclass: dns.ClassINET,
					Qtype:  dns.TypeA,
				},
			},
		})
		require.Equal(t, 0, len(d.GetMsgs()))
	})

	t.Run("Unknown Domain", func(t *testing.T) {
		d := &dummyResponseWriter{}
		p.ServeDNS(context.Background(), d, &dns.Msg{
			Question: []dns.Question{
				dns.Question{
					Name:   ".example2.com.",
					Qclass: dns.ClassINET,
					Qtype:  dns.TypeA,
				},
			},
		})
		require.Equal(t, 0, len(d.GetMsgs()))
	})

	t.Run("No Questions", func(t *testing.T) {
		d := &dummyResponseWriter{}
		p.ServeDNS(context.Background(), d, &dns.Msg{})
		require.Equal(t, 0, len(d.GetMsgs()))
	})

	t.Run("Invalid Class", func(t *testing.T) {
		d := &dummyResponseWriter{}
		p.ServeDNS(context.Background(), d, &dns.Msg{
			Question: []dns.Question{
				dns.Question{
					Name:   ".example1.com.",
					Qclass: dns.ClassANY,
					Qtype:  dns.TypeA,
				},
			},
		})
		require.Equal(t, 0, len(d.GetMsgs()))
	})

	t.Run("Multiple Questions", func(t *testing.T) {
		d := &dummyResponseWriter{}
		p.ServeDNS(context.Background(), d, &dns.Msg{
			Question: []dns.Question{
				dns.Question{
					Name:   "127.0.0.1.example1.com.",
					Qclass: dns.ClassINET,
					Qtype:  dns.TypeA,
				},
				dns.Question{
					Name:   "::1.example1.com.",
					Qclass: dns.ClassINET,
					Qtype:  dns.TypeAAAA,
				},
			},
		})
		require.Equal(t, 1, len(d.GetMsgs()))
		require.Equal(t, 2, len(d.GetMsgs()[0].Answer))
		require.Equal(t, dns.Class(dns.ClassINET), dns.Class(d.GetMsgs()[0].Answer[0].Header().Class))
		require.Equal(t, dns.Type(dns.TypeA), dns.Type(d.GetMsgs()[0].Answer[0].Header().Rrtype))
		require.Equal(t, "127.0.0.1.example1.com.", d.GetMsgs()[0].Answer[0].Header().Name)
		require.Equal(t, net.ParseIP("127.0.0.1"), d.GetMsgs()[0].Answer[0].(*dns.A).A)
		require.Equal(t, dns.Class(dns.ClassINET), dns.Class(d.GetMsgs()[0].Answer[1].Header().Class))
		require.Equal(t, dns.Type(dns.TypeAAAA), dns.Type(d.GetMsgs()[0].Answer[1].Header().Rrtype))
		require.Equal(t, "::1.example1.com.", d.GetMsgs()[0].Answer[1].Header().Name)
		require.Equal(t, net.ParseIP("::1"), d.GetMsgs()[0].Answer[1].(*dns.AAAA).AAAA)
	})
}
