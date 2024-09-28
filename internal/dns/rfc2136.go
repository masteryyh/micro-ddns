/*
Copyright © 2024 masteryyh <yyh991013@163.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dns

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/bodgit/tsig"
	"github.com/bodgit/tsig/gss"
	"github.com/masteryyh/micro-ddns/internal/config"
	"github.com/miekg/dns"
)

const RFC2136DefaultTTL = 120

type RFC2136DNSUpdateHandler struct {
	domain     string
	subdomain  string
	recordType RecordType
	server     string

	spec       *config.RFC2136Spec
	gssKeyName string
	keyName    string
	lastRR     string
	client     *dns.Client
	logger     *slog.Logger
}

func NewRFC2136DNSUpdateHandler(ddns *config.DDNSSpec, spec *config.RFC2136Spec, logger *slog.Logger) (*RFC2136DNSUpdateHandler, error) {
	port := 53
	if spec.Port != nil {
		port = *spec.Port
	}

	recordType := A
	if ddns.Stack == config.IPv6 {
		recordType = AAAA
	}

	server := spec.Address + ":" + strconv.Itoa(port)
	handler := &RFC2136DNSUpdateHandler{
		domain:     ddns.Domain,
		subdomain:  ddns.Subdomain,
		recordType: recordType,
		server:     server,
		spec:       spec,
		logger:     logger,
	}

	return handler, nil
}

func (h *RFC2136DNSUpdateHandler) negotiate(ctx context.Context) error {
	client := &dns.Client{}
	if h.spec.UseTCP != nil && *h.spec.UseTCP {
		client.Net = "tcp"
	}

	if h.spec.TSIG != nil {
		hmac := tsig.HMAC{
			dns.Fqdn(h.spec.TSIG.KeyName): h.spec.TSIG.Key,
		}
		h.client.TsigProvider = hmac
	} else if h.spec.GSSTSIG != nil {
		gssClient, err := gss.NewClient(client)
		if err != nil {
			return err
		}
		keyName, _, err := gssClient.NegotiateContextWithCredentials(h.server, h.spec.GSSTSIG.Domain, h.spec.GSSTSIG.Username, h.spec.GSSTSIG.Password)
		if err != nil {
			return err
		}

		h.gssKeyName = keyName
		client.TsigProvider = gssClient

		go func(client *gss.Client) {
			<-ctx.Done()
			if client != nil {
				client.Close()
				h.gssKeyName = ""
				h.client = nil
			}
		}(gssClient)
	}

	h.client = client
	return nil
}

func (h *RFC2136DNSUpdateHandler) Get(parentCtx context.Context) (string, error) {
	message := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Id:               dns.Id(),
			RecursionDesired: false,
		},
	}

	fqdn := dns.Fqdn(h.subdomain + "." + h.domain)
	qtype := dns.TypeA
	if h.recordType == AAAA {
		qtype = dns.TypeAAAA
	}
	message.Question = []dns.Question{
		{
			Name:   fqdn,
			Qtype:  qtype,
			Qclass: dns.ClassINET,
		},
	}

	h.logger.Debug("querying DNS server for current address")
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()
	if err := h.negotiate(ctx); err != nil {
		return "", err
	}
	result, _, err := h.client.ExchangeContext(ctx, message, h.server)
	if err != nil {
		return "", err
	}

	h.logger.Debug("got " + strconv.Itoa(len(result.Answer)) + " records")
	if len(result.Answer) == 0 {
		return "", nil
	}

	for _, ans := range result.Answer {
		str := ans.String()
		if rr, ok := ans.(*dns.A); ok {
			if str == h.lastRR {
				return rr.A.String(), nil
			}
		} else if rr, ok := ans.(*dns.AAAA); ok {
			if str == h.lastRR {
				return rr.AAAA.String(), nil
			}
		}
	}
	return "", nil
}

func (h *RFC2136DNSUpdateHandler) Create(parentCtx context.Context, address string) error {
	h.logger.Debug("creating DNS record for address " + address)
	message := &dns.Msg{}
	message.SetUpdate(dns.Fqdn(h.domain))

	fqdn := dns.Fqdn(h.subdomain + "." + h.domain)
	rrStr := fqdn + "\t" + strconv.Itoa(RFC2136DefaultTTL) + "\tIN\t" + string(h.recordType) + "\t" + address
	rr, err := dns.NewRR(rrStr)
	if err != nil {
		return err
	}
	message.Insert([]dns.RR{rr})

	h.logger.Debug("RR about to create: " + rr.String())
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()
	if err := h.negotiate(ctx); err != nil {
		return err
	}

	if h.gssKeyName != "" {
		message.SetTsig(h.gssKeyName, tsig.GSS, 300, time.Now().Unix())
	} else if h.keyName != "" {
		message.SetTsig(h.keyName, dns.HmacSHA256, 300, time.Now().Unix())
	}

	_, _, err = h.client.ExchangeContext(ctx, message, h.server)
	if err != nil {
		return err
	}

	h.lastRR = rrStr
	return nil
}

func (h *RFC2136DNSUpdateHandler) Update(parentCtx context.Context, newAddress string) error {
	h.logger.Debug("updating DNS record for new address " + newAddress)
	if h.lastRR == "" {
		return fmt.Errorf("last address unknown")
	}

	message := &dns.Msg{}
	message.SetUpdate(dns.Fqdn(h.domain))

	fqdn := dns.Fqdn(h.subdomain + "." + h.domain)
	rrStr := fqdn + "\t" + strconv.Itoa(RFC2136DefaultTTL) + "\tIN\t" + string(h.recordType) + "\t" + newAddress
	rr, err := dns.NewRR(rrStr)
	if err != nil {
		return err
	}
	message.Insert([]dns.RR{rr})

	oldRr, err := dns.NewRR(h.lastRR)
	if err != nil {
		return err
	}
	message.Remove([]dns.RR{oldRr})

	h.logger.Debug("RR about to update: " + rr.String())
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()
	if err := h.negotiate(ctx); err != nil {
		return err
	}

	if h.gssKeyName != "" {
		message.SetTsig(h.gssKeyName, tsig.GSS, 300, time.Now().Unix())
	} else if h.keyName != "" {
		message.SetTsig(h.keyName, dns.HmacSHA256, 300, time.Now().Unix())
	}

	_, _, err = h.client.ExchangeContext(ctx, message, h.server)
	if err != nil {
		return err
	}

	h.lastRR = rrStr
	return nil
}
