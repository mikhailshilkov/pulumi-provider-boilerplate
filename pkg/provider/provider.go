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

package provider

import (
	"context"
	"fmt"
	"github.com/pulumi/pulumi-xyz/pkg/resources"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	rpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbempty "github.com/golang/protobuf/ptypes/empty"
)

type xyzProvider struct {
	host    *provider.HostClient
	name    string
	version string
}

func makeProvider(host *provider.HostClient, name, version string) (rpc.ResourceProviderServer, error) {
	// Return the new provider
	return &xyzProvider{
		host:    host,
		name:    name,
		version: version,
	}, nil
}

// CheckConfig validates the configuration for this provider.
func (p *xyzProvider) CheckConfig(ctx context.Context, req *rpc.CheckRequest) (*rpc.CheckResponse, error) {
	return &rpc.CheckResponse{Inputs: req.GetNews()}, nil
}

// DiffConfig diffs the configuration for this provider.
func (p *xyzProvider) DiffConfig(ctx context.Context, req *rpc.DiffRequest) (*rpc.DiffResponse, error) {
	return &rpc.DiffResponse{}, nil
}

// Configure configures the resource provider with "globals" that control its behavior.
func (p *xyzProvider) Configure(_ context.Context, req *rpc.ConfigureRequest) (*rpc.ConfigureResponse, error) {
	return &rpc.ConfigureResponse{}, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (p *xyzProvider) Invoke(_ context.Context, req *rpc.InvokeRequest) (*rpc.InvokeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "StreamInvoke is not yet implemented")
}

// StreamInvoke dynamically executes a built-in function in the provider. The result is streamed
// back as a series of messages.
func (p *xyzProvider) StreamInvoke(req *rpc.InvokeRequest, server rpc.ResourceProvider_StreamInvokeServer) error {
	return status.Error(codes.Unimplemented, "StreamInvoke is not yet implemented")
}

// Check validates that the given property bag is valid for a resource of the given type and returns
// the inputs that should be passed to successive calls to Diff, Create, or Update for this
// resource. As a rule, the provider inputs returned by a call to Check should preserve the original
// representation of the properties as present in the program inputs. Though this rule is not
// required for correctness, violations thereof can negatively impact the end-user experience, as
// the provider inputs are using for detecting and rendering diffs.
func (p *xyzProvider) Check(ctx context.Context, req *rpc.CheckRequest) (*rpc.CheckResponse, error) {
	typ := resource.URN(req.GetUrn()).Type()

	res, ok := resources.Resources[typ.String()]
	if !ok {
		return nil, fmt.Errorf("unknown resource type %q", typ)
	}
	if res.Create == nil {
		return nil, fmt.Errorf("resource type %q has no Create operation defined", typ)
	}
	return &rpc.CheckResponse{Inputs: req.News, Failures: nil}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (p *xyzProvider) Diff(ctx context.Context, req *rpc.DiffRequest) (*rpc.DiffResponse, error) {
	typ := resource.URN(req.GetUrn()).Type()
	res := resources.Resources[typ.String()]
	if res.Update == nil {
		var replaces []string
		for prop := range res.Schema.InputProperties {
			replaces = append(replaces, prop)
		}
		return &rpc.DiffResponse{
			Replaces: replaces,
		}, nil
	}
	return &rpc.DiffResponse{}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.
func (p *xyzProvider) Create(ctx context.Context, req *rpc.CreateRequest) (*rpc.CreateResponse, error) {
	typ := resource.URN(req.GetUrn()).Type()

	inputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{SkipNulls: true})
	if err != nil {
		return nil, err
	}
	inputsMap := inputs.Mappable()

	res := resources.Resources[typ.String()]
	id, outputsMap, err := res.Create(ctx, inputsMap)
	if err != nil {
		return nil, err
	}

	outputs, err := plugin.MarshalProperties(
		resource.NewPropertyMapFromMap(outputsMap),
		plugin.MarshalOptions{SkipNulls: true},
	)
	if err != nil {
		return nil, err
	}

	return &rpc.CreateResponse{
		Id:         id,
		Properties: outputs,
	}, nil
}

// Read the current live state associated with a resource.
func (p *xyzProvider) Read(ctx context.Context, req *rpc.ReadRequest) (*rpc.ReadResponse, error) {
	typ := resource.URN(req.GetUrn()).Type()
	id := req.GetId()

	oldState, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, err
	}

	res := resources.Resources[typ.String()]
	if res.Read == nil {
		return &rpc.ReadResponse{Id: id, Properties: req.GetProperties(), Inputs: req.GetInputs()}, nil
	}

	outputsMap, exists, err := res.Read(ctx, oldState.Mappable())
	if err != nil {
		return nil, err
	}
	if !exists {
		return &rpc.ReadResponse{Id: ""}, nil
	}

	outputs, err := plugin.MarshalProperties(
		resource.NewPropertyMapFromMap(outputsMap),
		plugin.MarshalOptions{KeepSecrets: true, KeepUnknowns: true, SkipNulls: true},
	)
	if err != nil {
		return nil, err
	}

	return &rpc.ReadResponse{Id: id, Properties: outputs, Inputs: req.GetProperties()}, nil
}

