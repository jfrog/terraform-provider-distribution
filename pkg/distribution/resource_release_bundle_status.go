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
	StatusEndpoint = "distribution/api/v1/export/release_bundle/{name}/{version}/status"
	StatusEndpoint  = "distribution/api/v1/export/release_bundle/{name}/{version}/status"
)

func NewReleaseBundleStatusResource() resource.Resource {
	return &ReleaseBundleStatusResource{
		TypeName: "distribution_status",
	}
}

type ReleaseBundleStatusResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type ReleaseBundleStatusResourceModel struct {
	Name types.String `tfsdk:"name"`
	Version types.String `tfsdk:"version"`
	Status types.String `tfsdk:"status"`
	Message types.String `tfsdk:"message"`
	DownloadUrl types.String `tfsdk:"download_url"`
}

type StatusAPIModel struct {
	Name string `json:"name"`
	Version string `json:"version"`
	Status string `json:"status"`
	Message string `json:"message"`
	DownloadUrl string `json:"download_url"`
}

func (r *ReleaseBundleStatusResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *ReleaseBundleStatusResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"status": schema.StringAttribute{
				Optional: true,
				Description: "The status of the resource.",
			},
			"message": schema.StringAttribute{
				Optional: true,
				Description: "The message of the resource.",
			},
			"download_url": schema.StringAttribute{
				Optional: true,
				Description: "The download_url of the resource.",
			},
		},
		MarkdownDescription: "Manages status in JFrog Distribution.",
	}
}

func (r *ReleaseBundleStatusResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}



func (r *ReleaseBundleStatusResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state ReleaseBundleStatusResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result ReleaseBundleStatusAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"name": state.Name.ValueString(),
			"version": state.Version.ValueString(),
		}).
		SetResult(&result).
		Get(StatusEndpoint)
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




