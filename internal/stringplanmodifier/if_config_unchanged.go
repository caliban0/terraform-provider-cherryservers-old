package stringplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// IfConfigUnchanged compares the states of passed attributes to their configurations
// and sets UseState to false if, for any of them, the configuration is non-null
// and different from state. Otherwise, UseState is set to true.
func IfConfigUnchanged(exp ...path.Expression) UseStateForUnknownIfFunc {
	return func(ctx context.Context,
		req planmodifier.StringRequest,
		resp *UseStateForUnknownIfFuncResponse,
	) {
		exps := req.PathExpression.MergeExpressions(exp...)

		resp.UseState = true
		for _, e := range exps {
			// Find paths matching the expression in the configuration data.
			matchedPaths, diags := req.Config.PathMatches(ctx, e)

			resp.Diagnostics.Append(diags...)

			// Collect all errors.
			if diags.HasError() {
				continue
			}

			for _, matchedPath := range matchedPaths {
				var matchedPathConfigValue attr.Value
				var matchedPathStateValue attr.Value

				diags = req.Config.GetAttribute(ctx, matchedPath, &matchedPathConfigValue)
				resp.Diagnostics.Append(diags...)
				diags = req.State.GetAttribute(ctx, matchedPath, &matchedPathStateValue)
				resp.Diagnostics.Append(diags...)

				// Collect all errors.
				if diags.HasError() {
					continue
				}

				if !matchedPathConfigValue.IsNull() &&
					!matchedPathStateValue.Equal(matchedPathConfigValue) {
					resp.UseState = false
				}
			}

		}
	}
}
