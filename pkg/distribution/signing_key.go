package distribution

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
)

const (
	SigningKeysEndpoint = "distribution/api/v1/keys/{protocol}"
	SigningKeyEndpoint  = "distribution/api/v1/keys/{protocol}/{alias}"
)

type SigningKeyCommmonResourceModel struct {
	Protocol                 types.String `tfsdk:"protocol"`
	Alias                    types.String `tfsdk:"alias"`
	PropagateToEdgeNode      types.Bool   `tfsdk:"propagate_to_edge_nodes"`
	FailOnPropagationFailure types.Bool   `tfsdk:"fail_on_propagation_failure"`
	SetAsDefault             types.Bool   `tfsdk:"set_as_default"`
}

type SigningKeyCommonPostRequestAPIModel struct {
	PropagateToEdgeNode      bool `json:"propagate_to_edge_nodes"`
	FailOnPropagationFailure bool `json:"fail_on_propagation_failure"`
	SetAsDefault             bool `json:"set_as_default"`
}

type SigningKeyPostResponseAPIModel struct {
	Report SigningKeyReportPostResponseAPIModel `json:"report"`
}

type SigningKeyReportPostResponseAPIModel struct {
	Message string                                       `json:"message"`
	Status  string                                       `json:"status"`
	Details []SigningKeyReportDetailPostResponseAPIModel `json:"details"`
}

type SigningKeyReportDetailPostResponseAPIModel struct {
	JPDID    string `json:"jpd_id"`
	Name     string `json:"name"`
	KeyAlias string `json:"key_alias"`
	Status   string `json:"status"`
}

func (m SigningKeyReportPostResponseAPIModel) String() string {
	details := lo.Map(
		m.Details,
		func(detail SigningKeyReportDetailPostResponseAPIModel, _ int) string {
			return fmt.Sprintf("JPD ID: %s, Name: %s, Key alias: %s, Status: %s", detail.JPDID, detail.Name, detail.KeyAlias, detail.Status)
		},
	)
	return fmt.Sprintf("%s: %s - %s", m.Status, m.Message, strings.Join(details, ",\n"))
}

type SigningKeyPostErrorAPIModel struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Detail     string `json:"detail"`
}

func (m SigningKeyPostErrorAPIModel) String() string {
	return fmt.Sprintf("%d - %s: %s", m.StatusCode, m.Message, m.Detail)
}

type SigningKeyGetAPIModel struct {
	Alias     string `json:"alias"`
	PublicKey string `json:"public_key"`
}

type SigningKeyPutRequestAPIModel struct {
	NewAlias string `json:"new_alias"`
}

var commonSchemaAttributes = map[string]schema.Attribute{
	"alias": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
		Description: "Alias of the signing key",
	},
	"protocol": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("gpg", "pgp"),
		},
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
		Description: "Type of the signing key. Valid value: `gpg` or `pgp`",
	},
	"propagate_to_edge_nodes": schema.BoolAttribute{
		Optional: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.RequiresReplace(),
		},
		MarkdownDescription: "When set to `true`, the public key will be automatically propagated to the Edge Node just once.",
	},
	"fail_on_propagation_failure": schema.BoolAttribute{
		Optional: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.RequiresReplace(),
		},
		MarkdownDescription: "When set to `true`, the public key will be automatically propagated to the Edge Node just once.",
	},
	"set_as_default": schema.BoolAttribute{
		Optional: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.RequiresReplace(),
		},
		MarkdownDescription: "Set this to `true` if this is the first key that is set or if there is no default key in Artifactory.",
	},
}
