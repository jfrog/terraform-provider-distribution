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
	StopAllsEndpoint = "distribution/api/v1/maintenance/execute/stop_all"
	StopAllsEndpoint  = "distribution/api/v1/maintenance/execute/stop_all"
)

func NewExecuteStopAllResource() resource.Resource {
	return &ExecuteStopAllResource{
		TypeName: "distribution_stop_all",
	}
}

type ExecuteStopAllResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type ExecuteStopAllResourceModel struct {
	TrackersDbInfo types.Map `tfsdk:"trackers_db_info"`
	AffectedEdgeTrackerArtifacts types.Int64 `tfsdk:"affected_edge_tracker_artifacts"`
	AffectedEdgeTrackers types.Int64 `tfsdk:"affected_edge_trackers"`
	AffectedTrackers types.Int64 `tfsdk:"affected_trackers"`
	TasksQueueFlushed types.Bool `tfsdk:"tasks_queue_flushed"`
	DryRun types.String `tfsdk:"dryRun"`
}

type StopAllRequestAPIModel struct {
	TrackersDbInfo string `json:"trackers_db_info"`
	AffectedEdgeTrackerArtifacts string `json:"affected_edge_tracker_artifacts"`
	AffectedEdgeTrackers string `json:"affected_edge_trackers"`
	AffectedTrackers string `json:"affected_trackers"`
	TasksQueueFlushed string `json:"tasks_queue_flushed"`
	DryRun string `json:"dryRun"`
}

type StopAllAPIModel struct {
	TrackersDbInfo string `json:"trackers_db_info"`
	AffectedEdgeTrackerArtifacts string `json:"affected_edge_tracker_artifacts"`
	AffectedEdgeTrackers string `json:"affected_edge_trackers"`
	AffectedTrackers string `json:"affected_trackers"`
	TasksQueueFlushed string `json:"tasks_queue_flushed"`
	DryRun string `json:"dryRun"`
}

func (r *ExecuteStopAllResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *ExecuteStopAllResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"trackers_db_info": schema.StringAttribute{
				Optional: true,
				Description: "The trackers_db_info of the resource.",
			},
			"affected_edge_tracker_artifacts": schema.StringAttribute{
				Optional: true,
				Description: "The affected_edge_tracker_artifacts of the resource.",
			},
			"affected_edge_trackers": schema.StringAttribute{
				Optional: true,
				Description: "The affected_edge_trackers of the resource.",
			},
			"affected_trackers": schema.StringAttribute{
				Optional: true,
				Description: "The affected_trackers of the resource.",
			},
			"tasks_queue_flushed": schema.StringAttribute{
				Optional: true,
				Description: "The tasks_queue_flushed of the resource.",
			},
			"dryRun": schema.StringAttribute{
				Required: true,
				Description: "If true, only parses and validates.",
			},
		},
		MarkdownDescription: "Manages stop_all in JFrog Distribution.",
	}
}

func (r *ExecuteStopAllResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}


func (r *ExecuteStopAllResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	go util.SendUsageResourceCreate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan ExecuteStopAllResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := ExecuteStopAllRequestAPIModel{
		TrackersDbInfo: plan.TrackersDbInfo.ValueString(),
		AffectedEdgeTrackerArtifacts: plan.AffectedEdgeTrackerArtifacts.ValueString(),
		AffectedEdgeTrackers: plan.AffectedEdgeTrackers.ValueString(),
		AffectedTrackers: plan.AffectedTrackers.ValueString(),
		TasksQueueFlushed: plan.TasksQueueFlushed.ValueString(),
		DryRun: plan.DryRun.ValueString(),
	}

	var result ExecuteStopAllAPIModel

	response, err := r.ProviderData.Client.R().
		SetBody(requestBody).
		SetResult(&result).
		Post(StopAllsEndpoint)
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


func (r *ExecuteStopAllResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state ExecuteStopAllResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result ExecuteStopAllAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
		}).
		SetResult(&result).
		Get(StopAllsEndpoint)
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




