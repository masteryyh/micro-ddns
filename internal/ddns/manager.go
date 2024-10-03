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
	"log/slog"
	"sync"

	"github.com/go-co-op/gocron/v2"
	"github.com/masteryyh/micro-ddns/internal/config"
)

type DDNSInstanceManager struct {
	instances map[string]*DDNSInstance
	specs     []*config.DDNSSpec
	scheduler gocron.Scheduler
	logger    *slog.Logger
	wg        *sync.WaitGroup
}

func NewDDNSInstanceManager(specs []*config.DDNSSpec, scheduler gocron.Scheduler, logger *slog.Logger, wg *sync.WaitGroup) (*DDNSInstanceManager, error) {
	instances := make(map[string]*DDNSInstance, len(specs))
	for _, spec := range specs {
		if _, ok := instances[spec.Name]; ok {
			return nil, fmt.Errorf("DDNS instance %s already exists", spec.Name)
		}

		instanceLogger := logger.With(slog.Group("component", "type", "instance", "name", spec.Name))
		instance, err := NewDDNSInstance(spec, instanceLogger)
		if err != nil {
			return nil, err
		}
		instances[spec.Name] = instance
	}

	wg.Add(1)
	return &DDNSInstanceManager{
		instances: instances,
		specs:     specs,
		scheduler: scheduler,
		logger:    logger,
		wg:        wg,
	}, nil
}

func (m *DDNSInstanceManager) Start(parentCtx context.Context) {
	for name, instance := range m.instances {
		m.logger.Info("registering DDNS task", "name", name)
		job, err := m.scheduler.NewJob(gocron.CronJob(instance.spec.Cron, false), gocron.NewTask(func(ctx context.Context, instance *DDNSInstance) {
			err := instance.DoUpdate(ctx)
			if err != nil {
				m.logger.Error("failed to handle DNS update", "name", instance.spec.Name, "err", err)
				return
			}
			m.logger.Info("successfully updated DNS record", "name", instance.spec.Name)
		}, parentCtx, instance))
		if err != nil {
			m.logger.Error("failed to create job", "name", name, "err", err)
			return
		}

		m.logger.Info("created cron job", "name", name, "id", job.ID().String())
	}

	m.scheduler.Start()

	<-parentCtx.Done()
	m.logger.Info("shutting down ddns scheduler")
	if err := m.scheduler.Shutdown(); err != nil {
		m.logger.Error(fmt.Sprintf("failed shutting down ddns scheduler: %v", err))
	}
	m.wg.Done()
}
