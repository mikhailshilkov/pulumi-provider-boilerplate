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
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

// CustomResource is a manual SDK-based implementation of a (part of) resource when Azure API is missing some
// crucial operations.
type CustomResource struct {
	// Auxiliary types defined for this resource. Optional.
	Types map[string]schema.ComplexTypeSpec
	// Resource schema. Optional, by default the schema is assumed to be included in Azure Open API specs.
	Schema *schema.ResourceSpec
	// Create a new resource from a map of input values. Returns a map of resource outputs that match the schema shape.
	Create func(context.Context, map[string]interface{}) (string, map[string]interface{}, error)
	// Read the state of an existing resource. Constructs the resource ID based on input values. Returns a map of
	// resource outputs. If the requested resource does not exist, the second result is false.
	Read func(context.Context, map[string]interface{}) (map[string]interface{}, bool, error)
	// Update an existing resource with a map of input values. Returns a map of resource outputs that match the schema shape.
	Update func(context.Context, map[string]interface{}) (map[string]interface{}, error)
	// Delete an existing resource. Constructs the resource ID based on input values.
	Delete func(context.Context, map[string]interface{}) error
}

var Resources = map[string]*CustomResource{
	"xyz:index:RandomString": newRandomStringResource(),
}
