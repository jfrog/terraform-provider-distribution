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
	AbortEndpoint = "distribution/api/v1/distribution/{name}/{version}/abort"
	AbortEndpoint  = "distribution/api/v1/distribution/{name}/{version}/abort"
)

func NewDistributionAbortResource() resource.Resource {
	return &DistributionAbortResource{
		TypeName: "distribution_abort",
	}
}

type DistributionAbortResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type DistributionAbortResourceModel struct {
	Name types.String `tfsdk:"name"`
	Version types.String `tfsdk:"version"`
}

type AbortRequestAPIModel struct {
	Name string `json:"name"`
	Version string `json:"version"`
}

type AbortAPIModel struct {
	Name string `json:"name"`
	Version string `json:"version"`
}

func (r *DistributionAbortResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *DistributionAbortResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				Description: "The name of the resource.",
			},
			"version": schema.StringAttribute{
				Required: true,
				Description: "The version of the resource.",
			},
		},
		MarkdownDescription: "Manages abort in JFrog Distribution.",
	}
}

func (r *DistributionAbortResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}



func (r *DistributionAbortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state DistributionAbortResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result DistributionAbortAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"name": state.Name.ValueString(),
			"version": state.Version.ValueString(),
		}).
		SetResult(&result).
		Get(AbortEndpoint)
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


func (r *DistributionAbortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	go util.SendUsageResourceUpdate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan DistributionAbortResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := DistributionAbortRequestAPIModel{
		Name: plan.Name.ValueString(),
		Version: plan.Version.ValueString(),
	}

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"name": plan.Name.ValueString(),
			"version": plan.Version.ValueString(),
		}).
		SetBody(requestBody).
		Put(AbortEndpoint)
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



