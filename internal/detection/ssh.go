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

package detection

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/masteryyh/micro-ddns/internal/config"
	"github.com/masteryyh/micro-ddns/pkg/utils"
	"golang.org/x/crypto/ssh"
)

type SSHAddressDetector struct {
	address    string
	port       uint16
	user       string
	authMethod ssh.AuthMethod
	command    string

	stack    config.NetworkStack
	includes []*net.IPNet
	excludes []*net.IPNet
	logger   *slog.Logger
}

func NewSSHAddressDetector(detectionSpec *config.AddressDetectionSpec, stack config.NetworkStack, logger *slog.Logger) (*SSHAddressDetector, error) {
	spec := detectionSpec.SSH

	var command string
	if !utils.IsEmpty(spec.Command) {
		command = *spec.Command
	} else {
		if stack == config.IPv4 {
			command = fmt.Sprintf("ip addr show dev %s | grep -oE 'inet ([^/]+)' | awk '{print $2}'", *detectionSpec.SSH.Interface)
		} else {
			command = fmt.Sprintf("ip -6 addr show dev %s | grep -oE 'inet6 ([^/]+)' | awk '{print $2}'", *detectionSpec.SSH.Interface)
		}
	}

	var port uint16 = 22
	if spec.Host.Port != nil {
		port = *spec.Host.Port
	}

	var auth ssh.AuthMethod
	if !utils.IsEmpty(spec.Credential.PrivateKey) {
		var signer ssh.Signer
		if !utils.IsEmpty(spec.Credential.Passphrase) {
			s, err := ssh.ParsePrivateKeyWithPassphrase([]byte(*spec.Credential.PrivateKey), []byte(*spec.Credential.Passphrase))
			if err != nil {
				return nil, fmt.Errorf("failed to parse private key: %w", err)
			}
			signer = s
		} else {
			s, err := ssh.ParsePrivateKey([]byte(*spec.Credential.PrivateKey))
			if err != nil {
				return nil, fmt.Errorf("failed to parse private key: %w", err)
			}
			signer = s
		}
		auth = ssh.PublicKeys(signer)
	} else {
		auth = ssh.Password(*spec.Credential.Password)
	}

	includes := []*net.IPNet{}
	if detectionSpec.Selector != nil {
		includes = detectionSpec.Selector.GetIncludes()
	}

	excludes := []*net.IPNet{}
	if detectionSpec.Selector != nil {
		excludes = detectionSpec.Selector.GetExcludes()
	}

	return &SSHAddressDetector{
		address:    spec.Host.Address,
		port:       port,
		user:       spec.Credential.User,
		authMethod: auth,
		command:    command,
		stack:      stack,
		includes:   includes,
		excludes:   excludes,
		logger:     logger,
	}, nil
}

func (d *SSHAddressDetector) Detect(ctx context.Context) (string, error) {
	sshCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	results, err := utils.RunWithContext(sshCtx, func() (string, error) {
		d.logger.Debug("creating SSH connection to remote host", "address", d.address, "port", d.port)
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", d.address, d.port), &ssh.ClientConfig{
			User:            d.user,
			Auth:            []ssh.AuthMethod{d.authMethod},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		if err != nil {
			return "", fmt.Errorf("failed to dial SSH: %w", err)
		}
		defer client.Close()

		session, err := client.NewSession()
		if err != nil {
			return "", fmt.Errorf("failed to create SSH session: %w", err)
		}
		defer session.Close()

		d.logger.Debug("executing command", "command", d.command)
		output, err := session.Output(d.command)
		if err != nil {
			return "", fmt.Errorf("failed to execute command: %w", err)
		}
		return string(output), nil
	})
	if err != nil {
		return "", err
	}

	err, ok := results[1].(error)
	if ok && err != nil {
		return "", err
	}

	addressResult := results[0].(string)
	addresses := strings.Split(addressResult, "\n")
	for _, addr := range addresses {
		if addr == "" {
			continue
		}

		d.logger.Debug("checking address", "address", addr)
		var ip net.IP
		if d.stack == config.IPv4 {
			ip = IsValidV4(addr)
		} else {
			ip = IsValidV6(addr)
		}

		if ip != nil && !addressExcluded(ip, d.includes, d.excludes) {
			d.logger.Debug("valid address found", "address", addr)
			return addr, nil
		}
	}
	return "", fmt.Errorf("no valid address found")
}
