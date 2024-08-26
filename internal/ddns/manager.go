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
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/masteryyh/micro-ddns/internal/config"
	"log/slog"
)

type DDNSInstanceManager struct {
	instances map[string]*DDNSInstance
	specs     []*config.DDNSSpec
	ctx       context.Context
	cancel    context.CancelFunc
	scheduler gocron.Scheduler
	logger    *slog.Logger
}

func NewDDNSInstanceManager(parentCtx context.Context, specs []*config.DDNSSpec, scheduler gocron.Scheduler, logger *slog.Logger) (*DDNSInstanceManager, error) {
	ctx, cancel := context.WithCancel(parentCtx)

	if len(specs) == 0 {
		cancel()
		return nil, fmt.Errorf("no ddns specs provided")
	}

	instances := make(map[string]*DDNSInstance, len(specs))
	for _, spec := range specs {
		if _, ok := instances[spec.Name]; ok {
			cancel()
			return nil, fmt.Errorf("DDNS instance %s already exists", spec.Name)
		}

		instanceLogger := logger.With(slog.Group("component", "type", "instance", "name", spec.Name))
		instance, err := NewDDNSInstance(spec, ctx, instanceLogger)
		if err != nil {
			cancel()
			return nil, err
		}
		instances[spec.Name] = instance
	}

	return &DDNSInstanceManager{
		instances: instances,
		specs:     specs,
		ctx:       ctx,
		cancel:    cancel,
		scheduler: scheduler,
		logger:    logger,
	}, nil
}

func (m *DDNSInstanceManager) Start() error {
	for name, instance := range m.instances {
		m.logger.Info("registering DDNS task", "name", name)
		job, err := m.scheduler.NewJob(gocron.CronJob(instance.spec.Cron, false), gocron.NewTask(func(instance *DDNSInstance) {
			err := instance.DoUpdate()
			if err != nil {
				m.logger.Error("failed to handle DNS update", "name", instance.spec.Name, "err", err)
				return
			}
			m.logger.Info("successfully updated DNS record", "name", instance.spec.Name)
		}, instance))
		if err != nil {
			m.logger.Error("failed to create job", "name", name, "err", err)
			return err
		}

		m.logger.Info("created cron job", "name", name, "id", job.ID().String())
	}

	m.scheduler.Start()
	select {
	case <-m.ctx.Done():
	}

	m.logger.Info("shutting down ddns scheduler")
	return m.scheduler.Shutdown()
}
