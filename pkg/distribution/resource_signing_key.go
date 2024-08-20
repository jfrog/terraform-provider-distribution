package distribution

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jfrog/terraform-provider-shared/util"
	utilfw "github.com/jfrog/terraform-provider-shared/util/fw"
	"github.com/samber/lo"
)

func NewSigningKeyResource() resource.Resource {
	return &SigningKeyResource{
		TypeName: "distribution_signing_key",
	}
}

type SigningKeyResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type SigningKeyResourceModel struct {
	SigningKeyCommmonResourceModel
	PublicKey  types.String `tfsdk:"public_key"`
	PrivateKey types.String `tfsdk:"private_key"`
	Passphrase types.String `tfsdk:"passphrase"`
}

func (m SigningKeyResourceModel) toAPIModel(_ context.Context, apiModel *SigningKeyPostRequestAPIModel) (diags diag.Diagnostics) {
	apiModel.Key = SigningKeyKeyPostRequestAPIModel{
		Alias:      m.Alias.ValueString(),
		PublicKey:  m.PublicKey.ValueString(),
		PrivateKey: m.PrivateKey.ValueString(),
		Passphrase: m.Passphrase.ValueString(),
	}
	apiModel.PropagateToEdgeNode = m.PropagateToEdgeNode.ValueBool()
	apiModel.FailOnPropagationFailure = m.FailOnPropagationFailure.ValueBool()
	apiModel.SetAsDefault = m.SetAsDefault.ValueBool()

	return
}

type SigningKeyPostRequestAPIModel struct {
	SigningKeyCommonPostRequestAPIModel
	Key SigningKeyKeyPostRequestAPIModel `json:"key"`
}

type SigningKeyKeyPostRequestAPIModel struct {
	Alias      string `json:"alias"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	Passphrase string `json:"passphrase"`
}

func (r *SigningKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *SigningKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: lo.Assign(
			commonSchemaAttributes,
			map[string]schema.Attribute{
				"public_key": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Description: "Public key",
				},
				"private_key": schema.StringAttribute{
					Required:  true,
					Sensitive: true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Description: "Private key",
				},
				"passphrase": schema.StringAttribute{
					Optional:  true,
					Sensitive: true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Description: "Passphrase for key",
				},
			},
		),
		MarkdownDescription: "This resource enables you to upload and distribute GPG keys to sign Release Bundle V1. For more information, see [GPG Signing](https://jfrog.com/help/r/jfrog-distribution-documentation/gpg-signing) and [REST API](https://jfrog.com/help/r/jfrog-rest-apis/signing-keys).",
	}
}

func (r *SigningKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}

func (r *SigningKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	go util.SendUsageResourceCreate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan SigningKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var signingKey SigningKeyPostRequestAPIModel
	resp.Diagnostics.Append(plan.toAPIModel(ctx, &signingKey)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result SigningKeyPostResponseAPIModel
	var postErr SigningKeyPostErrorAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParam("protocol", plan.Protocol.ValueString()).
		SetBody(signingKey).
		SetResult(&result).
		SetError(&postErr).
		Post(SigningKeysEndpoint)
	if err != nil {
		utilfw.UnableToCreateResourceError(resp, err.Error())
		return
	}

	if response.IsError() {
		utilfw.UnableToCreateResourceError(resp, postErr.String())
		return
	}

	tflog.Info(ctx, "Signing Key Create", map[string]interface{}{
		"report": result.Report.String(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SigningKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state SigningKeyResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var signingKey SigningKeyGetAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"protocol": state.Protocol.ValueString(),
			"alias":    state.Alias.ValueString(),
		}).
		SetResult(&signingKey).
		Get(SigningKeyEndpoint)
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

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SigningKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	go util.SendUsageResourceUpdate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan SigningKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state SigningKeyResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	signingKey := SigningKeyPutRequestAPIModel{
		NewAlias: plan.Alias.ValueString(),
	}

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"protocol": state.Protocol.ValueString(),
			"alias":    state.Alias.ValueString(),
		}).
		SetBody(signingKey).
		Put(SigningKeyEndpoint)
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

func (r *SigningKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	go util.SendUsageResourceDelete(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state SigningKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"protocol": state.Protocol.ValueString(),
			"alias":    state.Alias.ValueString(),
		}).
		Delete(SigningKeyEndpoint)

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
