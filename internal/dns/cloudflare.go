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
	"github.com/cloudflare/cloudflare-go"
	"github.com/masteryyh/micro-ddns/internal/config"
	"github.com/masteryyh/micro-ddns/pkg/utils"
	"log/slog"
	"time"
)

const (
	PerPageCount = 500
	TTL          = 120
	Comment      = "Created by micro-ddns"
)

type CloudflareDNSUpdateHandler struct {
	domain     string
	subdomain  string
	recordType RecordType
	zoneId     string
	recordId   string

	ctx       context.Context
	cancel    context.CancelFunc
	apiClient *cloudflare.API
	logger    *slog.Logger
}

func NewCloudflareDNSUpdateHandler(ddns *config.DDNSSpec, cloudflareSpec *config.CloudflareSpec, parentCtx context.Context, logger *slog.Logger) (*CloudflareDNSUpdateHandler, error) {
	var client *cloudflare.API
	if !utils.IsEmpty(cloudflareSpec.APIToken) {
		c, err := cloudflare.NewWithAPIToken(utils.StringPtrToString(cloudflareSpec.APIToken))
		if err != nil {
			return nil, err
		}
		client = c
	} else if !utils.IsEmpty(cloudflareSpec.GlobalAPIKey) && !utils.IsEmpty(cloudflareSpec.Email) {
		c, err := cloudflare.New(utils.StringPtrToString(cloudflareSpec.GlobalAPIKey), utils.StringPtrToString(cloudflareSpec.Email))
		if err != nil {
			return nil, err
		}
		client = c
	}

	recordType := A
	if ddns.Stack == config.IPv6 {
		recordType = AAAA
	}

	ctx, cancel := context.WithCancel(parentCtx)
	return &CloudflareDNSUpdateHandler{
		domain:     ddns.Domain,
		subdomain:  ddns.Subdomain,
		recordType: recordType,
		apiClient:  client,
		ctx:        ctx,
		cancel:     cancel,
		logger:     logger,
	}, nil
}

func (h *CloudflareDNSUpdateHandler) fetchZoneId() error {
	zoneCtx, zoneCancel := context.WithTimeout(h.ctx, time.Second*30)
	defer zoneCancel()

	if h.zoneId == "" {
		h.logger.Debug("looking for user's DNS zone")
		zones, err := h.apiClient.ListZones(zoneCtx, h.domain)
		if err != nil {
			return err
		}

		for _, zone := range zones {
			if zone.Name == h.domain {
				h.zoneId = zone.ID
				break
			}
		}
	}

	if h.zoneId == "" {
		return fmt.Errorf("no corresponding DNS zone found")
	}
	h.logger.Debug("found DNS zone ID " + h.zoneId)
	return nil
}

func (h *CloudflareDNSUpdateHandler) Get() (string, error) {
	ctx, cancel := context.WithCancel(h.ctx)
	defer cancel()

	if h.zoneId == "" {
		if err := h.fetchZoneId(); err != nil {
			return "", err
		}
	}

	if h.recordId == "" {
		h.logger.Debug("looking for current DNS record ID")
		recordCtx, recordCancel := context.WithTimeout(ctx, time.Second*30)
		records, resultInfo, err := h.apiClient.ListDNSRecords(recordCtx, cloudflare.ZoneIdentifier(h.zoneId), cloudflare.ListDNSRecordsParams{
			Type: string(h.recordType),
			ResultInfo: cloudflare.ResultInfo{
				Page:    1,
				PerPage: PerPageCount,
			},
		})
		if err != nil {
			recordCancel()
			return "", err
		}

		if resultInfo.Count > PerPageCount {
			pages := resultInfo.Count/PerPageCount - 1
			if resultInfo.Count%PerPageCount != 0 {
				pages += 1
			}

			for i := 0; i < pages; i++ {
				additionalRecordCtx, additionalRecordCancel := context.WithTimeout(ctx, time.Second*30)
				additionalRecords, _, err := h.apiClient.ListDNSRecords(additionalRecordCtx, cloudflare.ZoneIdentifier(h.zoneId), cloudflare.ListDNSRecordsParams{
					Type: string(h.recordType),
					ResultInfo: cloudflare.ResultInfo{
						Page:    i + 2,
						PerPage: PerPageCount,
					},
				})

				if err != nil {
					additionalRecordCancel()
					recordCancel()
					return "", err
				}
				if len(additionalRecords) > 0 {
					records = append(records, additionalRecords...)
				}
				additionalRecordCancel()
			}
		}

		var id string
		for _, record := range records {
			if record.Name == h.subdomain {
				id = record.ID
				break
			}
		}

		if id == "" {
			recordCancel()
			return "", nil
		}

		h.logger.Debug("found DNS record id: " + id)
		recordCancel()
	}

	dnsCtx, dnsCancel := context.WithTimeout(ctx, time.Second*30)
	defer dnsCancel()
	h.logger.Debug("getting DNS record detail for record ID " + h.recordId)
	record, err := h.apiClient.GetDNSRecord(dnsCtx, cloudflare.ZoneIdentifier(h.zoneId), h.recordId)
	if err != nil {
		return "", err
	}
	return record.Content, nil
}

func (h *CloudflareDNSUpdateHandler) Create(address string) error {
	if h.zoneId == "" {
		h.logger.Debug("DNS zone ID is empty, searching")
		if err := h.fetchZoneId(); err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(h.ctx, time.Second*30)
	defer cancel()

	if h.recordId == "" {
		h.logger.Debug("creating DNS record")
		record, err := h.apiClient.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(h.zoneId), cloudflare.CreateDNSRecordParams{
			Type:    string(h.recordType),
			Name:    h.subdomain,
			Content: address,
			ZoneID:  h.zoneId,
			TTL:     TTL,
			Proxied: utils.BoolPtr(false),
			Comment: Comment,
		})

		if err != nil {
			return err
		}
		h.recordId = record.ID
	}
	return nil
}

func (h *CloudflareDNSUpdateHandler) Update(newAddress string) error {
	if h.zoneId == "" {
		return fmt.Errorf("zoneId is empty")
	}

	if h.recordId == "" {
		return fmt.Errorf("recordId is empty")
	}

	ctx, cancel := context.WithTimeout(h.ctx, time.Second*30)
	defer cancel()

	h.logger.Debug("updating DNS record for record ID " + h.recordId)
	_, err := h.apiClient.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(h.zoneId), cloudflare.UpdateDNSRecordParams{
		Type:    string(h.recordType),
		Name:    h.subdomain,
		Content: newAddress,
		TTL:     TTL,
		Proxied: utils.BoolPtr(false),
	})
	return err
}
