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
	ReleaseBundleEndpoint = "distribution/api/v1/export/release_bundle/{name}/{version}"
	ReleaseBundleEndpoint  = "distribution/api/v1/export/release_bundle/{name}/{version}"
)

func NewExportReleaseBundleResource() resource.Resource {
	return &ExportReleaseBundleResource{
		TypeName: "distribution_release_bundle",
	}
}

type ExportReleaseBundleResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type ExportReleaseBundleResourceModel struct {
	Name types.String `tfsdk:"name"`
	Version types.String `tfsdk:"version"`
}

type ReleaseBundleRequestAPIModel struct {
	Name string `json:"name"`
	Version string `json:"version"`
}

type ReleaseBundleAPIModel struct {
	Name string `json:"name"`
	Version string `json:"version"`
}

func (r *ExportReleaseBundleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *ExportReleaseBundleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
		MarkdownDescription: "Manages release_bundle in JFrog Distribution.",
	}
}

func (r *ExportReleaseBundleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}


func (r *ExportReleaseBundleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	go util.SendUsageResourceCreate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan ExportReleaseBundleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := ExportReleaseBundleRequestAPIModel{
		Name: plan.Name.ValueString(),
		Version: plan.Version.ValueString(),
	}

	var result ExportReleaseBundleAPIModel

	response, err := r.ProviderData.Client.R().
		SetBody(requestBody).
		SetResult(&result).
		Post(ReleaseBundleEndpoint)
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


func (r *ExportReleaseBundleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state ExportReleaseBundleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result ExportReleaseBundleAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"name": state.Name.ValueString(),
			"version": state.Version.ValueString(),
		}).
		SetResult(&result).
		Get(ReleaseBundleEndpoint)
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




func (r *ExportReleaseBundleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	go util.SendUsageResourceDelete(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state ExportReleaseBundleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"name": state.Name.ValueString(),
			"version": state.Version.ValueString(),
		}).
		Delete(ReleaseBundleEndpoint)
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

