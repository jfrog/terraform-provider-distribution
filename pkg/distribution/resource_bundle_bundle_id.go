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
	BundleIdsEndpoint = "distribution/api/v1/system/support/bundle/bundle_id"
	BundleIdsEndpoint  = "distribution/api/v1/system/support/bundle/bundle_id"
)

func NewBundleBundleIdResource() resource.Resource {
	return &BundleBundleIdResource{
		TypeName: "distribution_bundle_id",
	}
}

type BundleBundleIdResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type BundleBundleIdResourceModel struct {
	Parameters types.Map `tfsdk:"parameters"`
	Configuration types.Bool `tfsdk:"configuration"`
	Logs types.Map `tfsdk:"logs"`
	Include types.Bool `tfsdk:"include"`
	StartDate types.String `tfsdk:"start_date"`
	EndDate types.String `tfsdk:"end_date"`
	System types.Bool `tfsdk:"system"`
	ThreadDump types.Map `tfsdk:"thread_dump"`
	Count types.Int64 `tfsdk:"count"`
	Interval types.Int64 `tfsdk:"interval"`
	Description types.String `tfsdk:"description"`
	Name types.String `tfsdk:"name"`
	Id types.Int64 `tfsdk:"id"`
	Artifactory types.Map `tfsdk:"artifactory"`
	ServiceId types.String `tfsdk:"service_id"`
	Parameters.Configuration types.Bool `tfsdk:"parameters.configuration"`
	Parameters.Logs types.Map `tfsdk:"parameters.logs"`
	Parameters.Logs.Include types.Bool `tfsdk:"parameters.logs.include"`
	Parameters.Logs.StartDate types.String `tfsdk:"parameters.logs.start_date"`
	Parameters.Logs.EndDate types.String `tfsdk:"parameters.logs.end_date"`
	ThreadDump.Count types.Int64 `tfsdk:"thread_dump.count"`
	ThreadDump.Interval types.Int64 `tfsdk:"thread_dump.interval"`
}

type BundleIdRequestAPIModel struct {
	Parameters string `json:"parameters"`
	Configuration string `json:"configuration"`
	Logs string `json:"logs"`
	Include string `json:"include"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	System string `json:"system"`
	ThreadDump string `json:"thread_dump"`
	Count string `json:"count"`
	Interval string `json:"interval"`
	Description string `json:"description"`
	Name string `json:"name"`
	Id string `json:"id"`
	Artifactory string `json:"artifactory"`
	ServiceId string `json:"service_id"`
	Parameters.Configuration string `json:"parameters.configuration"`
	Parameters.Logs string `json:"parameters.logs"`
	Parameters.Logs.Include string `json:"parameters.logs.include"`
	Parameters.Logs.StartDate string `json:"parameters.logs.start_date"`
	Parameters.Logs.EndDate string `json:"parameters.logs.end_date"`
	ThreadDump.Count string `json:"thread_dump.count"`
	ThreadDump.Interval string `json:"thread_dump.interval"`
}

type BundleIdAPIModel struct {
	Parameters string `json:"parameters"`
	Configuration string `json:"configuration"`
	Logs string `json:"logs"`
	Include string `json:"include"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	System string `json:"system"`
	ThreadDump string `json:"thread_dump"`
	Count string `json:"count"`
	Interval string `json:"interval"`
	Description string `json:"description"`
	Name string `json:"name"`
	Id string `json:"id"`
	Artifactory string `json:"artifactory"`
	ServiceId string `json:"service_id"`
	Parameters.Configuration string `json:"parameters.configuration"`
	Parameters.Logs string `json:"parameters.logs"`
	Parameters.Logs.Include string `json:"parameters.logs.include"`
	Parameters.Logs.StartDate string `json:"parameters.logs.start_date"`
	Parameters.Logs.EndDate string `json:"parameters.logs.end_date"`
	ThreadDump.Count string `json:"thread_dump.count"`
	ThreadDump.Interval string `json:"thread_dump.interval"`
}

