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

package dns

type RecordType string

const (
	A    RecordType = "A"
	AAAA RecordType = "AAAA"

	PerPageCount = 500
)

type DNSUpdateHandler interface {
	// Get will get current IP address registered in DNS record
	Get() (string, error)

	// Create will create new DNS record with address
	Create(address string) error

	// Update will update DNS record with new address
	Update(newAddress string) error
}
