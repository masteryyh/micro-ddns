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
	"github.com/go-co-op/gocron/v2"
	"github.com/masteryyh/micro-ddns/internal/config"
	"github.com/masteryyh/micro-ddns/internal/ddns"
	"github.com/masteryyh/micro-ddns/internal/signal"
	"log/slog"
	"os"
)

type App struct {
	manager *ddns.DDNSInstanceManager
	ctx     context.Context
	cancel  context.CancelFunc
	logger  *slog.Logger

	configFile string
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

	return &App{
		ctx:        ctx,
		cancel:     cancel,
		logger:     logger,
		configFile: configFile,
	}, nil
}

func (a *App) Run() error {
	a.logger.Info("starting app")
	a.logger.Info("reading config file from " + a.configFile)
	configs, err := config.ReadConfigOrGet(a.configFile)
	if err != nil {
		return err
	}

	a.logger.Info("starting cron task scheduler")
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	manager, err := ddns.NewDDNSInstanceManager(a.ctx, configs.DDNS, scheduler, a.logger)
	if err != nil {
		return err
	}

	return manager.Start()
}
