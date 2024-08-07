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
	"github.com/masteryyh/micro-ddns/internal/version"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

var jsonFormat bool

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information about micro-ddns.",
	Run: func(cmd *cobra.Command, args []string) {
		if jsonFormat {
			if err := version.PrintVersionJson(); err != nil {
				slog.Error("error when printing version in JSON: ", err)
				os.Exit(1)
			}
		} else {
			version.PrintVersionNormal()
		}
	},
}

func init() {
	versionCmd.Flags().BoolVar(&jsonFormat, "json", false, "output version in JSON format")
	rootCmd.AddCommand(versionCmd)
}
