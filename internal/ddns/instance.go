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

package ddns

import (
	"context"
	"log/slog"

	"github.com/masteryyh/micro-ddns/internal/config"
	"github.com/masteryyh/micro-ddns/internal/dns"
	"github.com/masteryyh/micro-ddns/internal/ip"
)

type DDNSInstance struct {
	spec *config.DDNSSpec

	dnsHandler      dns.DNSUpdateHandler
	addressDetector ip.AddressDetector
	logger          *slog.Logger
}

func NewDDNSInstance(ddnsSpec *config.DDNSSpec, logger *slog.Logger) (*DDNSInstance, error) {
	var handler dns.DNSUpdateHandler

	providerSpec := ddnsSpec.GetProviderSpec()
	providerType := providerSpec.GetType()
	switch providerType {
	case config.DNSProviderCloudflare:
		spec := providerSpec.Cloudflare
		h, err := dns.NewCloudflareDNSUpdateHandler(ddnsSpec, spec, logger)
		if err != nil {
			return nil, err
		}
		handler = h
	case config.DNSProviderAliCloud:
		spec := providerSpec.AliCloud
		h, err := dns.NewAliCloudDNSUpdateHandler(ddnsSpec, spec, logger)
		if err != nil {
			return nil, err
		}
		handler = h
	case config.DNSProviderDNSPod:
		spec := providerSpec.DNSPod
		h, err := dns.NewDNSPodDNSUpdateHandler(ddnsSpec, spec, logger)
		if err != nil {
			return nil, err
		}
		handler = h
	case config.DNSProviderHuaweiCloud:
		spec := providerSpec.Huawei
		h, err := dns.NewHuaweiCloudDNSUpdateHandler(ddnsSpec, spec, logger)
		if err != nil {
			return nil, err
		}
		handler = h
	case config.DNSProviderJDCloud:
		spec := providerSpec.JD
		h, err := dns.NewJDCloudDNSUpdateHandler(ddnsSpec, spec, logger)
		if err != nil {
			return nil, err
		}
		handler = h
	case config.DNSProviderRFC2136:
		spec := providerSpec.RFC2136
		h, err := dns.NewRFC2136DNSUpdateHandler(ddnsSpec, spec, logger)
		if err != nil {
			return nil, err
		}
		handler = h
	}

	var addrDetector ip.AddressDetector

	detectionSpec := ddnsSpec.GetDetectionSpec()
	detectionType := detectionSpec.GetDetectionType()
	switch detectionType {
	case config.AddressDetectionIface:
		addrDetector = ip.NewIfaceAddressDetector(detectionSpec, ddnsSpec.Stack, logger)
	case config.AddressDetectionThirdParty:
		addrDetector = ip.NewThirdPartyAddressDetector(detectionSpec, ddnsSpec.Stack, logger)
	}

	return &DDNSInstance{
		spec:            ddnsSpec,
		dnsHandler:      handler,
		addressDetector: addrDetector,
		logger:          logger,
	}, nil
}

func (n *DDNSInstance) DoUpdate(parentCtx context.Context) error {
	n.logger.Info("detecting current address", "name", n.spec.Name)
	addr, err := n.addressDetector.Detect(parentCtx)
	if err != nil {
		n.logger.Error("error detecting address", "name", n.spec.Name, "err", err)
		return err
	}

	n.logger.Info("getting current address registered with DNS provider", "name", n.spec.Name)
	recordAddr, err := n.dnsHandler.Get(parentCtx)
	if err != nil {
		n.logger.Error("error getting current address", "name", n.spec.Name, "err", err)
		return err
	}
	if recordAddr == "" {
		n.logger.Info("DNS record for this subdomain not found or ignored, creating", "name", n.spec.Name, "domain", n.spec.Domain, "subdomain", n.spec.Subdomain)
		return n.dnsHandler.Create(parentCtx, addr)
	}

	if recordAddr != addr {
		n.logger.Info("address changed, updating DNS record", "name", n.spec.Name, "domain", n.spec.Domain, "subdomain", n.spec.Subdomain, "address", addr)
		return n.dnsHandler.Update(parentCtx, addr)
	}
	n.logger.Info("address not changed, skipping")
	return nil
}
