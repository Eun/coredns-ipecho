package ipecho

import (
	"log"
	"net"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

type ipecho struct {
	Next   plugin.Handler
	Config *config
}

// ServeDNS implements the middleware.Handler interface.
func (p ipecho) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	p.echoIP(w, r)
	return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (ipecho) Name() string { return "IPEcho" }

func (p *ipecho) echoIP(w dns.ResponseWriter, r *dns.Msg) {
	if len(r.Question) <= 0 {
		return
	}

	var rrs []dns.RR

	for i := 0; i < len(r.Question); i++ {
		question := r.Question[i]
		if question.Qclass != dns.ClassINET {
			continue
		}

		if question.Qtype == dns.TypeA || question.Qtype == dns.TypeAAAA {
			ip := p.parseIP(&question)
			if ip == nil {
				if p.Config.Debug {
					log.Printf("Parsed IP of '%s' is nil\n", question.Name)
				}
				continue
			}
			// not an ip4
			if ip4 := ip.To4(); ip4 != nil {
				if p.Config.Debug {
					log.Printf("Parsed IP of '%s' is an IPv4 address\n", question.Name)
				}
				rrs = append(rrs, &dns.A{
					Hdr: dns.RR_Header{
						Name:   question.Name,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    p.Config.TTL,
					},
					A: ip,
				})
			} else {
				if p.Config.Debug {
					log.Printf("Parsed IP of '%s' is an IPv6 address\n", question.Name)
				}
				rrs = append(rrs, &dns.AAAA{
					Hdr: dns.RR_Header{
						Name:   question.Name,
						Rrtype: dns.TypeAAAA,
						Class:  dns.ClassINET,
						Ttl:    p.Config.TTL,
					},
					AAAA: ip,
				})
			}
		}
	}

	if len(rrs) > 0 {
		if p.Config.Debug {
			log.Printf("Answering with %d rr's\n", len(rrs))
		}
		w.WriteMsg(&dns.Msg{
			Answer: rrs,
		})
	}
}

func (p *ipecho) parseIP(question *dns.Question) net.IP {
	if p.Config.Debug {
		log.Printf("Query for '%s'", question.Name)
	}

	for _, domain := range p.Config.Domains {
		if strings.HasSuffix(strings.ToLower(question.Name), domain) == true {
			subdomain := question.Name[:len(question.Name)-len(domain)]
			if len(subdomain) <= 0 {
				if p.Config.Debug {
					log.Printf("Query ('%s') has no subomain\n", question.Name)
				}
				return nil
			}
			subdomain = strings.Trim(subdomain, ".")
			if len(subdomain) <= 0 {
				if p.Config.Debug {
					log.Printf("Parsed Subdomain of '%s' is empty\n", question.Name)
				}
				return nil
			}
			if p.Config.Debug {
				log.Printf("Parsed Subdomain of '%s' is '%s'\n", question.Name, subdomain)
			}
			return net.ParseIP(subdomain)
		}
	}

	if p.Config.Debug {
		log.Printf("Query ('%s') does not end with one of the domains (%s)\n", question.Name, strings.Join(p.Config.Domains, ", "))
	}
	return nil
}
