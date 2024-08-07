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

package version

import (
	"encoding/json"
	"fmt"
)

type AppVersionInfo struct {
	Version    string `json:"Version"`
	BuildTime  string `json:"BuildTime"`
	GoVersion  string `json:"GoVersion"`
	CommitHash string `json:"CommitHash"`
}

func PrintVersionNormal() {
	fmt.Printf("micro-ddns version: %s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Go Version: %s\n", GoVersion)
	fmt.Printf("Latest Commit Hash: %s\n", CommitHash)
}

func PrintVersionJson() error {
	ver := AppVersionInfo{
		Version:    Version,
		BuildTime:  BuildTime,
		GoVersion:  GoVersion,
		CommitHash: CommitHash,
	}

	bytes, err := json.Marshal(ver)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}
