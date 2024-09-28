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

package ip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/itchyny/gojq"
	"github.com/masteryyh/micro-ddns/internal/config"
	"github.com/masteryyh/micro-ddns/pkg/utils"
)

type ThirdPartyAddressDetector struct {
	url                string
	jsonPath           string
	params             map[string]string
	headers            map[string]string
	username           string
	password           string
	localAddressPolicy config.LocalAddressPolicy
	stack              config.NetworkStack

	logger *slog.Logger
}

func NewThirdPartyAddressDetector(detectionSpec config.AddressDetectionSpec, stack config.NetworkStack, logger *slog.Logger) *ThirdPartyAddressDetector {
	spec := detectionSpec.API

	var policy config.LocalAddressPolicy
	if detectionSpec.LocalAddressPolicy == nil {
		policy = config.LocalAddressPolicyIgnore
	} else {
		policy = *detectionSpec.LocalAddressPolicy
	}
	return &ThirdPartyAddressDetector{
		url:                spec.URL,
		jsonPath:           utils.StringPtrToString(spec.JsonPath),
		params:             utils.MapPtrToMap(spec.Params),
		headers:            utils.MapPtrToMap(spec.Headers),
		username:           utils.StringPtrToString(spec.Username),
		password:           utils.StringPtrToString(spec.Password),
		localAddressPolicy: policy,
		stack:              stack,
		logger:             logger,
	}
}

func (d *ThirdPartyAddressDetector) requestAddress(parentCtx context.Context) (string, error) {
	client := http.DefaultClient
	params := url.Values{}
	for k, v := range d.params {
		params.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(parentCtx, 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.url, nil)
	if err != nil {
		return "", err
	}
	for k, v := range d.headers {
		req.Header.Set(k, v)
	}

	d.logger.Debug("requesting address " + req.URL.String())
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	header := res.Header.Get("Content-Type")
	if header == "" {
		header = res.Header.Get("content-type")
	}

	if strings.Contains(header, "application/json") || d.jsonPath != "" {
		d.logger.Debug("extracting address from JSON data using jsonpath", "jsonpath", d.jsonPath)
		if d.jsonPath == "" {
			return "", fmt.Errorf("no jsonpath specified")
		}
		return d.extractIP(parentCtx, body, d.jsonPath)
	}

	d.logger.Debug("use body as address directly", "body", body)
	return string(body), nil
}

func (d *ThirdPartyAddressDetector) extractIP(parentCtx context.Context, response []byte, jsonpath string) (string, error) {
	if len(response) == 0 {
		return "", fmt.Errorf("response is empty")
	}

	if len(jsonpath) == 0 {
		return "", fmt.Errorf("jsonpath is empty")
	}

	var body map[string]interface{}
	if err := json.Unmarshal(response, &body); err != nil {
		return "", err
	}

	query, err := gojq.Parse(jsonpath)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(parentCtx, 30*time.Second)
	defer cancel()
	iter := query.RunWithContext(ctx, body)
	var val string
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := v.(error); ok {
			var haltError *gojq.HaltError
			if errors.As(err, &haltError) && haltError.Value() == nil {
				break
			}
			return "", err
		}
		val = v.(string)
	}

	return val, nil
}

func (d *ThirdPartyAddressDetector) detectV4(parentCtx context.Context) (string, error) {
	val, err := d.requestAddress(parentCtx)
	if err != nil {
		return "", err
	}

	if !IsValidV4(val) {
		return "", fmt.Errorf("invalid address: %s", val)
	}

	if IsPrivate(val) {
		if d.localAddressPolicy == config.LocalAddressPolicyIgnore {
			return "", fmt.Errorf("local address is ignored: %s", val)
		}
	}

	return val, nil
}

func (d *ThirdPartyAddressDetector) detectV6(parentCtx context.Context) (string, error) {
	val, err := d.requestAddress(parentCtx)
	if err != nil {
		return "", err
	}

	if !IsValidV6(val) {
		return "", fmt.Errorf("invalid address: %s", val)
	}

	if IsPrivate(val) {
		if d.localAddressPolicy == config.LocalAddressPolicyIgnore {
			return "", fmt.Errorf("ULA address is ignored: %s", val)
		}
	}

	return val, nil
}

func (d *ThirdPartyAddressDetector) Detect(parentCtx context.Context) (string, error) {
	if d.stack == config.IPv6 {
		return d.detectV6(parentCtx)
	}
	return d.detectV4(parentCtx)
}
