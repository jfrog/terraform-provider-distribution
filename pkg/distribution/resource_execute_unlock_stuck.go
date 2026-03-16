package distribution

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jfrog/terraform-provider-shared/util"
	utilfw "github.com/jfrog/terraform-provider-shared/util/fw"
)

const (
	UnlockStucksEndpoint = "distribution/api/v1/maintenance/execute/unlock_stuck"
	UnlockStucksEndpoint  = "distribution/api/v1/maintenance/execute/unlock_stuck"
)

func NewExecuteUnlockStuckResource() resource.Resource {
	return &ExecuteUnlockStuckResource{
		TypeName: "distribution_unlock_stuck",
	}
}

type ExecuteUnlockStuckResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type ExecuteUnlockStuckResourceModel struct {
	StuckDistributions types.String `tfsdk:"stuck_distributions"`
	BundleName types.String `tfsdk:"bundle_name"`
	BundleVersion types.String `tfsdk:"bundle_version"`
	ReleaseBundlesSchema types.String `tfsdk:"release_bundles_schema"`
	DistributionFriendlyId types.String `tfsdk:"distribution_friendly_id"`
	Status types.String `tfsdk:"status"`
	DistributionTrackerId types.String `tfsdk:"distribution_tracker_id"`
	DryRun types.String `tfsdk:"dryRun"`
}

type UnlockStuckRequestAPIModel struct {
	StuckDistributions string `json:"stuck_distributions"`
	BundleName string `json:"bundle_name"`
	BundleVersion string `json:"bundle_version"`
	ReleaseBundlesSchema string `json:"release_bundles_schema"`
	DistributionFriendlyId string `json:"distribution_friendly_id"`
	Status string `json:"status"`
	DistributionTrackerId string `json:"distribution_tracker_id"`
	DryRun string `json:"dryRun"`
}

type UnlockStuckAPIModel struct {
	StuckDistributions string `json:"stuck_distributions"`
	BundleName string `json:"bundle_name"`
	BundleVersion string `json:"bundle_version"`
	ReleaseBundlesSchema string `json:"release_bundles_schema"`
	DistributionFriendlyId string `json:"distribution_friendly_id"`
	Status string `json:"status"`
	DistributionTrackerId string `json:"distribution_tracker_id"`
	DryRun string `json:"dryRun"`
}

func (r *ExecuteUnlockStuckResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *ExecuteUnlockStuckResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"stuck_distributions": schema.StringAttribute{
				Optional: true,
				Description: "The stuck_distributions of the resource.",
			},
			"bundle_name": schema.StringAttribute{
				Optional: true,
				Description: "The bundle_name of the resource.",
			},
			"bundle_version": schema.StringAttribute{
				Optional: true,
				Description: "The bundle_version of the resource.",
			},
			"release_bundles_schema": schema.StringAttribute{
				Optional: true,
				Description: "The release_bundles_schema of the resource.",
			},
			"distribution_friendly_id": schema.StringAttribute{
				Optional: true,
				Description: "The distribution_friendly_id of the resource.",
			},
			"status": schema.StringAttribute{
				Optional: true,
				Description: "The status of the resource.",
			},
			"distribution_tracker_id": schema.StringAttribute{
				Optional: true,
				Description: "The distribution_tracker_id of the resource.",
			},
			"dryRun": schema.StringAttribute{
				Required: true,
				Description: "If set to true , parses, validates, and returns the list of tasks that would be unlocked, but does not execute the operation.",
			},
		},
		MarkdownDescription: "Manages unlock_stuck in JFrog Distribution.",
	}
}

func (r *ExecuteUnlockStuckResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}


func (r *ExecuteUnlockStuckResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	go util.SendUsageResourceCreate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan ExecuteUnlockStuckResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := ExecuteUnlockStuckRequestAPIModel{
		StuckDistributions: plan.StuckDistributions.ValueString(),
		BundleName: plan.BundleName.ValueString(),
		BundleVersion: plan.BundleVersion.ValueString(),
		ReleaseBundlesSchema: plan.ReleaseBundlesSchema.ValueString(),
		DistributionFriendlyId: plan.DistributionFriendlyId.ValueString(),
		Status: plan.Status.ValueString(),
		DistributionTrackerId: plan.DistributionTrackerId.ValueString(),
		DryRun: plan.DryRun.ValueString(),
	}

	var result ExecuteUnlockStuckAPIModel

	response, err := r.ProviderData.Client.R().
		SetBody(requestBody).
		SetResult(&result).
		Post(UnlockStucksEndpoint)
	if err != nil {
		utilfw.UnableToCreateResourceError(resp, err.Error())
		return
	}

	if response.IsError() {
		utilfw.UnableToCreateResourceError(resp, response.String())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}


func (r *ExecuteUnlockStuckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state ExecuteUnlockStuckResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result ExecuteUnlockStuckAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
		}).
		SetResult(&result).
		Get(UnlockStucksEndpoint)
	if err != nil {
		utilfw.UnableToRefreshResourceError(resp, err.Error())
		return
	}

	if response.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	if response.IsError() {
		utilfw.UnableToRefreshResourceError(resp, response.String())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}




