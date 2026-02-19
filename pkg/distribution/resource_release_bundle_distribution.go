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
	DistributionEndpoint = "distribution/api/v1/release_bundle/{name}/{version}/distribution"
	DistributionEndpoint  = "distribution/api/v1/release_bundle/{name}/{version}/distribution"
)

func NewReleaseBundleDistributionResource() resource.Resource {
	return &ReleaseBundleDistributionResource{
		TypeName: "distribution_distribution",
	}
}

type ReleaseBundleDistributionResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type ReleaseBundleDistributionResourceModel struct {
	Name types.String `tfsdk:"name"`
	Version types.String `tfsdk:"version"`
}

type DistributionAPIModel struct {
	Name string `json:"name"`
	Version string `json:"version"`
}

func (r *ReleaseBundleDistributionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *ReleaseBundleDistributionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
		MarkdownDescription: "Manages distribution in JFrog Distribution.",
	}
}

func (r *ReleaseBundleDistributionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}



func (r *ReleaseBundleDistributionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state ReleaseBundleDistributionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result ReleaseBundleDistributionAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"name": state.Name.ValueString(),
			"version": state.Version.ValueString(),
		}).
		SetResult(&result).
		Get(DistributionEndpoint)
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




