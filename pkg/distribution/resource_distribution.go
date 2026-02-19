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
	DistributionsEndpoint = "distribution/api/v1"
	DistributionEndpoint  = "distribution/api/v1/distribution/{name}/{version}"
)

func NewDistributionResource() resource.Resource {
	return &DistributionResource{
		TypeName: "distribution_distribution",
	}
}

type DistributionResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type DistributionResourceModel struct {
	Name types.String `tfsdk:"name"`
	Version types.String `tfsdk:"version"`
	DistributionDestinations types.String `tfsdk:"distribution_destinations"`
	SiteName types.String `tfsdk:"site_name"`
	CityName types.String `tfsdk:"city_name"`
	CountryCodes types.String `tfsdk:"country_codes"`
	ResourceType types.String `tfsdk:"resource_type"`
	Principals types.Map `tfsdk:"principals"`
	Users types.Map `tfsdk:"users"`
	Anonymous types.String `tfsdk:"anonymous"`
	Groups types.Map `tfsdk:"groups"`
	Readers types.String `tfsdk:"readers"`
	User1 types.String `tfsdk:"user1"`
	Group1 types.String `tfsdk:"group1"`
}

type DistributionRequestAPIModel struct {
	Name string `json:"name"`
	Version string `json:"version"`
	DistributionDestinations string `json:"distribution_destinations"`
	SiteName string `json:"site_name"`
	CityName string `json:"city_name"`
	CountryCodes string `json:"country_codes"`
	ResourceType string `json:"resource_type"`
	Principals string `json:"principals"`
	Users string `json:"users"`
	Anonymous string `json:"anonymous"`
	Groups string `json:"groups"`
	Readers string `json:"readers"`
	User1 string `json:"user1"`
	Group1 string `json:"group1"`
}

type DistributionAPIModel struct {
	Name string `json:"name"`
	Version string `json:"version"`
	DistributionDestinations string `json:"distribution_destinations"`
	SiteName string `json:"site_name"`
	CityName string `json:"city_name"`
	CountryCodes string `json:"country_codes"`
	ResourceType string `json:"resource_type"`
	Principals string `json:"principals"`
	Users string `json:"users"`
	Anonymous string `json:"anonymous"`
	Groups string `json:"groups"`
	Readers string `json:"readers"`
	User1 string `json:"user1"`
	Group1 string `json:"group1"`
}

func (r *DistributionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *DistributionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"distribution_destinations": schema.StringAttribute{
				Optional: true,
				Description: "The distribution_destinations of the resource.",
			},
			"site_name": schema.StringAttribute{
				Optional: true,
				Description: "The site_name of the resource.",
			},
			"city_name": schema.StringAttribute{
				Optional: true,
				Description: "The city_name of the resource.",
			},
			"country_codes": schema.StringAttribute{
				Optional: true,
				Description: "The country_codes of the resource.",
			},
			"resource_type": schema.StringAttribute{
				Optional: true,
				Description: "The resource_type of the resource.",
			},
			"principals": schema.StringAttribute{
				Optional: true,
				Description: "The principals of the resource.",
			},
			"users": schema.StringAttribute{
				Optional: true,
				Description: "The users of the resource.",
			},
			"anonymous": schema.StringAttribute{
				Optional: true,
				Description: "The anonymous of the resource.",
			},
			"groups": schema.StringAttribute{
				Optional: true,
				Description: "The groups of the resource.",
			},
			"readers": schema.StringAttribute{
				Optional: true,
				Description: "The readers of the resource.",
			},
			"user1": schema.StringAttribute{
				Optional: true,
				Description: "The user1 of the resource.",
			},
			"group1": schema.StringAttribute{
				Optional: true,
				Description: "The group1 of the resource.",
			},
		},
		MarkdownDescription: "Manages distribution in JFrog Distribution.",
	}
}

func (r *DistributionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}


func (r *DistributionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	go util.SendUsageResourceCreate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan DistributionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := DistributionRequestAPIModel{
		Name: plan.Name.ValueString(),
		Version: plan.Version.ValueString(),
		DistributionDestinations: plan.DistributionDestinations.ValueString(),
		SiteName: plan.SiteName.ValueString(),
		CityName: plan.CityName.ValueString(),
		CountryCodes: plan.CountryCodes.ValueString(),
		ResourceType: plan.ResourceType.ValueString(),
		Principals: plan.Principals.ValueString(),
		Users: plan.Users.ValueString(),
		Anonymous: plan.Anonymous.ValueString(),
		Groups: plan.Groups.ValueString(),
		Readers: plan.Readers.ValueString(),
		User1: plan.User1.ValueString(),
		Group1: plan.Group1.ValueString(),
	}

	var result DistributionAPIModel

	response, err := r.ProviderData.Client.R().
		SetBody(requestBody).
		SetResult(&result).
		Post(DistributionsEndpoint)
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


func (r *DistributionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state DistributionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result DistributionAPIModel

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


func (r *DistributionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	go util.SendUsageResourceUpdate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan DistributionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody := DistributionRequestAPIModel{
		Name: plan.Name.ValueString(),
		Version: plan.Version.ValueString(),
		DistributionDestinations: plan.DistributionDestinations.ValueString(),
		SiteName: plan.SiteName.ValueString(),
		CityName: plan.CityName.ValueString(),
		CountryCodes: plan.CountryCodes.ValueString(),
		ResourceType: plan.ResourceType.ValueString(),
		Principals: plan.Principals.ValueString(),
		Users: plan.Users.ValueString(),
		Anonymous: plan.Anonymous.ValueString(),
		Groups: plan.Groups.ValueString(),
		Readers: plan.Readers.ValueString(),
		User1: plan.User1.ValueString(),
		Group1: plan.Group1.ValueString(),
	}

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"name": plan.Name.ValueString(),
			"version": plan.Version.ValueString(),
		}).
		SetBody(requestBody).
		Put(DistributionEndpoint)
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



func (r *DistributionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	go util.SendUsageResourceDelete(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state DistributionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"name": state.Name.ValueString(),
			"version": state.Version.ValueString(),
		}).
		Delete(DistributionEndpoint)
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