// Update updates an existing resource with new values.
func (p *xyzProvider) Update(ctx context.Context, req *rpc.UpdateRequest) (*rpc.UpdateResponse, error) {
	typ := resource.URN(req.GetUrn()).Type()

	inputs, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{SkipNulls: true})
	if err != nil {
		return nil, err
	}
	inputsMap := inputs.Mappable()

	res := resources.Resources[typ.String()]
	if res.Update == nil {
		return nil, fmt.Errorf("resource type %q has no Update operation defined", typ)
	}
	outputsMap, err := res.Update(ctx, inputsMap)
	if err != nil {
		return nil, err
	}

	outputs, err := plugin.MarshalProperties(
		resource.NewPropertyMapFromMap(outputsMap),
		plugin.MarshalOptions{SkipNulls: true},
	)
	if err != nil {
		return nil, err
	}

	return &rpc.UpdateResponse{
		Properties: outputs,
	}, nil
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed
// to still exist.
func (p *xyzProvider) Delete(ctx context.Context, req *rpc.DeleteRequest) (*pbempty.Empty, error) {
	typ := resource.URN(req.GetUrn()).Type()
	res := resources.Resources[typ.String()]
	if res.Delete != nil {
		inputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{SkipNulls: true})
		if err != nil {
			return nil, err
		}
		inputsMap := inputs.Mappable()

		err = res.Delete(ctx, inputsMap)
		if err != nil {
			return nil, err
		}
	}

	return &pbempty.Empty{}, nil
}

// Construct creates a new component resource.
func (p *xyzProvider) Construct(_ context.Context, _ *rpc.ConstructRequest) (*rpc.ConstructResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Construct is not yet implemented")
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (p *xyzProvider) GetPluginInfo(context.Context, *pbempty.Empty) (*rpc.PluginInfo, error) {
	return &rpc.PluginInfo{
		Version: p.version,
	}, nil
}

// GetSchema returns the JSON-serialized schema for the provider.
func (p *xyzProvider) GetSchema(ctx context.Context, req *rpc.GetSchemaRequest) (*rpc.GetSchemaResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetSchema is not yet implemented")
}

// Cancel signals the provider to gracefully shut down and abort any ongoing resource operations.
// Operations aborted in this way will return an error (e.g., `Update` and `Create` will either a
// creation error or an initialization error). Since Cancel is advisory and non-blocking, it is up
// to the host to decide how long to wait after Cancel is called before (e.g.)
// hard-closing any gRPC connection.
func (p *xyzProvider) Cancel(context.Context, *pbempty.Empty) (*pbempty.Empty, error) {
	return &pbempty.Empty{}, nil
}
