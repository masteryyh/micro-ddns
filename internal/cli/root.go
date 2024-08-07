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

package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	logLevel int
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "micro-ddns",
		Short: "A simple tool to update DNS dynamically",
		Long: `A simple tool to update DNS dynamically if you want to create DNS records
for hosts that IP address might changes. Support multiple DNS providers.

Still under construction. ;)`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&logLevel, "verbose", "v",
		4, "log level, available options are -4 (DEBUG), 0 (INFO), 4 (WARN) and 8 (ERROR)")
}
