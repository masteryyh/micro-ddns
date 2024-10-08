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

package metrics

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type MetricsServer struct {
	server *http.Server
	logger *slog.Logger
	wg     *sync.WaitGroup
}

func NewMetricsServer(logger *slog.Logger, wg *sync.WaitGroup) *MetricsServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("pong"))
	})
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	wg.Add(1)
	return &MetricsServer{
		server: server,
		logger: logger,
		wg:     wg,
	}
}

func (s *MetricsServer) Serve(parentCtx context.Context) {
	s.logger.Info("starting metrics server")
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(fmt.Sprintf("metrics server error: %v", err))
		}
	}()

	<-parentCtx.Done()
	s.logger.Info("shutting down metrics server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error(fmt.Sprintf("failed shutting down metrics server: %v", err))
	}
	s.wg.Done()
}
