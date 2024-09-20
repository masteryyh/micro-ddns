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
	"github.com/jdcloud-api/jdcloud-sdk-go/core"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/domainservice/apis"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/domainservice/client"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/domainservice/models"
	"github.com/masteryyh/micro-ddns/internal/config"
	"github.com/masteryyh/micro-ddns/pkg/utils"
	"log/slog"
	"strconv"
)

const (
	JDCloudDefaultTTL = 120
	JDCloudPageSize   = 50
)

type JDCloudDNSUpdateHandler struct {
	domain     string
	subdomain  string
	recordType RecordType
	domainId   *int
	recordId   *int
	viewId     int

	ctx    context.Context
	cancel context.CancelFunc
	client *client.DomainserviceClient
	logger *slog.Logger
}

func NewJDCloudDNSUpdateHandler(ddns *config.DDNSSpec, spec *config.JDCloudSpec, parentCtx context.Context, logger *slog.Logger) (*JDCloudDNSUpdateHandler, error) {
	cred := core.NewCredentials(spec.AccessKey, spec.SecretKey)
	dnsClient := client.NewDomainserviceClient(cred)
	dnsClient.SetLogger(core.NewDefaultLogger(core.LogWarn))

	recordType := A
	if ddns.Stack == config.IPv6 {
		recordType = AAAA
	}

	view := -1
	if spec.ViewID != nil {
		view = *spec.ViewID
	}

	ctx, cancel := context.WithCancel(parentCtx)
	return &JDCloudDNSUpdateHandler{
		domain:     ddns.Domain,
		subdomain:  ddns.Subdomain,
		recordType: recordType,
		viewId:     view,
		ctx:        ctx,
		cancel:     cancel,
		client:     dnsClient,
		logger:     logger,
	}, nil
}

func (h *JDCloudDNSUpdateHandler) Get() (string, error) {
	if h.domainId == nil {
		h.logger.Debug("domain id is empty, searching")

		request := apis.NewDescribeDomainsRequestWithAllParams("jdcloud-api", 1, JDCloudPageSize, &h.domain, nil)
		result, err := h.client.DescribeDomains(request)
		if err != nil {
			return "", err
		}
		if result.Error.Code != 0 {
			return "", fmt.Errorf(result.Error.Message)
		}

		id := utils.IntPtr(-1)
		for _, domain := range result.Result.DataList {
			if h.domain == domain.DomainName {
				*id = domain.Id
				h.logger.Debug("got domain id " + strconv.Itoa(*id))
				break
			}
		}

		if *id == -1 {
			return "", fmt.Errorf("domain " + h.domain + " not exists")
		}
		h.domainId = id
	}

	request := apis.NewDescribeResourceRecordRequestWithAllParams("jdcloud-api", strconv.Itoa(*h.domainId), utils.IntPtr(1), utils.IntPtr(JDCloudPageSize), &h.subdomain)
	result, err := h.client.DescribeResourceRecord(request)
	if err != nil {
		return "", err
	}
	if result.Error.Code != 0 {
		return "", fmt.Errorf(result.Error.Message)
	}

	for _, record := range result.Result.DataList {
		if h.subdomain == record.HostRecord {
			h.logger.Debug("got record id " + strconv.Itoa(record.Id))
			h.recordId = &record.Id
			return record.HostValue, nil
		}
	}
	return "", nil
}

func (h *JDCloudDNSUpdateHandler) Create(address string) error {
	if h.domainId == nil {
		return fmt.Errorf("domain id is empty")
	}

	request := apis.NewCreateResourceRecordRequestWithAllParams("jdcloud-api", strconv.Itoa(*h.domainId), &models.AddRR{
		HostRecord: h.subdomain,
		HostValue:  address,
		Type:       string(h.recordType),
		ViewValue:  h.viewId,
		Ttl:        JDCloudDefaultTTL,
	})
	result, err := h.client.CreateResourceRecord(request)
	if err != nil {
		return err
	}
	if result.Error.Code != 0 {
		return fmt.Errorf(result.Error.Message)
	}

	h.recordId = &result.Result.DataList.Id
	return nil
}

func (h *JDCloudDNSUpdateHandler) Update(newAddress string) error {
	if h.domainId == nil {
		return fmt.Errorf("domain id is empty")
	}

	if h.recordId == nil {
		return fmt.Errorf("record id is empty")
	}

	request := apis.NewModifyResourceRecordRequestWithAllParams("jdcloud-api", strconv.Itoa(*h.domainId), strconv.Itoa(*h.recordId), &models.UpdateRR{
		DomainName: h.domain,
		HostRecord: h.subdomain,
		HostValue:  newAddress,
		Ttl:        JDCloudDefaultTTL,
		Type:       string(h.recordType),
		ViewValue:  h.viewId,
	})
	result, err := h.client.ModifyResourceRecord(request)
	if err != nil {
		return err
	}
	if result.Error.Code != 0 {
		return fmt.Errorf(result.Error.Message)
	}
	return nil
}
