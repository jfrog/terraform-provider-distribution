package distribution

import (
	"context"
	"net/http"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jfrog/terraform-provider-shared/util"
	utilfw "github.com/jfrog/terraform-provider-shared/util/fw"
	"github.com/samber/lo"
)

func NewPermissionResource() resource.Resource {
	return &PermissionResource{
		TypeName: "distribution_permission_target",
	}
}

type PermissionResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

func (r *PermissionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *PermissionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes:          permissionSchemaAttributes,
		MarkdownDescription: "This resource enables you to manage permissions in JFrog Distribution. For more information, see [Permission Management](https://jfrog.com/help/r/jfrog-platform-administration-documentation/permission-management) and [REST API](https://jfrog.com/help/r/jfrog-rest-apis/permission-management).",
	}
}

func (r *PermissionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}

func (r *PermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	go util.SendUsageResourceCreate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan PermissionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiModel, diags := r.toAPIModel(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result PermissionAPIModel
	var putErr PermissionErrorAPIModel

	response, err := r.ProviderData.Client.R().
		SetBody(apiModel).
		SetResult(&result).
		SetError(&putErr).
		Put(PermissionsEndpoint + apiModel.Name)
	if err != nil {
		utilfw.UnableToCreateResourceError(resp, err.Error())
		return
	}

	if response.IsError() {
		utilfw.UnableToCreateResourceError(resp, putErr.String())
		return
	}

	tflog.Info(ctx, "Permission Create", map[string]interface{}{
		"name": result.Name,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state PermissionResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var permission PermissionAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParam("permissionName", state.Name.ValueString()).
		SetResult(&permission).
		Get(PermissionEndpoint)
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

	// Update state with API response
	state = r.fromAPIModel(ctx, permission, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	go util.SendUsageResourceUpdate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan PermissionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiModel, diags := r.toAPIModel(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.ProviderData.Client.R().
		SetPathParam("permissionName", plan.Name.ValueString()).
		SetBody(apiModel).
		Put(PermissionEndpoint)
	if err != nil {
		utilfw.UnableToUpdateResourceError(resp, err.Error())
		return
	}

	if response.IsError() {
		utilfw.UnableToUpdateResourceError(resp, response.String())
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	go util.SendUsageResourceDelete(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state PermissionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	response, err := r.ProviderData.Client.R().
		SetPathParam("permissionName", state.Name.ValueString()).
		Delete(PermissionEndpoint)

	if err != nil {
		utilfw.UnableToDeleteResourceError(resp, err.Error())
		return
	}

	if response.IsError() {
		utilfw.UnableToDeleteResourceError(resp, response.String())
		return
	}

	// If the logic reaches here, it implicitly succeeded and will remove
	// the resource from state if there are no other errors.
}

func (r *PermissionResource) toAPIModel(ctx context.Context, model PermissionResourceModel) (PermissionAPIModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	apiModel := PermissionAPIModel{
		Name: model.Name.ValueString(),
	}

	if !model.ResourceType.IsNull() && !model.ResourceType.IsUnknown() {
		apiModel.ResourceType = model.ResourceType.ValueString()
	}

	// Handle distribution_destinations
	if !model.DistributionDestinations.IsNull() && !model.DistributionDestinations.IsUnknown() {
		var destinations []DistributionDestination
		diags.Append(model.DistributionDestinations.ElementsAs(ctx, &destinations, false)...)
		apiModel.DistributionDestinations = destinations
	}

	// Handle principals
	if !model.Principals.IsNull() && !model.Principals.IsUnknown() {
		var principals PermissionPrincipals
		diags.Append(model.Principals.As(ctx, &principals, basetypes.ObjectAsOptions{})...)

		// Sort permissions for consistency
		sortedPrincipals := PermissionPrincipals{
			Users:  make(map[string][]string),
			Groups: make(map[string][]string),
		}

		for user, perms := range principals.Users {
			sortedPerms := make([]string, len(perms))
			copy(sortedPerms, perms)
			sort.Strings(sortedPerms)
			sortedPrincipals.Users[user] = sortedPerms
		}

		for group, perms := range principals.Groups {
			sortedPerms := make([]string, len(perms))
			copy(sortedPerms, perms)
			sort.Strings(sortedPerms)
			sortedPrincipals.Groups[group] = sortedPerms
		}

		apiModel.Principals = sortedPrincipals
	}

	return apiModel, diags
}

func (r *PermissionResource) fromAPIModel(ctx context.Context, apiModel PermissionAPIModel, diags *diag.Diagnostics) PermissionResourceModel {
	model := PermissionResourceModel{
		Name: types.StringValue(apiModel.Name),
	}

	if apiModel.ResourceType != "" {
		model.ResourceType = types.StringValue(apiModel.ResourceType)
	} else {
		model.ResourceType = types.StringNull()
	}

	// Handle distribution_destinations
	if len(apiModel.DistributionDestinations) > 0 {
		destinations := lo.Map(apiModel.DistributionDestinations, func(dest DistributionDestination, _ int) attr.Value {
			val, _ := types.ObjectValue(map[string]attr.Type{
				"site_name":     types.StringType,
				"city_name":     types.StringType,
				"country_codes": types.ListType{ElemType: types.StringType},
			}, map[string]attr.Value{
				"site_name": types.StringValue(dest.SiteName),
				"city_name": types.StringValue(dest.CityName),
				"country_codes": types.ListValueMust(types.StringType, lo.Map(dest.CountryCodes, func(code string, _ int) attr.Value {
					return types.StringValue(code)
				})),
			})
			return val
		})
		model.DistributionDestinations = types.ListValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"site_name":     types.StringType,
				"city_name":     types.StringType,
				"country_codes": types.ListType{ElemType: types.StringType},
			},
		}, destinations)
	} else {
		model.DistributionDestinations = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"site_name":     types.StringType,
				"city_name":     types.StringType,
				"country_codes": types.ListType{ElemType: types.StringType},
			},
		})
	}

	// Handle principals - always create the object structure with empty maps
	principals := PermissionPrincipals{
		Users:  make(map[string][]string),
		Groups: make(map[string][]string),
	}

	// Populate with API data if available and sort for consistency
	if apiModel.Principals.Users != nil {
		principals.Users = make(map[string][]string)
		for user, perms := range apiModel.Principals.Users {
			sortedPerms := make([]string, len(perms))
			copy(sortedPerms, perms)
			sort.Strings(sortedPerms)
			principals.Users[user] = sortedPerms
		}
	}
	if apiModel.Principals.Groups != nil {
		principals.Groups = make(map[string][]string)
		for group, perms := range apiModel.Principals.Groups {
			sortedPerms := make([]string, len(perms))
			copy(sortedPerms, perms)
			sort.Strings(sortedPerms)
			principals.Groups[group] = sortedPerms
		}
	}

	principalsValue, diagVal := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"users":  types.MapType{ElemType: types.ListType{ElemType: types.StringType}},
		"groups": types.MapType{ElemType: types.ListType{ElemType: types.StringType}},
	}, principals)
	diags.Append(diagVal...)
	model.Principals = principalsValue

	return model
}

// ImportState imports the resource into the Terraform state.
func (r *PermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// For permission targets, the import ID is just the permission name
	permissionName := req.ID

	if permissionName == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			"Expected permission name",
		)
		return
	}

	// Set the name attribute in the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), permissionName)...)
}
