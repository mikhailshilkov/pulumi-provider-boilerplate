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

package main

import (
	"encoding/json"
	"fmt"
	"github.com/pulumi/pulumi-xyz/pkg/resources"
	"os"
	"path"

	"github.com/pulumi/pulumi/sdk/v3/go/common/tools"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"

	"github.com/pkg/errors"
	dotnetgen "github.com/pulumi/pulumi/pkg/v3/codegen/dotnet"
	gogen "github.com/pulumi/pulumi/pkg/v3/codegen/go"
	nodejsgen "github.com/pulumi/pulumi/pkg/v3/codegen/nodejs"
	pygen "github.com/pulumi/pulumi/pkg/v3/codegen/python"
	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: pulumi-sdkgen-xyz <target-sdk-folder>\n")
		return
	}

	targetSdkFolder := os.Args[1]

	err := emitPackage(targetSdkFolder)
	if err != nil {
		fmt.Printf("Failed: %s", err.Error())
	}
}

// emitPackage emits an entire package pack into the configured output directory with the configured settings.
func emitPackage(targetSdkFolder string) error {
	spec := pschema.PackageSpec{
		Name:      "xyz",
		Resources: map[string]pschema.ResourceSpec{},
		Language: map[string]json.RawMessage{
			"nodejs": rawMessage(map[string]interface{}{
				"dependencies": map[string]string{
					"@pulumi/pulumi": "^3.0.0",
				},
			}),
			"python": rawMessage(map[string]interface{}{
				"usesIOClasses": true,
			}),
			"csharp": rawMessage(map[string]interface{}{
				"packageReferences": map[string]string{
					"Pulumi":                       "3.*",
					"System.Collections.Immutable": "1.6.0",
				},
			}),
			"go": rawMessage(map[string]interface{}{}),
		},
	}

	for tok, res := range resources.Resources {
		spec.Resources[tok] = *res.Schema
	}

	ppkg, err := pschema.ImportSpec(spec, nil)
	if err != nil {
		return errors.Wrap(err, "reading schema")
	}

	toolDescription := "the Pulumi SDK Generator"
	extraFiles := map[string][]byte{}

	sdkGenerators := map[string]func() (map[string][]byte, error){
		"python": func() (map[string][]byte, error) {
			return pygen.GeneratePackage(toolDescription, ppkg, extraFiles)
		},
		"nodejs": func() (map[string][]byte, error) {
			return nodejsgen.GeneratePackage(toolDescription, ppkg, extraFiles)
		},
		"go": func() (map[string][]byte, error) {
			return gogen.GeneratePackage(toolDescription, ppkg)
		},
		"dotnet": func() (map[string][]byte, error) {
			return dotnetgen.GeneratePackage(toolDescription, ppkg, extraFiles)
		},
	}

	for sdkName, generator := range sdkGenerators {
		files, err := generator()
		if err != nil {
			return errors.Wrapf(err, "generating %s package", sdkName)
		}

		for f, contents := range files {
			if err := emitFile(path.Join(targetSdkFolder, sdkName), f, contents); err != nil {
				return errors.Wrapf(err, "emitting file %v", f)
			}
		}
	}

	return nil
}

func emitFile(outDir, relPath string, contents []byte) error {
	p := path.Join(outDir, relPath)
	if err := tools.EnsureDir(path.Dir(p)); err != nil {
		return errors.Wrap(err, "creating directory")
	}

	f, err := os.Create(p)
	if err != nil {
		return errors.Wrap(err, "creating file")
	}
	defer contract.IgnoreClose(f)

	_, err = f.Write(contents)
	return err
}

func rawMessage(v interface{}) json.RawMessage {
	bytes, err := json.Marshal(v)
	contract.Assert(err == nil)
	return bytes
}
