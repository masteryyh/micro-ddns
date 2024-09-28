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
	"log/slog"
	"time"

	alidns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/masteryyh/micro-ddns/internal/config"
	"github.com/masteryyh/micro-ddns/pkg/utils"
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

	client *alidns.Client
	logger *slog.Logger
}

func NewAliCloudDNSUpdateHandler(ddns *config.DDNSSpec, aliSpec *config.AliCloudSpec, logger *slog.Logger) (*AliCloudDNSUpdateHandler, error) {
	clientConfig := &openapi.Config{
		AccessKeyId:     &aliSpec.AccessKeyID,
		AccessKeySecret: &aliSpec.AccessKeySecret,
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

	return &AliCloudDNSUpdateHandler{
		domain:     ddns.Domain,
		subdomain:  ddns.Subdomain,
		recordType: recordType,
		line:       line,
		client:     client,
		logger:     logger,
	}, nil
}

func (h *AliCloudDNSUpdateHandler) Get(parentCtx context.Context) (string, error) {
	if h.recordId != "" {
		h.logger.Debug("record id present, getting record info")

		ctx, cancel := context.WithTimeout(parentCtx, 30*time.Second)
		defer cancel()
		result, err := utils.RunWithContext(ctx, func() (string, error) {
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
			return *result.Body.Value, nil
		})
		if err != nil {
			return "", err
		}

		val := result[0].(string)
		if result[1] != nil {
			return "", result[1].(error)
		}

		h.logger.Debug("got current ip address registered: " + val)
		return val, nil
	}

	h.logger.Debug("no record id present, searching for records already exists")

	ctx, cancel := context.WithTimeout(parentCtx, 30*time.Second)
	defer cancel()
	result, err := utils.RunWithContext(ctx, func() (string, error) {
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

		for _, record := range result.Body.DomainRecords.Record {
			if *record.RR == h.subdomain && *record.DomainName == h.domain {
				h.logger.Debug("got existing DNS record", "id", *record.RecordId)
				h.recordId = *record.RecordId
				return *record.Value, nil
			}
		}

		return "", nil
	})
	if err != nil {
		return "", err
	}

	addr := result[0].(string)
	if result[1] != nil {
		return "", result[1].(error)
	}

	if addr == "" {
		h.logger.Debug("no record with subdomain " + h.subdomain + " found")
	}
	return addr, nil
}

func (h *AliCloudDNSUpdateHandler) Create(parentCtx context.Context, address string) error {
	h.logger.Debug("creating record for address " + address)

	ctx, cancel := context.WithTimeout(parentCtx, 30*time.Second)
	defer cancel()
	result, err := utils.RunWithContext(ctx, func() (string, error) {
		result, err := h.client.AddDomainRecord(&alidns.AddDomainRecordRequest{
			DomainName: &h.domain,
			RR:         &h.subdomain,
			Type:       utils.StringPtr(string(h.recordType)),
			Value:      &address,
			Line:       &h.line,
			TTL:        utils.Int64Ptr(AliCloudDefaultTTL),
		})
		if err != nil {
			return "", err
		}
		return *result.Body.RecordId, nil
	})
	if err != nil {
		return err
	}

	val := result[0].(string)
	if result[1] != nil {
		return result[1].(error)
	}

	h.recordId = val
	h.logger.Debug("created DNS record", "id", val)
	return nil
}

func (h *AliCloudDNSUpdateHandler) Update(parentCtx context.Context, newAddress string) error {
	if h.recordId == "" {
		return fmt.Errorf("no record id present")
	}

	h.logger.Debug("updating DNS record", "id", h.recordId, "address", newAddress)

	ctx, cancel := context.WithTimeout(parentCtx, 30*time.Second)
	defer cancel()
	result, err := utils.RunWithContext(ctx, func() error {
		_, err := h.client.UpdateDomainRecord(&alidns.UpdateDomainRecordRequest{
			RecordId: &h.recordId,
			RR:       &h.subdomain,
			Type:     utils.StringPtr(string(h.recordType)),
			Value:    &newAddress,
			Line:     &h.line,
			TTL:      utils.Int64Ptr(AliCloudDefaultTTL),
		})
		return err
	})
	if err != nil {
		return err
	}

	if result[0] != nil {
		return result[0].(error)
	}
	return nil
}
