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

package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-co-op/gocron/v2"
	"github.com/masteryyh/micro-ddns/internal/config"
	"github.com/masteryyh/micro-ddns/internal/ddns"
	"github.com/masteryyh/micro-ddns/internal/metrics"
	"github.com/masteryyh/micro-ddns/internal/signal"
)

type App struct {
	manager *ddns.DDNSInstanceManager
	ctx     context.Context
	cancel  context.CancelFunc
	logger  *slog.Logger
	metrics *metrics.MetricsServer
}

func initLogger(level int) (*slog.Logger, error) {
	if level != -4 && level != 0 && level != 4 && level != 8 {
		return nil, fmt.Errorf("invalid log level: %d", level)
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(level),
	})
	logger := slog.New(handler)
	return logger, nil
}

func NewApp(logLevel int, configFile string) (*App, error) {
	ctx, cancel := signal.SetupContext()

	logger, err := initLogger(logLevel)
	if err != nil {
		return nil, err
	}

	logger.Info("reading config file from " + configFile)
	configs, err := config.ReadConfigOrGet(configFile)
	if err != nil {
		return nil, err
	}

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	manager, err := ddns.NewDDNSInstanceManager(ctx, configs.DDNS, scheduler, logger)
	if err != nil {
		return nil, err
	}

	metricsLogger := logger.With(slog.Group("component", "type", "metrics"))
	return &App{
		ctx:     ctx,
		cancel:  cancel,
		logger:  logger,
		manager: manager,
		metrics: metrics.NewMetricsServer(ctx, metricsLogger),
	}, nil
}

func (a *App) Run() {
	a.logger.Info("starting app")

	go a.metrics.Serve()
	go a.manager.Start()

	<-a.ctx.Done()
	a.logger.Info("shutting down")
}
