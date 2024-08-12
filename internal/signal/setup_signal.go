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

package signal

import (
	"context"
	"os"
	"os/signal"
)

var callOnce = make(chan struct{})

// SetupContext will receive SIGTERM and SIGINT for graceful shutdown
// returns a context for scheduler and HTTP server
func SetupContext() (context.Context, context.CancelFunc) {
	close(callOnce)

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)

	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1)
	}()

	return ctx, cancel
}
