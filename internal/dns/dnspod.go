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
	"github.com/masteryyh/micro-ddns/internal/config"
	"github.com/masteryyh/micro-ddns/pkg/utils"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"log/slog"
	"time"
)

const (
	DNSPodDefaultTTL = 120
	Comment          = "Created/Updated by micro-ddns"
)

type DNSPodDNSUpdateHandler struct {
	domain     string
	subdomain  string
	recordType RecordType
	line       string

	ctx    context.Context
	cancel context.CancelFunc
	client *dnspod.Client
	logger *slog.Logger

	domainId *uint64
	recordId *uint64
}

func NewDNSPodDNSUpdateHandler(ddns *config.DDNSSpec, spec *config.DNSPodSpec, parentCtx context.Context, logger *slog.Logger) (*DNSPodDNSUpdateHandler, error) {
	credential := common.NewCredential(spec.SecretID, spec.SecretKey)
	client, err := dnspod.NewClient(credential, spec.Region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	recordType := AAAA
	if ddns.Stack == config.IPv4 {
		recordType = A
	}

	line := "0"
	if spec.LineID != nil {
		line = *spec.LineID
	}

	ctx, cancel := context.WithCancel(parentCtx)
	return &DNSPodDNSUpdateHandler{
		domain:     ddns.Domain,
		subdomain:  ddns.Subdomain,
		recordType: recordType,
		line:       line,
		ctx:        ctx,
		cancel:     cancel,
		client:     client,
		logger:     logger,
	}, nil
}

func (h *DNSPodDNSUpdateHandler) findDomainId() error {
	ctx, cancel := context.WithTimeout(h.ctx, 30*time.Second)
	defer cancel()
	result, err := h.client.DescribeDomainListWithContext(ctx, &dnspod.DescribeDomainListRequest{
		Keyword: &h.domain,
		Limit:   utils.Int64Ptr(PerPageCount),
	})
	if err != nil {
		return err
	}

	list := result.Response.DomainList
	for _, domain := range list {
		if *domain.Name == h.domain {
			h.logger.Debug("got domain id", "id", domain.DomainId)
			h.domainId = domain.DomainId
			break
		}
	}
	if h.domainId == nil {
		if *result.Response.DomainCountInfo.DomainTotal <= PerPageCount {
			return fmt.Errorf("domain " + h.domain + " not exists in the account")
		}

		total := int(*result.Response.DomainCountInfo.DomainTotal)
		pages := (total - PerPageCount) / PerPageCount
		if pages%PerPageCount != 0 {
			pages++
		}

		for i := 1; i <= pages; i++ {
			pageCtx, pageCancel := context.WithTimeout(h.ctx, 30*time.Second)
			pageResult, err := h.client.DescribeDomainListWithContext(pageCtx, &dnspod.DescribeDomainListRequest{
				Keyword: &h.domain,
				Limit:   utils.Int64Ptr(PerPageCount),
				Offset:  utils.Int64Ptr(int64(PerPageCount * i)),
			})
			if err != nil {
				pageCancel()
				return err
			}

			list = pageResult.Response.DomainList
			for _, domain := range list {
				if *domain.Name == h.domain {
					h.logger.Debug("got domain id", "id", *domain.DomainId)
					h.domainId = domain.DomainId
					break
				}
			}
			pageCancel()
		}
	}

	if h.recordId == nil {
		return fmt.Errorf("domain " + h.domain + " not exists in the account")
	}
	return nil
}

func (h *DNSPodDNSUpdateHandler) findRecordId() (string, error) {
	ctx, cancel := context.WithTimeout(h.ctx, 30*time.Second)
	defer cancel()
	result, err := h.client.DescribeRecordListWithContext(ctx, &dnspod.DescribeRecordListRequest{
		DomainId:     h.domainId,
		Subdomain:    &h.subdomain,
		RecordType:   utils.StringPtr(string(h.recordType)),
		RecordLineId: &h.line,
	})
	if err != nil {
		return "", err
	}

	list := result.Response.RecordList
	for _, record := range list {
		if *record.Name == h.subdomain {
			h.recordId = record.RecordId
			h.logger.Debug("got record id", "id", *record.RecordId)
			return *record.Value, nil
		}
	}

	if h.recordId == nil {
		if *result.Response.RecordCountInfo.SubdomainCount <= PerPageCount {
			return "", nil
		}

		total := int(*result.Response.RecordCountInfo.SubdomainCount)
		pages := (total - PerPageCount) / PerPageCount
		if pages%PerPageCount != 0 {
			pages++
		}

		for i := 1; i <= pages; i++ {
			pageCtx, pageCancel := context.WithTimeout(h.ctx, 30*time.Second)
			pageResult, err := h.client.DescribeRecordListWithContext(pageCtx, &dnspod.DescribeRecordListRequest{
				DomainId:     h.domainId,
				Subdomain:    &h.subdomain,
				RecordType:   utils.StringPtr(string(h.recordType)),
				RecordLineId: &h.line,
				Offset:       utils.Uint64Ptr(uint64(PerPageCount * i)),
			})
			if err != nil {
				pageCancel()
				return "", err
			}

			list = pageResult.Response.RecordList
			for _, record := range list {
				if *record.Name == h.subdomain {
					h.recordId = record.RecordId
					h.logger.Debug("got record id", "id", *record.RecordId)
					pageCancel()
					return *record.Value, nil
				}
			}
			pageCancel()
		}
	}
	return "", nil
}

func (h *DNSPodDNSUpdateHandler) Get() (string, error) {
	if h.domainId == nil {
		h.logger.Debug("no domain id present, searching")
		if err := h.findDomainId(); err != nil {
			return "", err
		}
	}

	if h.recordId == nil {
		h.logger.Debug("no record id present, searching")
		return h.findRecordId()
	}

	ctx, cancel := context.WithTimeout(h.ctx, 30*time.Second)
	defer cancel()
	result, err := h.client.DescribeRecordWithContext(ctx, &dnspod.DescribeRecordRequest{
		DomainId: h.domainId,
		RecordId: h.recordId,
	})
	if err != nil {
		return "", err
	}
	return *result.Response.RecordInfo.Value, nil
}

func (h *DNSPodDNSUpdateHandler) Create(address string) error {
	if h.domainId == nil {
		return fmt.Errorf("domain id is empty")
	}

	h.logger.Debug("creating DNS record for domain " + h.subdomain + "." + h.domain)
	ctx, cancel := context.WithTimeout(h.ctx, 30*time.Second)
	defer cancel()
	result, err := h.client.CreateRecordWithContext(ctx, &dnspod.CreateRecordRequest{
		DomainId:     h.domainId,
		RecordType:   utils.StringPtr(string(h.recordType)),
		RecordLineId: &h.line,
		SubDomain:    &h.subdomain,
		TTL:          utils.Uint64Ptr(DNSPodDefaultTTL),
		Value:        &address,
		Remark:       utils.StringPtr(Comment),
	})
	if err != nil {
		return err
	}

	h.recordId = result.Response.RecordId
	return nil
}

func (h *DNSPodDNSUpdateHandler) Update(newAddress string) error {
	if h.domainId == nil {
		return fmt.Errorf("domain id is empty")
	}

	if h.recordId == nil {
		return fmt.Errorf("record id is empty")
	}

	h.logger.Debug("updating DNS record for domain " + h.subdomain + "." + h.domain)
	ctx, cancel := context.WithTimeout(h.ctx, 30*time.Second)
	defer cancel()
	_, err := h.client.ModifyRecordWithContext(ctx, &dnspod.ModifyRecordRequest{
		DomainId:     h.domainId,
		RecordId:     h.recordId,
		RecordType:   utils.StringPtr(string(h.recordType)),
		RecordLineId: &h.line,
		TTL:          utils.Uint64Ptr(DNSPodDefaultTTL),
		Value:        &newAddress,
		Remark:       utils.StringPtr(Comment),
	})
	return err
}
