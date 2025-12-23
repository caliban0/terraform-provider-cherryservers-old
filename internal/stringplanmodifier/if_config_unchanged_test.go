package stringplanmodifier

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestIfConfigUnchangedAtomic(t *testing.T) {
	t.Parallel()

	const testAttr = "test_attr"

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			testAttr: schema.StringAttribute{},
		},
	}

	testConfig := func(value types.String) tfsdk.Config {
		tfValue, err := value.ToTerraformValue(context.Background())
		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.Config{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(context.Background()),
				map[string]tftypes.Value{
					testAttr: tfValue,
				},
			),
		}
	}

	testState := func(value types.String) tfsdk.State {
		tfValue, err := value.ToTerraformValue(context.Background())
		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.State{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(context.Background()),
				map[string]tftypes.Value{
					testAttr: tfValue,
				},
			),
		}
	}

	testCases := map[string]struct {
		req  planmodifier.StringRequest
		want *UseStateForUnknownIfFuncResponse
	}{
		"state-different-from-config": {
			req: planmodifier.StringRequest{
				Config: testConfig(types.StringValue("a")),
				State:  testState(types.StringValue("b")),
			},
			want: &UseStateForUnknownIfFuncResponse{
				UseState:    false,
				Diagnostics: nil,
			},
		},
		"state-null-config-non-null": {
			req: planmodifier.StringRequest{
				Config: testConfig(types.StringValue("a")),
				State:  testState(types.StringNull()),
			},
			want: &UseStateForUnknownIfFuncResponse{
				UseState:    false,
				Diagnostics: nil,
			},
		},
		"state-empty-and-different-from-config": {
			req: planmodifier.StringRequest{
				Config: testConfig(types.StringValue("a")),
				State:  testState(types.StringValue("")),
			},
			want: &UseStateForUnknownIfFuncResponse{
				UseState:    false,
				Diagnostics: nil,
			},
		},
		"config-unknown": {
			req: planmodifier.StringRequest{
				Config: testConfig(types.StringUnknown()),
				State:  testState(types.StringValue("b")),
			},
			want: &UseStateForUnknownIfFuncResponse{
				UseState:    false,
				Diagnostics: nil,
			},
		},
		"state-same-as-config": {
			req: planmodifier.StringRequest{
				Config: testConfig(types.StringValue("a")),
				State:  testState(types.StringValue("a")),
			},
			want: &UseStateForUnknownIfFuncResponse{
				UseState:    true,
				Diagnostics: nil,
			},
		},
		"config-null-and-state-known": {
			req: planmodifier.StringRequest{
				Config: testConfig(types.StringNull()),
				State:  testState(types.StringValue("a")),
			},
			want: &UseStateForUnknownIfFuncResponse{
				UseState:    true,
				Diagnostics: nil,
			},
		},
		"state-null-config-null": {
			req: planmodifier.StringRequest{
				Config: testConfig(types.StringNull()),
				State:  testState(types.StringNull()),
			},
			want: &UseStateForUnknownIfFuncResponse{
				UseState:    true,
				Diagnostics: nil,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			f := IfConfigUnchanged(path.MatchRoot(testAttr))
			got := &UseStateForUnknownIfFuncResponse{}
			f(context.Background(), testCase.req, got)

			if diff := cmp.Diff(testCase.want, got); diff != "" {
				t.Errorf("unexpected difference %s", diff)
			}
		})
	}
}

func TestIfConfigUnchangedComplex(t *testing.T) {
	t.Parallel()

	const testAttr = "test_attr"

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			testAttr: schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						testAttr: schema.StringAttribute{},
					},
				},
			},
		},
	}

	testConfig := func(value types.Set) tfsdk.Config {
		tfValue, err := value.ToTerraformValue(context.Background())
		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.Config{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(context.Background()),
				map[string]tftypes.Value{
					testAttr: tfValue,
				},
			),
		}
	}

	testState := func(value types.Set) tfsdk.State {
		tfValue, err := value.ToTerraformValue(context.Background())
		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.State{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(context.Background()),
				map[string]tftypes.Value{
					testAttr: tfValue,
				},
			),
		}
	}

	objType := types.ObjectNull(map[string]attr.Type{testAttr: types.StringType}).Type(context.Background())

	testCases := map[string]struct {
		req  planmodifier.StringRequest
		want *UseStateForUnknownIfFuncResponse
	}{
		"child-state-different-from-config": {
			req: planmodifier.StringRequest{
				Config: testConfig(
					types.SetValueMust(
						objType,
						[]attr.Value{types.ObjectValueMust(
							map[string]attr.Type{testAttr: types.StringType},
							map[string]attr.Value{testAttr: types.StringValue("a")})})),
				State: testState(
					types.SetValueMust(
						objType,
						[]attr.Value{types.ObjectValueMust(
							map[string]attr.Type{testAttr: types.StringType},
							map[string]attr.Value{testAttr: types.StringValue("b")})})),
			},
			want: &UseStateForUnknownIfFuncResponse{
				UseState:    false,
				Diagnostics: nil,
			},
		},
		"config-has-no-elements-but-state-does": {
			req: planmodifier.StringRequest{
				Config: testConfig(
					types.SetValueMust(
						objType,
						[]attr.Value{})),
				State: testState(
					types.SetValueMust(
						objType,
						[]attr.Value{types.ObjectValueMust(
							map[string]attr.Type{testAttr: types.StringType},
							map[string]attr.Value{testAttr: types.StringValue("b")})})),
			},
			want: &UseStateForUnknownIfFuncResponse{
				UseState:    false,
				Diagnostics: nil,
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			f := IfConfigUnchanged(path.MatchRoot(testAttr))
			got := &UseStateForUnknownIfFuncResponse{}
			f(context.Background(), testCase.req, got)

			if diff := cmp.Diff(testCase.want, got); diff != "" {
				t.Errorf("unexpected difference %s", diff)
			}
		})
	}
}
