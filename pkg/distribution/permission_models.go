package distribution

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// API Endpoints
const (
	PermissionsEndpoint = "distribution/api/v1/security/permissions/"
	PermissionEndpoint  = "distribution/api/v1/security/permissions/{permissionName}"
)

// DistributionDestination matches the JSON structure for distribution_destinations
type DistributionDestination struct {
	SiteName     string   `tfsdk:"site_name" json:"site_name"`
	CityName     string   `tfsdk:"city_name" json:"city_name"`
	CountryCodes []string `tfsdk:"country_codes" json:"country_codes"`
}

// Permission Models
type PermissionResourceModel struct {
	Name                     types.String `tfsdk:"name"`
	ResourceType             types.String `tfsdk:"resource_type"`
	DistributionDestinations types.List   `tfsdk:"distribution_destinations"`
	Principals               types.Object `tfsdk:"principals"`
}

type PermissionAPIModel struct {
	Name                     string                    `json:"name"`
	ResourceType             string                    `json:"resource_type"`
	DistributionDestinations []DistributionDestination `json:"distribution_destinations"`
	Principals               PermissionPrincipals      `json:"principals"`
}

type PermissionPrincipals struct {
	Users  map[string][]string `tfsdk:"users" json:"users,omitempty"`
	Groups map[string][]string `tfsdk:"groups" json:"groups,omitempty"`
}

// Permission Schema Attributes
var permissionSchemaAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
		Description: "Name of the permission",
	},
	"resource_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			validateResourceType(), // Only "destination" is allowed
		},
		Description: "Resource type for the permission (only 'destination' is allowed)",
	},
	"distribution_destinations": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"site_name": schema.StringAttribute{
					Required:    true,
					Description: "Site name for the distribution destination",
				},
				"city_name": schema.StringAttribute{
					Required:    true,
					Description: "City name for the distribution destination",
				},
				"country_codes": schema.ListAttribute{
					ElementType: types.StringType,
					Required:    true,
					Description: "Country codes for the distribution destination",
				},
			},
			Validators: []validator.Object{
				validateDistributionDestination(), // Ensure all required fields are provided
			},
		},
		Required: true, // Change from Optional to Required
		Validators: []validator.List{
			validateDistributionDestinations(), // At least one destination required
		},
		Description: "Distribution destinations for the permission (at least one required)",
	},
	"principals": schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"users": schema.MapAttribute{
				ElementType: types.ListType{ElemType: types.StringType},
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Map{
					planModifierForEmptyMap{},
				},
				Description: "User principals for the permission",
			},
			"groups": schema.MapAttribute{
				ElementType: types.ListType{ElemType: types.StringType},
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Map{
					planModifierForEmptyMap{},
				},
				Description: "Group principals for the permission",
			},
		},
		Required: true, // Change from Optional to Required
		Validators: []validator.Object{
			validatePrincipals(), // At least one user or group required
		},
		Description: "Principals for the permission (at least one user or group required)",
	},
}

// Error Models
type PermissionErrorAPIModel struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Detail     string `json:"detail"`
}

func (m PermissionErrorAPIModel) String() string {
	return fmt.Sprintf("%d - %s: %s", m.StatusCode, m.Message, m.Detail)
}

// Custom validators

// validateResourceType ensures only "destination" is allowed
func validateResourceType() validator.String {
	return stringvalidator.OneOf("destination")
}

// validatePrincipals ensures at least one user or group is provided
func validatePrincipals() validator.Object {
	return objectvalidator.AtLeastOneOf(
		path.MatchRelative().AtName("users"),
		path.MatchRelative().AtName("groups"),
	)
}

// validateDistributionDestinations ensures at least one destination is provided
func validateDistributionDestinations() validator.List {
	return listvalidator.SizeAtLeast(1)
}

// validateDistributionDestination ensures site_name, city_name, and country_codes are provided
func validateDistributionDestination() validator.Object {
	return objectvalidator.AlsoRequires(
		path.MatchRelative().AtName("site_name"),
		path.MatchRelative().AtName("city_name"),
		path.MatchRelative().AtName("country_codes"),
	)
}

// planModifierForEmptyMap ensures empty maps are set when not provided
type planModifierForEmptyMap struct{}

func (m planModifierForEmptyMap) Description(ctx context.Context) string {
	return "Sets empty map when not provided"
}

func (m planModifierForEmptyMap) MarkdownDescription(ctx context.Context) string {
	return "Sets empty map when not provided"
}

func (m planModifierForEmptyMap) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	// If the value is null or unknown, set it to an empty map
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		resp.PlanValue = types.MapValueMust(
			types.ListType{ElemType: types.StringType},
			map[string]attr.Value{},
		)
	}
}
