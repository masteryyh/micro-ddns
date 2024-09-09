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
	"errors"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/as/v1/region"
	huaweiv2 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/model"
	"github.com/masteryyh/micro-ddns/internal/config"
	"github.com/masteryyh/micro-ddns/pkg/utils"
	"log/slog"
)

const HuaweiCloudDefaultTTL = 120

type HuaweiCloudDNSUpdateHandler struct {
	domain      string
	subdomain   string
	recordType  RecordType
	zoneId      string
	recordSetId string

	ctx    context.Context
	cancel context.CancelFunc
	client *huaweiv2.DnsClient
	logger *slog.Logger
}

func NewHuaweiCloudDNSUpdateHandler(ddns *config.DDNSSpec, spec *config.HuaweiCloudSpec, parentCtx context.Context, logger *slog.Logger) (*HuaweiCloudDNSUpdateHandler, error) {
	cred, err := basic.NewCredentialsBuilder().WithAk(spec.AccessKey).WithSk(spec.SecretAccessKey).SafeBuild()
	if err != nil {
		return nil, err
	}

	// Use this for error check only for now
	_, err = region.SafeValueOf(spec.Region)
	if err != nil {
		return nil, err
	}

	// Directly specify API endpoint for now to avoid SDK bugs :(
	endpoint := "https://dns." + spec.Region + ".myhuaweicloud.com"
	hcClient, err := huaweiv2.DnsClientBuilder().WithCredential(cred).WithEndpoints([]string{endpoint}).SafeBuild()
	if err != nil {
		return nil, err
	}
	client := huaweiv2.NewDnsClient(hcClient)

	recordType := A
	if ddns.Stack == config.IPv6 {
		recordType = AAAA
	}

	ctx, cancel := context.WithCancel(parentCtx)
	return &HuaweiCloudDNSUpdateHandler{
		// Add a dot at the end of the domain for compatibility
		domain:     ddns.Domain + ".",
		subdomain:  ddns.Subdomain,
		recordType: recordType,
		ctx:        ctx,
		cancel:     cancel,
		client:     client,
		logger:     logger,
	}, nil
}

func (h *HuaweiCloudDNSUpdateHandler) Get() (string, error) {
	if h.zoneId == "" {
		h.logger.Debug("zone id not present, searching")
		result, err := h.client.ListPublicZones(&model.ListPublicZonesRequest{
			Name:       utils.StringPtr(h.domain),
			SearchMode: utils.StringPtr("equal"),
		})
		if err != nil {
			return "", err
		}

		var id string
		for _, zone := range *result.Zones {
			if h.domain == *zone.Name {
				id = *zone.Id
				break
			}
		}

		if id == "" {
			return "", fmt.Errorf("zone " + h.domain + " not exists")
		}
		h.logger.Debug("got zone id " + id)
		h.zoneId = id
	}

	if h.recordSetId == "" {
		h.logger.Debug("record id not present, searching")
		result, err := h.client.ListRecordSetsByZone(&model.ListRecordSetsByZoneRequest{
			ZoneId:     h.zoneId,
			Type:       utils.StringPtr(string(h.recordType)),
			Name:       utils.StringPtr(h.subdomain),
			SearchMode: utils.StringPtr("equal"),
		})
		if err != nil {
			return "", err
		}

		for _, record := range *result.Recordsets {
			if h.subdomain == *record.Name {
				h.logger.Debug("got record id " + *record.Id)
				h.recordSetId = *record.Id
				return (*record.Records)[0], nil
			}
		}
		return "", nil
	}

	result, err := h.client.ShowRecordSet(&model.ShowRecordSetRequest{
		ZoneId:      h.zoneId,
		RecordsetId: h.recordSetId,
	})
	if err != nil {
		hwErr := &sdkerr.ServiceResponseError{}
		if errors.As(err, &hwErr) {
			if hwErr.StatusCode == 404 {
				return "", nil
			}
		}
		return "", err
	}
	return (*result.Records)[0], nil
}

func (h *HuaweiCloudDNSUpdateHandler) Create(address string) error {
	if h.zoneId == "" {
		return fmt.Errorf("zone id is empty")
	}

	h.logger.Debug("creating DNS record for address " + address)
	fqdn := h.subdomain + "." + h.domain
	result, err := h.client.CreateRecordSet(&model.CreateRecordSetRequest{
		ZoneId: h.zoneId,
		Body: &model.CreateRecordSetRequestBody{
			Name:        fqdn,
			Description: utils.StringPtr(Comment),
			Type:        string(h.recordType),
			Records:     []string{address},
			Ttl:         utils.Int32Ptr(HuaweiCloudDefaultTTL),
		},
	})
	if err != nil {
		return err
	}

	h.recordSetId = *result.Id
	return nil
}

func (h *HuaweiCloudDNSUpdateHandler) Update(newAddress string) error {
	if h.zoneId == "" {
		return fmt.Errorf("zone id is empty")
	}

	if h.recordSetId == "" {
		return fmt.Errorf("record id is empty")
	}

	h.logger.Debug("updating DNS record for address " + newAddress)
	fqdn := h.subdomain + "." + h.domain
	_, err := h.client.UpdateRecordSet(&model.UpdateRecordSetRequest{
		ZoneId:      h.zoneId,
		RecordsetId: h.recordSetId,
		Body: &model.UpdateRecordSetReq{
			Name:        &fqdn,
			Description: utils.StringPtr(Comment),
			Type:        utils.StringPtr(string(h.recordType)),
			Ttl:         utils.Int32Ptr(HuaweiCloudDefaultTTL),
			Records:     &[]string{newAddress},
		},
	})
	return err
}
