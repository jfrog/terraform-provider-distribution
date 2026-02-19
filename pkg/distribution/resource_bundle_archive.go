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
	ArchiveEndpoint = "distribution/api/v1/system/support/bundle/{bundle_id}/archive"
	ArchiveEndpoint  = "distribution/api/v1/system/support/bundle/{bundle_id}/archive"
)

func NewBundleArchiveResource() resource.Resource {
	return &BundleArchiveResource{
		TypeName: "distribution_archive",
	}
}

type BundleArchiveResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type BundleArchiveResourceModel struct {
	BundleId types.String `tfsdk:"bundle_id"`
}

type ArchiveAPIModel struct {
	BundleId string `json:"bundle_id"`
}

func (r *BundleArchiveResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *BundleArchiveResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bundle_id": schema.StringAttribute{
				Required: true,
				Description: "The bundle_id of the resource.",
			},
		},
		MarkdownDescription: "Manages archive in JFrog Distribution.",
	}
}

func (r *BundleArchiveResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}



func (r *BundleArchiveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state BundleArchiveResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result BundleArchiveAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"bundle_id": state.BundleId.ValueString(),
		}).
		SetResult(&result).
		Get(ArchiveEndpoint)
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




