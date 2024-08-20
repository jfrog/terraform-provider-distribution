package distribution

import (
	"context"
	"fmt"
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

func NewVaultSigningKeyResource() resource.Resource {
	return &VaultSigningKeyResource{
		TypeName: "distribution_vault_signing_key",
	}
}

type VaultSigningKeyResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type VaultSigningKeyResourceModel struct {
	SigningKeyCommmonResourceModel
	VaultID    types.String `tfsdk:"vault_id"`
	PublicKey  types.Object `tfsdk:"public_key"`
	PrivateKey types.Object `tfsdk:"private_key"`
}

func (m VaultSigningKeyResourceModel) toAPIModel(_ context.Context, apiModel *VaultSigningKeyPostRequestAPIModel) (diags diag.Diagnostics) {
	publicKeyAttrs := m.PublicKey.Attributes()
	privateAttrs := m.PrivateKey.Attributes()
	apiModel.Key = VaultSigningKeyKeyPostRequestAPIModel{
		VaultData: VaultSigningKeyKeyVaultDataPostRequestAPIModel{
			ID: m.VaultID.ValueString(),
			PublicKey: VaultSigningKeyVaultDataPathKeyRequestAPIModel{
				Path: publicKeyAttrs["path"].(types.String).ValueString(),
				Key:  publicKeyAttrs["key"].(types.String).ValueString(),
			},
			PrivateKey: VaultSigningKeyVaultDataPathKeyRequestAPIModel{
				Path: privateAttrs["path"].(types.String).ValueString(),
				Key:  privateAttrs["key"].(types.String).ValueString(),
			},
		},
	}
	apiModel.PropagateToEdgeNode = m.PropagateToEdgeNode.ValueBool()
	apiModel.FailOnPropagationFailure = m.FailOnPropagationFailure.ValueBool()
	apiModel.SetAsDefault = m.SetAsDefault.ValueBool()

	return
}

type VaultSigningKeyPostRequestAPIModel struct {
	SigningKeyCommonPostRequestAPIModel
	Key VaultSigningKeyKeyPostRequestAPIModel `json:"key"`
}

type VaultSigningKeyKeyPostRequestAPIModel struct {
	VaultData VaultSigningKeyKeyVaultDataPostRequestAPIModel `json:"vault_data"`
}

type VaultSigningKeyKeyVaultDataPostRequestAPIModel struct {
	ID         string                                         `json:"vault_id"`
	PublicKey  VaultSigningKeyVaultDataPathKeyRequestAPIModel `json:"public_key"`
	PrivateKey VaultSigningKeyVaultDataPathKeyRequestAPIModel `json:"private_key"`
}

type VaultSigningKeyVaultDataPathKeyRequestAPIModel struct {
	Path string `json:"path"`
	Key  string `json:"key"`
}

func (r *VaultSigningKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

func (r *VaultSigningKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: lo.Assign(
			commonSchemaAttributes,
			map[string]schema.Attribute{
				"alias": schema.StringAttribute{
					Computed: true,
				},
				"vault_id": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Description: "Name of the Vault integration in Artifactory",
				},
				"public_key": schema.SingleNestedAttribute{
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							MarkdownDescription: "Path to the key, e.g. `secret/my-key`",
						},
						"key": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							MarkdownDescription: "Field name of the key, e.g. `public`",
						},
					},
					Required: true,
				},
				"private_key": schema.SingleNestedAttribute{
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							MarkdownDescription: "Path to the key, e.g. `secret/my-key`",
						},
						"key": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							MarkdownDescription: "Field name of the key, e.g. `private`",
						},
					},
					Required: true,
				},
			},
		),
		MarkdownDescription: "This resource enables you to distribute GPG keys (store in HashiCorp Vault) to sign Release Bundle V1. For more information, see [GPG Signing](https://jfrog.com/help/r/jfrog-distribution-documentation/gpg-signing), [Vault integration](https://jfrog.com/help/r/jfrog-platform-administration-documentation/vault), and [REST API](https://jfrog.com/help/r/jfrog-rest-apis/signing-keys).",
	}
}

func (r *VaultSigningKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}

func (r *VaultSigningKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	go util.SendUsageResourceCreate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan VaultSigningKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var signingKey VaultSigningKeyPostRequestAPIModel
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

	tflog.Info(ctx, "Vault Signing Key Create", map[string]interface{}{
		"report": result.Report.String(),
	})

	successJPD, found := lo.Find(
		result.Report.Details,
		func(detail SigningKeyReportDetailPostResponseAPIModel) bool {
			return detail.Status == "SUCCESS"
		},
	)

	if !found {
		utilfw.UnableToCreateResourceError(resp, fmt.Sprintf("Failed to deploy signing key: %s", result.Report.Status))
		return
	}

	plan.Alias = types.StringValue(successJPD.KeyAlias)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VaultSigningKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state VaultSigningKeyResourceModel

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

func (r *VaultSigningKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	go util.SendUsageResourceUpdate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan VaultSigningKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state VaultSigningKeyResourceModel

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

func (r *VaultSigningKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	go util.SendUsageResourceDelete(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state VaultSigningKeyResourceModel

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
