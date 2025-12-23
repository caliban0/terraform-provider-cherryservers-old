package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var (
	_ planmodifier.String = warnIfChangedModifier{}
	_ planmodifier.Set    = warnIfChangedModifier{}
	_ planmodifier.Int64  = warnIfChangedModifier{}
)

// WarnIfChangedString returns a plan modifier that displays a warning if an attribute will be changed on update.
func WarnIfChangedString(warningSummary, warningDetail string) planmodifier.String {
	return warnIfChangedModifier{
		warningSummary: warningSummary,
		warningDetail:  warningDetail,
	}
}

// WarnIfChangedSet returns a plan modifier that displays a warning if an attribute will be changed on update.
func WarnIfChangedSet(warningSummary, warningDetail string) planmodifier.Set {
	return warnIfChangedModifier{
		warningSummary: warningSummary,
		warningDetail:  warningDetail,
	}
}

// WarnIfChangedInt64 returns a plan modifier that displays a warning if an attribute will be changed on update.
func WarnIfChangedInt64(warningSummary, warningDetail string) planmodifier.Int64 {
	return warnIfChangedModifier{
		warningSummary: warningSummary,
		warningDetail:  warningDetail,
	}
}

type warnIfChangedModifier struct {
	warningSummary string
	warningDetail  string
}

func (d warnIfChangedModifier) Description(ctx context.Context) string {
	return "Display a warning, if the attribute will be changed on update."
}

func (d warnIfChangedModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d warnIfChangedModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Ignore create or destroy cases.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Ignore if attribute has not changed.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	resp.Diagnostics.AddAttributeWarning(req.Path, d.warningSummary, d.warningDetail)
}

func (d warnIfChangedModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// Ignore create or destroy cases.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Ignore if attribute has not changed.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	resp.Diagnostics.AddAttributeWarning(req.Path, d.warningSummary, d.warningDetail)
}

func (d warnIfChangedModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	// Ignore create or destroy cases.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Ignore if attribute has not changed.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	resp.Diagnostics.AddAttributeWarning(req.Path, d.warningSummary, d.warningDetail)
}
