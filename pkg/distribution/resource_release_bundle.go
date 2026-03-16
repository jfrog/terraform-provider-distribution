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
	ReleaseBundlesEndpoint = "distribution/api/v1/release_bundle[?start_pos=:position]"
	ReleaseBundlesEndpoint  = "distribution/api/v1/release_bundle[?start_pos=:position]"
)

func NewReleaseBundleResource() resource.Resource {
	return &ReleaseBundleResource{
		TypeName: "distribution_release_bundle",
	}
}

type ReleaseBundleResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type ReleaseBundleResourceModel struct {
	Name types.String `tfsdk:"name"`
	Version types.String `tfsdk:"version"`
	State types.String `tfsdk:"state"`
	Description types.String `tfsdk:"description"`
	ReleaseNotes types.Map `tfsdk:"release_notes"`
	Content types.String `tfsdk:"content"`
	Syntax types.String `tfsdk:"syntax"`
	Created types.String `tfsdk:"created"`
	CreatedBy types.String `tfsdk:"created_by"`
	DistributedBy types.String `tfsdk:"distributed_by"`
	Artifacts types.String `tfsdk:"artifacts"`
	Archived types.Bool `tfsdk:"archived"`
}

type ReleaseBundleAPIModel struct {
	Name string `json:"name"`
	Version string `json:"version"`
	State string `json:"state"`
	Description string `json:"description"`
	ReleaseNotes string `json:"release_notes"`
	Content string `json:"content"`
	Syntax string `json:"syntax"`
	Created string `json:"created"`
	CreatedBy string `json:"created_by"`
	DistributedBy string `json:"distributed_by"`
	Artifacts string `json:"artifacts"`
	Archived string `json:"archived"`
}

func (r *ReleaseBundleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *ReleaseBundleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
				Description: "The name of the resource.",
			},
			"version": schema.StringAttribute{
				Optional: true,
				Description: "The version of the resource.",
			},
			"state": schema.StringAttribute{
				Optional: true,
				Description: "The state of the resource.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Description: "The description of the resource.",
			},
			"release_notes": schema.StringAttribute{
				Optional: true,
				Description: "The release_notes of the resource.",
			},
			"content": schema.StringAttribute{
				Optional: true,
				Description: "The content of the resource.",
			},
			"syntax": schema.StringAttribute{
				Optional: true,
				Description: "The syntax of the resource.",
			},
			"created": schema.StringAttribute{
				Optional: true,
				Description: "The created of the resource.",
			},
			"created_by": schema.StringAttribute{
				Optional: true,
				Description: "The created_by of the resource.",
			},
			"distributed_by": schema.StringAttribute{
				Optional: true,
				Description: "The distributed_by of the resource.",
			},
			"artifacts": schema.StringAttribute{
				Optional: true,
				Description: "The artifacts of the resource.",
			},
			"archived": schema.StringAttribute{
				Optional: true,
				Description: "The archived of the resource.",
			},
		},
		MarkdownDescription: "Manages release_bundle in JFrog Distribution.",
	}
}

func (r *ReleaseBundleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}



func (r *ReleaseBundleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state ReleaseBundleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result ReleaseBundleAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
		}).
		SetResult(&result).
		Get(ReleaseBundlesEndpoint)
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




