// *** WARNING: this file was generated by the Pulumi SDK Generator. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package xyz

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// A string of random characters of a given length.
type RandomString struct {
	pulumi.CustomResourceState

	// Length of the generated string.
	Length pulumi.IntOutput `pulumi:"length"`
	// Random string that is stored in the state and is persistent across multiple runs.
	Result pulumi.StringOutput `pulumi:"result"`
}

// NewRandomString registers a new resource with the given unique name, arguments, and options.
func NewRandomString(ctx *pulumi.Context,
	name string, args *RandomStringArgs, opts ...pulumi.ResourceOption) (*RandomString, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Length == nil {
		return nil, errors.New("invalid value for required argument 'Length'")
	}
	var resource RandomString
	err := ctx.RegisterResource("xyz:index:RandomString", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetRandomString gets an existing RandomString resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetRandomString(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *RandomStringState, opts ...pulumi.ResourceOption) (*RandomString, error) {
	var resource RandomString
	err := ctx.ReadResource("xyz:index:RandomString", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering RandomString resources.
type randomStringState struct {
	// Length of the generated string.
	Length *int `pulumi:"length"`
	// Random string that is stored in the state and is persistent across multiple runs.
	Result *string `pulumi:"result"`
}

type RandomStringState struct {
	// Length of the generated string.
	Length pulumi.IntPtrInput
	// Random string that is stored in the state and is persistent across multiple runs.
	Result pulumi.StringPtrInput
}

func (RandomStringState) ElementType() reflect.Type {
	return reflect.TypeOf((*randomStringState)(nil)).Elem()
}

type randomStringArgs struct {
	// Length of the string to generate.
	Length int `pulumi:"length"`
}

// The set of arguments for constructing a RandomString resource.
type RandomStringArgs struct {
	// Length of the string to generate.
	Length pulumi.IntInput
}

func (RandomStringArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*randomStringArgs)(nil)).Elem()
}

type RandomStringInput interface {
	pulumi.Input

	ToRandomStringOutput() RandomStringOutput
	ToRandomStringOutputWithContext(ctx context.Context) RandomStringOutput
}

func (*RandomString) ElementType() reflect.Type {
	return reflect.TypeOf((*RandomString)(nil))
}

func (i *RandomString) ToRandomStringOutput() RandomStringOutput {
	return i.ToRandomStringOutputWithContext(context.Background())
}

func (i *RandomString) ToRandomStringOutputWithContext(ctx context.Context) RandomStringOutput {
	return pulumi.ToOutputWithContext(ctx, i).(RandomStringOutput)
}

type RandomStringOutput struct {
	*pulumi.OutputState
}

func (RandomStringOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*RandomString)(nil))
}

func (o RandomStringOutput) ToRandomStringOutput() RandomStringOutput {
	return o
}

func (o RandomStringOutput) ToRandomStringOutputWithContext(ctx context.Context) RandomStringOutput {
	return o
}

func init() {
	pulumi.RegisterOutputType(RandomStringOutput{})
}
