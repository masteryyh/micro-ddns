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

import (
	"context"
	"fmt"
	"reflect"
)

func RunWithContext(ctx context.Context, function interface{}, args ...interface{}) ([]interface{}, error) {
	fn := reflect.ValueOf(function)
	if fn.Kind() != reflect.Func {
		return nil, fmt.Errorf("not a function")
	}

	input := make([]reflect.Value, len(args))
	for i, arg := range args {
		input[i] = reflect.ValueOf(arg)
	}

	resultChan := make(chan []reflect.Value)
	go func() {
		resultChan <- fn.Call(input)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case results := <-resultChan:
		returns := make([]interface{}, len(results))
		for i, v := range results {
			returns[i] = v.Interface()
		}
		return returns, nil
	}
}
