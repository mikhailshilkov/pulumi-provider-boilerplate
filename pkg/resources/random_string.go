// Copyright 2016-2021, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"context"
	"fmt"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"math/rand"
	"strconv"
	"time"
)

func newRandomStringResource() *CustomResource {
	return &CustomResource{
		Schema: &schema.ResourceSpec{
			ObjectTypeSpec: schema.ObjectTypeSpec{
				Description: "A string of random characters of a given length.",
				Type:        "object",
				Properties: map[string]schema.PropertySpec{
					"length": {
						Description: "Length of the generated string.",
						TypeSpec: schema.TypeSpec{Type: "integer"},
					},
					"result": {
						Description: "Random string that is stored in the state and is persistent across multiple runs.",
						TypeSpec: schema.TypeSpec{Type: "string"},
					},
				},
				Required: []string{"length", "result"},
			},
			InputProperties: map[string]schema.PropertySpec{
				"length": {
					Description: "Length of the string to generate.",
					TypeSpec: schema.TypeSpec{Type: "integer"},
				},
			},
			RequiredInputs: []string{"length"},
		},
		Create: create,
	}
}

func create(_ context.Context, inputs map[string]interface{}) (string, map[string]interface{}, error) {
	length, ok := inputs["length"].(float64)
	if !ok {
		return "", nil, fmt.Errorf("expected input property 'length' of type 'number' but got '%s", inputs["length"])
	}

	n := int(length)

	// Actually "create" the random string.
	result := makeRandom(n)

	outputs := map[string]interface{}{
		"length": n,
		"result": result,
	}

	// Id defines the identity of the resource: pick stable properties to build it.
	id := strconv.Itoa(n)

	return id, outputs, nil
}

func makeRandom(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	result := make([]rune, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}
