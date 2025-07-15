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
	"fmt"
	"os"

	"github.com/masteryyh/micro-ddns/internal/config"
	"github.com/spf13/cobra"
)

// validateCmd validates a configuration file
var validateCmd = &cobra.Command{
	Use:   "validate [config-file]",
	Short: "Validate a configuration file",
	Long:  `Validate the syntax and content of a micro-ddns configuration file.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var configPath string
		if len(args) > 0 {
			configPath = args[0]
		} else {
			configPath = "/etc/micro-ddns/config.yaml"
		}

		fmt.Printf("Validating configuration file: %s\n", configPath)

		if err := config.InitializeConfig(configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Configuration validation failed: %v\n", err)
			return err
		}

		cfg := config.GetConfig()
		fmt.Println("✅ Configuration is valid!")
		fmt.Printf("  - DDNS specs: %d\n", len(cfg.DDNS))
		fmt.Printf("  - Detection specs: %d\n", len(cfg.Detection))
		fmt.Printf("  - Provider specs: %d\n", len(cfg.Provider))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