func (r *BundleBundleIdResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *BundleBundleIdResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"parameters": schema.StringAttribute{
				Optional: true,
				Description: "Support bundle parameters",
			},
			"configuration": schema.StringAttribute{
				Optional: true,
				Description: "The configuration of the resource.",
			},
			"logs": schema.StringAttribute{
				Optional: true,
				Description: "The logs of the resource.",
			},
			"include": schema.StringAttribute{
				Optional: true,
				Description: "The include of the resource.",
			},
			"start_date": schema.StringAttribute{
				Optional: true,
				Description: "The start_date of the resource.",
			},
			"end_date": schema.StringAttribute{
				Optional: true,
				Description: "The end_date of the resource.",
			},
			"system": schema.StringAttribute{
				Optional: true,
				Description: "Information about your system including storage, system properties, cpu, and JVM information",
			},
			"thread_dump": schema.StringAttribute{
				Optional: true,
				Description: "Create a thread dump for all running threads",
			},
			"count": schema.StringAttribute{
				Optional: true,
				Description: "The count of the resource.",
			},
			"interval": schema.StringAttribute{
				Optional: true,
				Description: "The interval of the resource.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Description: "Support bundle description",
			},
			"name": schema.StringAttribute{
				Optional: true,
				Description: "The name of the resource.",
			},
			"id": schema.StringAttribute{
				Optional: true,
				Description: "The id of the resource.",
			},
			"artifactory": schema.StringAttribute{
				Optional: true,
				Description: "The artifactory of the resource.",
			},
			"service_id": schema.StringAttribute{
				Optional: true,
				Description: "The service_id of the resource.",
			},
			"parameters.configuration": schema.StringAttribute{
				Optional: true,
				Description: "Collect configuration files",
			},
			"parameters.logs": schema.StringAttribute{
				Optional: true,
				Description: "Collect all system logs, if this field is not specified support bundle will collect logs from day before today until today.",
			},
			"parameters.logs.include": schema.StringAttribute{
				Optional: true,
				Description: "Collect system logs.",
			},
			"parameters.logs.start_date": schema.StringAttribute{
				Optional: true,
				Description: "Start date from which to fetch the logs",
			},
			"parameters.logs.end_date": schema.StringAttribute{
				Optional: true,
				Description: "End date until which to fetch the logs",
			},
			"thread_dump.count": schema.StringAttribute{
				Optional: true,
				Description: "Number of thread dumps to collect",
			},
			"thread_dump.interval": schema.StringAttribute{
				Optional: true,
				Description: "Interval between times of collect thread dump in milliseconds",
			},
		},
		MarkdownDescription: "Manages bundle_id in JFrog Distribution.",
	}
}

func (r *BundleBundleIdResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}



func (r *BundleBundleIdResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state BundleBundleIdResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result BundleBundleIdAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
		}).
		SetResult(&result).
		Get(BundleIdsEndpoint)
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


func (r *BundleBundleIdResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	go util.SendUsageResourceUpdate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan BundleBundleIdResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := BundleBundleIdRequestAPIModel{
		Parameters: plan.Parameters.ValueString(),
		Configuration: plan.Configuration.ValueString(),
		Logs: plan.Logs.ValueString(),
		Include: plan.Include.ValueString(),
		StartDate: plan.StartDate.ValueString(),
		EndDate: plan.EndDate.ValueString(),
		System: plan.System.ValueString(),
		ThreadDump: plan.ThreadDump.ValueString(),
		Count: plan.Count.ValueString(),
		Interval: plan.Interval.ValueString(),
		Description: plan.Description.ValueString(),
		Name: plan.Name.ValueString(),
		Id: plan.Id.ValueString(),
		Artifactory: plan.Artifactory.ValueString(),
		ServiceId: plan.ServiceId.ValueString(),
		Parameters.Configuration: plan.Parameters.Configuration.ValueString(),
		Parameters.Logs: plan.Parameters.Logs.ValueString(),
		Parameters.Logs.Include: plan.Parameters.Logs.Include.ValueString(),
		Parameters.Logs.StartDate: plan.Parameters.Logs.StartDate.ValueString(),
		Parameters.Logs.EndDate: plan.Parameters.Logs.EndDate.ValueString(),
		ThreadDump.Count: plan.ThreadDump.Count.ValueString(),
		ThreadDump.Interval: plan.ThreadDump.Interval.ValueString(),
	}

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
		}).
		SetBody(requestBody).
		Put(BundleIdsEndpoint)
	if err != nil {
		utilfw.UnableToUpdateResourceError(resp, err.Error())
		return
	}

	if response.IsError() {
		utilfw.UnableToUpdateResourceError(resp, response.String())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}



func (r *BundleBundleIdResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	go util.SendUsageResourceDelete(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state BundleBundleIdResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
		}).
		Delete(BundleIdsEndpoint)
	if err != nil {
		utilfw.UnableToDeleteResourceError(resp, err.Error())
		return
	}

	if response.StatusCode() == http.StatusNotFound {
		return
	}

	if response.IsError() {
		utilfw.UnableToDeleteResourceError(resp, response.String())
		return
	}
}

