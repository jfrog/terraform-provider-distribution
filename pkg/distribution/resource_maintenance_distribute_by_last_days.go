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
	DistributeByLastDayssEndpoint = "distribution/api/v1/maintenance/distributeByLastDays"
	DistributeByLastDayssEndpoint  = "distribution/api/v1/maintenance/distributeByLastDays"
)

func NewMaintenanceDistributeByLastDaysResource() resource.Resource {
	return &MaintenanceDistributeByLastDaysResource{
		TypeName: "distribution_distribute_by_last_days",
	}
}

type MaintenanceDistributeByLastDaysResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type MaintenanceDistributeByLastDaysResourceModel struct {
	EdgeName types.String `tfsdk:"edge_name"`
	LastNumberOfDays types.Int64 `tfsdk:"last_number_of_days"`
	Date types.String `tfsdk:"date"`
	AutoCreateMissingRepositories types.Bool `tfsdk:"auto_create_missing_repositories"`
	TotalTriggered types.Int64 `tfsdk:"total_triggered"`
	Triggered types.String `tfsdk:"triggered"`
}

type DistributeByLastDaysRequestAPIModel struct {
	EdgeName string `json:"edge_name"`
	LastNumberOfDays string `json:"last_number_of_days"`
	Date string `json:"date"`
	AutoCreateMissingRepositories string `json:"auto_create_missing_repositories"`
	TotalTriggered string `json:"total_triggered"`
	Triggered string `json:"triggered"`
}

type DistributeByLastDaysAPIModel struct {
	EdgeName string `json:"edge_name"`
	LastNumberOfDays string `json:"last_number_of_days"`
	Date string `json:"date"`
	AutoCreateMissingRepositories string `json:"auto_create_missing_repositories"`
	TotalTriggered string `json:"total_triggered"`
	Triggered string `json:"triggered"`
}

func (r *MaintenanceDistributeByLastDaysResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *MaintenanceDistributeByLastDaysResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"edge_name": schema.StringAttribute{
				Optional: true,
				Description: "The edge_name of the resource.",
			},
			"last_number_of_days": schema.StringAttribute{
				Optional: true,
				Description: "The last_number_of_days of the resource.",
			},
			"date": schema.StringAttribute{
				Optional: true,
				Description: "The date of the resource.",
			},
			"auto_create_missing_repositories": schema.StringAttribute{
				Optional: true,
				Description: "The auto_create_missing_repositories of the resource.",
			},
			"total_triggered": schema.StringAttribute{
				Optional: true,
				Description: "The total number of Release Bundle versions that were redistributed.",
			},
			"triggered": schema.StringAttribute{
				Optional: true,
				Description: "The list of Release Bundle versions that were redistributed to the specified distribution target (Edge node).",
			},
		},
		MarkdownDescription: "Manages distributeByLastDays in JFrog Distribution.",
	}
}

func (r *MaintenanceDistributeByLastDaysResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}


func (r *MaintenanceDistributeByLastDaysResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	go util.SendUsageResourceCreate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan MaintenanceDistributeByLastDaysResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := MaintenanceDistributeByLastDaysRequestAPIModel{
		EdgeName: plan.EdgeName.ValueString(),
		LastNumberOfDays: plan.LastNumberOfDays.ValueString(),
		Date: plan.Date.ValueString(),
		AutoCreateMissingRepositories: plan.AutoCreateMissingRepositories.ValueString(),
		TotalTriggered: plan.TotalTriggered.ValueString(),
		Triggered: plan.Triggered.ValueString(),
	}

	var result MaintenanceDistributeByLastDaysAPIModel

	response, err := r.ProviderData.Client.R().
		SetBody(requestBody).
		SetResult(&result).
		Post(DistributeByLastDayssEndpoint)
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


func (r *MaintenanceDistributeByLastDaysResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state MaintenanceDistributeByLastDaysResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result MaintenanceDistributeByLastDaysAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
		}).
		SetResult(&result).
		Get(DistributeByLastDayssEndpoint)
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




