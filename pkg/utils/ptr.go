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

package utils

func IsEmpty(v *string) bool {
	return v == nil || *v == ""
}

func StringPtrToString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func MapPtrToMap(v *map[string]string) map[string]string {
	if v == nil {
		return map[string]string{}
	}
	return *v
}

func BoolPtr(v bool) *bool {
	ptr := new(bool)
	*ptr = v
	return ptr
}

func IntPtr(v int) *int {
	ptr := new(int)
	*ptr = v
	return ptr
}

func Int32Ptr(v int32) *int32 {
	ptr := new(int32)
	*ptr = v
	return ptr
}

func Int64Ptr(v int64) *int64 {
	ptr := new(int64)
	*ptr = v
	return ptr
}

func StringPtr(v string) *string {
	ptr := new(string)
	*ptr = v
	return ptr
}

func Uint64Ptr(v uint64) *uint64 {
	ptr := new(uint64)
	*ptr = v
	return ptr
}
