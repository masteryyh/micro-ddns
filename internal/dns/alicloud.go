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
	alidns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/masteryyh/micro-ddns/internal/config"
	"github.com/masteryyh/micro-ddns/pkg/utils"
	"log/slog"
	"time"
)

const (
	AliCloudDefaultTTL = 600
)

type AliCloudDNSUpdateHandler struct {
	domain     string
	subdomain  string
	recordType RecordType
	line       string
	recordId   string

	ctx    context.Context
	cancel context.CancelFunc
	client *alidns.Client
	logger *slog.Logger
}

func NewAliCloudDNSUpdateHandler(ddns *config.DDNSSpec, aliSpec *config.AliCloudSpec, parentCtx context.Context, logger *slog.Logger) (*AliCloudDNSUpdateHandler, error) {
	clientConfig := &openapi.Config{
		AccessKeyId:     &aliSpec.AccessKeyID,
		AccessKeySecret: &aliSpec.AccessKeySecret,
		RegionId:        &aliSpec.RegionID,
	}
	client, err := alidns.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	recordType := A
	if ddns.Stack == config.IPv6 {
		recordType = AAAA
	}

	line := "default"
	if aliSpec.Line != nil {
		line = *aliSpec.Line
	}

	ctx, cancel := context.WithCancel(parentCtx)
	return &AliCloudDNSUpdateHandler{
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

func (h *AliCloudDNSUpdateHandler) Get() (string, error) {
	if h.recordId != "" {
		h.logger.Debug("record id present, getting record info")
		result, err := h.client.DescribeDomainRecordInfo(&alidns.DescribeDomainRecordInfoRequest{
			RecordId: &h.recordId,
		})
		if err != nil {
			aliErr := &tea.SDKError{}
			if errors.As(err, &aliErr) {
				if *aliErr.Code == "InvalidRR.NoExist" {
					return "", nil
				}
			}
			return "", err
		}

		val := *result.Body.Value
		h.logger.Debug("got current ip address registered: " + val)
		return val, nil
	}

	h.logger.Debug("no record id present, searching for records already exists")
	result, err := h.client.DescribeDomainRecords(&alidns.DescribeDomainRecordsRequest{
		DomainName: &h.domain,
		RRKeyWord:  &h.subdomain,
		Type:       utils.StringPtr(string(h.recordType)),
		Line:       &h.line,
		PageNumber: utils.Int64Ptr(1),
		PageSize:   utils.Int64Ptr(PerPageCount),
	})
	if err != nil {
		return "", err
	}

	if *result.Body.TotalCount == 0 {
		h.logger.Debug("no record with subdomain " + h.subdomain + " found")
		return "", nil
	}

	thisPageCount := (int)(*result.Body.PageSize)
	for i := 0; i < thisPageCount; i++ {
		record := result.Body.DomainRecords.Record[i]
		if *record.RR == h.subdomain && *record.DomainName == h.domain {
			h.logger.Debug("got existing DNS record", "id", *record.RecordId)
			h.recordId = *record.RecordId
			return *record.Value, nil
		}
	}

	total := *result.Body.TotalCount
	if total > PerPageCount {
		time.Sleep(500 * time.Millisecond)

		total -= PerPageCount
		pages := (int)(total / PerPageCount)
		if total%PerPageCount != 0 {
			pages++
		}

		for i := 1; i <= pages; i++ {
			time.Sleep(500 * time.Millisecond)
			recordsResult, err := h.client.DescribeDomainRecords(&alidns.DescribeDomainRecordsRequest{
				DomainName: &h.domain,
				RRKeyWord:  &h.subdomain,
				Type:       utils.StringPtr(string(h.recordType)),
				Line:       &h.line,
				PageNumber: utils.Int64Ptr(int64(i)),
				PageSize:   utils.Int64Ptr(int64(pages)),
			})
			if err != nil {
				return "", err
			}

			if *recordsResult.Body.PageSize > 0 {
				for _, record := range recordsResult.Body.DomainRecords.Record {
					if *record.RR == h.subdomain && *record.DomainName == h.domain {
						h.recordId = *record.RecordId
						return *record.Value, nil
					}
				}
			}
		}
	}

	return "", nil
}

func (h *AliCloudDNSUpdateHandler) Create(address string) error {
	h.logger.Debug("creating record for address " + address)
	result, err := h.client.AddDomainRecord(&alidns.AddDomainRecordRequest{
		DomainName: &h.domain,
		RR:         &h.subdomain,
		Type:       utils.StringPtr(string(h.recordType)),
		Value:      &address,
		Line:       &h.line,
		TTL:        utils.Int64Ptr(AliCloudDefaultTTL),
	})
	if err != nil {
		return err
	}

	h.recordId = *result.Body.RecordId
	h.logger.Debug("created DNS record", "id", *result.Body.RecordId)
	return nil
}

func (h *AliCloudDNSUpdateHandler) Update(newAddress string) error {
	if h.recordId == "" {
		return fmt.Errorf("no record id present")
	}

	h.logger.Debug("updating DNS record", "id", h.recordId, "address", newAddress)
	_, err := h.client.UpdateDomainRecord(&alidns.UpdateDomainRecordRequest{
		RecordId: &h.recordId,
		RR:       &h.subdomain,
		Type:     utils.StringPtr(string(h.recordType)),
		Value:    &newAddress,
		Line:     &h.line,
		TTL:      utils.Int64Ptr(AliCloudDefaultTTL),
	})
	return err
}
