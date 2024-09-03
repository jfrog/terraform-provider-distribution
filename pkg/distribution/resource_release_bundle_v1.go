package distribution

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jfrog/terraform-provider-shared/util"
	utilfw "github.com/jfrog/terraform-provider-shared/util/fw"
	"github.com/samber/lo"
)

const (
	ReleaseBundlesV1Endpoint = "distribution/api/v1/release_bundle"
	ReleaseBundleV1Endpoint  = "distribution/api/v1/release_bundle/{name}/{version}"
)

func NewReleaseBundleV1Resource() resource.Resource {
	return &ReleaseBundleV1Resource{
		TypeName: "distribution_release_bundle_v1",
	}
}

type ReleaseBundleV1Resource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type ReleaseBundleV1ResourceModel struct {
	Name              types.String `tfsdk:"name"`
	Version           types.String `tfsdk:"version"`
	GPGPassphase      types.String `tfsdk:"gpg_passphase"`
	DryRun            types.Bool   `tfsdk:"dry_run"`
	SignImmediately   types.Bool   `tfsdk:"sign_immediately"`
	StoringRepository types.String `tfsdk:"storing_repository"`
	Description       types.String `tfsdk:"description"`
	ReleaseNotes      types.Object `tfsdk:"release_notes"`
	Spec              types.Object `tfsdk:"spec"`
	State             types.String `tfsdk:"state"`
	Created           types.String `tfsdk:"created"`
	CreatedBy         types.String `tfsdk:"created_by"`
	DistributedBy     types.String `tfsdk:"distributed_by"`
	Artifacts         types.Set    `tfsdk:"artifacts"`
	ArtifactsSize     types.Int64  `tfsdk:"artifacts_size"`
	Archived          types.Bool   `tfsdk:"archived"`
}

func (m ReleaseBundleV1ResourceModel) toAPIModel(ctx context.Context, apiModel *ReleaseBundleV1APIModel) (diags diag.Diagnostics) {
	apiModel.Name = m.Name.ValueString()
	apiModel.Version = m.Version.ValueString()
	apiModel.DryRun = m.DryRun.ValueBool()
	apiModel.SignImmediately = m.SignImmediately.ValueBool()
	apiModel.StoringRepository = m.StoringRepository.ValueString()
	apiModel.Description = m.Description.ValueString()

	releaseNotesAttrs := m.ReleaseNotes.Attributes()
	apiModel.ReleaseNotes = ReleaseBundleV1ReleaseNotesAPIModel{
		Syntax:  releaseNotesAttrs["syntax"].(types.String).ValueString(),
		Content: releaseNotesAttrs["content"].(types.String).ValueString(),
	}

	specAttrs := m.Spec.Attributes()
	queries := lo.Map(
		specAttrs["queries"].(types.Set).Elements(),
		func(elem attr.Value, _ int) ReleaseBundleV1SpecQueryAPIModel {
			attrs := elem.(types.Object).Attributes()

			mappings := lo.Map(
				attrs["mappings"].(types.Set).Elements(),
				func(elem attr.Value, _ int) ReleaseBundleV1SpecQueryMappingAPIModel {
					attrs := elem.(types.Object).Attributes()

					return ReleaseBundleV1SpecQueryMappingAPIModel{
						Input:  attrs["input"].(types.String).ValueString(),
						Output: attrs["output"].(types.String).ValueString(),
					}
				},
			)

			addedProps := lo.Map(
				attrs["added_props"].(types.Set).Elements(),
				func(elem attr.Value, _ int) ReleaseBundleV1PropAPIModel {
					attrs := elem.(types.Object).Attributes()

					values := lo.Map(
						attrs["values"].(types.Set).Elements(),
						func(elem attr.Value, _ int) string {
							return elem.(types.String).ValueString()
						},
					)

					return ReleaseBundleV1PropAPIModel{
						Key:    attrs["key"].(types.String).ValueString(),
						Values: values,
					}
				},
			)

			var excludedPropsPatterns []string
			diags.Append(attrs["exclude_props_patterns"].(types.Set).ElementsAs(ctx, &excludedPropsPatterns, false)...)

			return ReleaseBundleV1SpecQueryAPIModel{
				AQL:                   attrs["aql"].(types.String).ValueString(),
				QueryName:             attrs["query_name"].(types.String).ValueString(),
				Mappings:              mappings,
				AddedProps:            addedProps,
				ExcludedPropsPatterns: excludedPropsPatterns,
			}
		},
	)
	apiModel.Spec = ReleaseBundleV1SpecAPIModel{
		Queries: queries,
	}

	return
}

var artifactsAttrType = map[string]attr.Type{
	"checksum":         types.StringType,
	"source_repo_path": types.StringType,
	"target_repo_path": types.StringType,
	"props": types.SetType{
		ElemType: artifactsPropsType,
	},
}

var artifactsObjectType = types.ObjectType{
	AttrTypes: artifactsAttrType,
}

var artifactsPropsAttrType = map[string]attr.Type{
	"key":    types.StringType,
	"values": types.SetType{ElemType: types.StringType},
}

var artifactsPropsType = types.ObjectType{
	AttrTypes: artifactsPropsAttrType,
}

var artifactsFromAPIModel = func(ctx context.Context, artifactsAPIModel []ReleaseBundleV1ArtifactAPIModel) (types.Set, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	artifacts := lo.Map(
		artifactsAPIModel,
		func(artifact ReleaseBundleV1ArtifactAPIModel, _ int) attr.Value {
			props := lo.Map(
				artifact.Props,
				func(prop ReleaseBundleV1PropAPIModel, _ int) attr.Value {
					values, d := types.SetValueFrom(ctx, types.StringType, prop.Values)
					if d.HasError() {
						diags.Append(d...)
					}

					p, d := types.ObjectValue(
						artifactsPropsAttrType,
						map[string]attr.Value{
							"key":    types.StringValue(prop.Key),
							"values": values,
						},
					)
					if d.HasError() {
						diags.Append(d...)
					}

					return p
				},
			)

			propsSet, d := types.SetValue(
				artifactsPropsType,
				props,
			)
			if d.HasError() {
				diags.Append(d...)
			}

			a, d := types.ObjectValue(
				artifactsAttrType,
				map[string]attr.Value{
					"checksum":         types.StringValue(artifact.Checksum),
					"source_repo_path": types.StringValue(artifact.SourceRepoPath),
					"target_repo_path": types.StringValue(artifact.TargetRepoPath),
					"props":            propsSet,
				},
			)
			if d.HasError() {
				diags.Append(d...)
			}

			return a
		},
	)

	artifactsSet, d := types.SetValue(
		artifactsObjectType,
		artifacts,
	)
	if d.HasError() {
		diags.Append(d...)
	}

	return artifactsSet, diags
}

func (m *ReleaseBundleV1ResourceModel) fromPostAPIModel(ctx context.Context, apiModel ReleaseBundleV1PostResponseAPIModel) (diags diag.Diagnostics) {
	m.StoringRepository = types.StringValue(apiModel.StoringRepository)
	m.State = types.StringValue(apiModel.State)
	m.Created = types.StringValue(apiModel.Created)
	m.CreatedBy = types.StringValue(apiModel.CreatedBy)

	m.DistributedBy = types.StringNull()
	if apiModel.DistributedBy != nil {
		m.DistributedBy = types.StringPointerValue(apiModel.DistributedBy)
	}

	artifactsSet, d := artifactsFromAPIModel(ctx, apiModel.Artifacts)
	if d.HasError() {
		diags.Append(d...)
	}
	m.Artifacts = artifactsSet

	m.ArtifactsSize = types.Int64Value(apiModel.ArtifactsSize)
	m.Archived = types.BoolValue(apiModel.Archived)

	return
}

var mappingsAttrType = map[string]attr.Type{
	"input":  types.StringType,
	"output": types.StringType,
}

var mappingsObjectType = types.ObjectType{
	AttrTypes: mappingsAttrType,
}

var addedPropsAttrType = map[string]attr.Type{
	"key":    types.StringType,
	"values": types.SetType{ElemType: types.StringType},
}

var addedPropsObjectType = types.ObjectType{
	AttrTypes: addedPropsAttrType,
}

var queriesAttrType = map[string]attr.Type{
	"aql":        types.StringType,
	"query_name": types.StringType,
	"mappings": types.SetType{
		ElemType: mappingsObjectType,
	},
	"added_props": types.SetType{
		ElemType: addedPropsObjectType,
	},
	"exclude_props_patterns": types.SetType{
		ElemType: types.StringType,
	},
}

var queriesObjectType = types.ObjectType{
	AttrTypes: queriesAttrType,
}

var specAttrType = map[string]attr.Type{
	"queries": types.SetType{
		ElemType: queriesObjectType,
	},
}

func (m *ReleaseBundleV1ResourceModel) fromGetAPIModel(ctx context.Context, apiModel ReleaseBundleV1GetAPIModel) (diags diag.Diagnostics) {
	m.Name = types.StringValue(apiModel.Name)
	m.Version = types.StringValue(apiModel.Version)
	m.StoringRepository = types.StringValue(apiModel.StoringRepository)
	m.State = types.StringValue(apiModel.State)
	m.Description = types.StringValue(apiModel.Description)

	releaseNotes, d := types.ObjectValue(
		map[string]attr.Type{
			"content": types.StringType,
			"syntax":  types.StringType,
		},
		map[string]attr.Value{
			"content": types.StringValue(apiModel.ReleaseNotes.Content),
			"syntax":  types.StringValue(apiModel.ReleaseNotes.Syntax),
		},
	)
	if d.HasError() {
		diags.Append(d...)
	}
	m.ReleaseNotes = releaseNotes

	m.Created = types.StringValue(apiModel.Created)
	m.CreatedBy = types.StringValue(apiModel.CreatedBy)

	m.DistributedBy = types.StringNull()
	if apiModel.DistributedBy != nil {
		m.DistributedBy = types.StringPointerValue(apiModel.DistributedBy)
	}

	artifactsSet, d := artifactsFromAPIModel(ctx, apiModel.Artifacts)
	if d.HasError() {
		diags.Append(d...)
	}
	m.Artifacts = artifactsSet

	m.ArtifactsSize = types.Int64Value(apiModel.ArtifactsSize)
	m.Archived = types.BoolValue(apiModel.Archived)

	queries := lo.Map(
		apiModel.Spec.Queries,
		func(query ReleaseBundleV1SpecQueryAPIModel, _ int) attr.Value {
			mappings := lo.Map(
				query.Mappings,
				func(mapping ReleaseBundleV1SpecQueryMappingAPIModel, _ int) attr.Value {
					m, d := types.ObjectValue(
						mappingsAttrType,
						map[string]attr.Value{
							"input":  types.StringValue(mapping.Input),
							"output": types.StringValue(mapping.Output),
						},
					)
					if d.HasError() {
						diags.Append(d...)
					}
					return m
				},
			)
			mappingsSet, d := types.SetValue(
				mappingsObjectType,
				mappings,
			)
			if d.HasError() {
				diags.Append(d...)
			}

			addedProps := lo.Map(
				query.AddedProps,
				func(addedProp ReleaseBundleV1PropAPIModel, _ int) attr.Value {
					values, d := types.SetValueFrom(ctx, types.StringType, addedProp.Values)
					if d.HasError() {
						diags.Append(d...)
					}
					p, d := types.ObjectValue(
						addedPropsAttrType,
						map[string]attr.Value{
							"key":    types.StringValue(addedProp.Key),
							"values": values,
						},
					)
					if d.HasError() {
						diags.Append(d...)
					}
					return p
				},
			)
			addedPropsSet, d := types.SetValue(
				addedPropsObjectType,
				addedProps,
			)
			if d.HasError() {
				diags.Append(d...)
			}

			excludePropsPatterns, d := types.SetValueFrom(ctx, types.StringType, query.ExcludedPropsPatterns)
			if d.HasError() {
				diags.Append(d...)
			}

			q, d := types.ObjectValue(
				queriesAttrType,
				map[string]attr.Value{
					"aql":                    types.StringValue(query.AQL),
					"query_name":             types.StringValue(query.QueryName),
					"mappings":               mappingsSet,
					"added_props":            addedPropsSet,
					"exclude_props_patterns": excludePropsPatterns,
				},
			)
			if d.HasError() {
				diags.Append(d...)
			}
			return q
		},
	)

	queriesSet, d := types.SetValue(
		queriesObjectType,
		queries,
	)
	if d.HasError() {
		diags.Append(d...)
	}

	spec, d := types.ObjectValue(
		specAttrType,
		map[string]attr.Value{
			"queries": queriesSet,
		},
	)
	if d.HasError() {
		diags.Append(d...)
	}
	m.Spec = spec

	return
}

type ReleaseBundleV1APIModel struct {
	Name              string                              `json:"name"`
	Version           string                              `json:"version"`
	DryRun            bool                                `json:"dry_run"`
	SignImmediately   bool                                `json:"sign_immediately"`
	StoringRepository string                              `json:"storing_repository,omitempty"`
	Description       string                              `json:"description"`
	ReleaseNotes      ReleaseBundleV1ReleaseNotesAPIModel `json:"release_notes"`
	Spec              ReleaseBundleV1SpecAPIModel         `json:"spec"`
}

type ReleaseBundleV1ReleaseNotesAPIModel struct {
	Syntax  string `json:"syntax"`
	Content string `json:"content"`
}

type ReleaseBundleV1SpecAPIModel struct {
	Queries []ReleaseBundleV1SpecQueryAPIModel `json:"queries"`
}

type ReleaseBundleV1SpecQueryAPIModel struct {
	AQL                   string                                    `json:"aql"`
	QueryName             string                                    `json:"query_name"`
	Mappings              []ReleaseBundleV1SpecQueryMappingAPIModel `json:"mappings"`
	AddedProps            []ReleaseBundleV1PropAPIModel             `json:"added_props"`
	ExcludedPropsPatterns []string                                  `json:"exclude_props_patterns"`
}

type ReleaseBundleV1SpecQueryMappingAPIModel struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type ReleaseBundleV1PostResponseAPIModel struct {
	StoringRepository string                            `json:"storing_repository"`
	State             string                            `json:"state"`
	Created           string                            `json:"created"`
	CreatedBy         string                            `json:"created_by"`
	DistributedBy     *string                           `json:"distributed_by,omitempty"`
	Artifacts         []ReleaseBundleV1ArtifactAPIModel `json:"artifacts"`
	ArtifactsSize     int64                             `json:"artifacts_size"`
	Archived          bool                              `json:"archived"`
}

type ReleaseBundleV1ArtifactAPIModel struct {
	Checksum       string                        `json:"checksum"`
	SourceRepoPath string                        `json:"sourceRepoPath"`
	TargetRepoPath string                        `json:"targetRepoPath"`
	Props          []ReleaseBundleV1PropAPIModel `json:"props"`
}

type ReleaseBundleV1PropAPIModel struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type ReleaseBundleV1GetAPIModel struct {
	Name              string                              `json:"name"`
	Version           string                              `json:"version"`
	StoringRepository string                              `json:"storing_repository,omitempty"`
	State             string                              `json:"state"`
	Description       string                              `json:"description"`
	ReleaseNotes      ReleaseBundleV1ReleaseNotesAPIModel `json:"release_notes"`
	Created           string                              `json:"created"`
	CreatedBy         string                              `json:"created_by"`
	DistributedBy     *string                             `json:"distributed_by,omitempty"`
	Artifacts         []ReleaseBundleV1ArtifactAPIModel   `json:"artifacts"`
	ArtifactsSize     int64                               `json:"artifacts_size"`
	Archived          bool                                `json:"archived"`
	Spec              ReleaseBundleV1SpecAPIModel         `json:"spec"`
}

func (r *ReleaseBundleV1Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName
}

var nameVersionRegexValidator = stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_\.\-:]+$`), "must begin with a letter or digit and consist only of letters, digits, underscores, periods, hyphens, and colons.")

func (r *ReleaseBundleV1Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
					nameVersionRegexValidator,
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Release bundle name. Must begin with a letter or digit and consist only of letters, digits, underscores, periods, hyphens, and colons.",
			},
			"version": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
					stringvalidator.NoneOf("LATEST"),
					nameVersionRegexValidator,
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Release bundle version. Must begin with a letter or digit and consist only of letters, digits, underscores, periods, hyphens, and colons. The string `LATEST` is prohibited.",
			},
			"gpg_passphase": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Description: "Passphrase for the signing key, if applicable",
			},
			"dry_run": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "When set to `true`, only parses and validates.",
			},
			"sign_immediately": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "When set to `true`, automatically signs the release bundle version.",
			},
			"storing_repository": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Description: "A repository name at source Artifactory to store release bundle artifacts in. If not provided, Artifactory will use the default one (requires Artifactory 6.5 or later).",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Description: "Description of the release bundle.",
			},
			"release_notes": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"syntax": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("plain_text"),
						Validators: []validator.String{
							stringvalidator.OneOf("markdown", "asciidoc", "plain_text"),
						},
						MarkdownDescription: "The syntax for the release notes. Options include: `markdown`, `asciidoc`, `plain_text` (default).",
					},
					"content": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
						MarkdownDescription: "The content of the release notes.",
					},
				},
				Optional:    true,
				Description: "Describes the release notes for the release bundle version.",
			},
			"spec": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"queries": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"aql": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.LengthAtLeast(1),
									},
									Description: "AQL query for gathering the artifacts from Artifactory.",
								},
								"query_name": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.LengthBetween(2, 32),
										stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_\-\.:]+$`), "must start with alphabetic character followed by an alphanumeric or '_-.:' characters only"),
									},
									MarkdownDescription: "A name to be used when displaying the query object. Note that the release bundle query name length must be between 2 and 32 characters long and must start with alphabetic character followed by an alphanumeric or `_-.:` characters only.",
								},
								"added_props": schema.SetNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
												Description: "Property key to be created or updated on the distributed artifacts.",
											},
											"values": schema.SetAttribute{
												ElementType: types.StringType,
												Optional:    true,
												Description: "List of values to be added to the property key after distribution of the release bundle",
											},
										},
									},
									Optional:    true,
									Description: "List of added properties which will be added to the artifacts after distribution of the release bundle",
								},
								"mappings": schema.SetNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"input": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
												Description: "Regex matcher for artifact paths.",
											},
											"output": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.LengthAtLeast(1),
												},
												Description: "Replacement for artifact paths matched by the `input` matcher. Capture groups can be used as `$1`.",
											},
										},
									},
									Optional:    true,
									Description: "List of mappings, which are applied to the artifact paths after distribution of the release bundle",
								},
								"exclude_props_patterns": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Description: "List of patterns for Properties keys to exclude after distribution of the release bundle. This will not have an effect on the `added_props` attribute.",
								},
							},
						},
						Required:    true,
						Description: "List of query objects to gather artifacts by.",
					},
				},
				Required:    true,
				Description: "Describes the specification by artifacts are gathered and distributed in this release bundle.",
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"created_by": schema.StringAttribute{
				Computed: true,
			},
			"distributed_by": schema.StringAttribute{
				Computed: true,
			},
			"artifacts": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"checksum": schema.StringAttribute{
							Computed: true,
						},
						"source_repo_path": schema.StringAttribute{
							Computed: true,
						},
						"target_repo_path": schema.StringAttribute{
							Computed: true,
						},
						"props": schema.SetNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Computed: true,
									},
									"values": schema.SetAttribute{
										ElementType: types.StringType,
										Computed:    true,
									},
								},
							},
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"artifacts_size": schema.Int64Attribute{
				Computed: true,
			},
			"archived": schema.BoolAttribute{
				Computed: true,
			},
		},
		MarkdownDescription: "This resource enables you to create a new Release Bundle V1 version. For more information, see [Create Release Bundle V1](https://jfrog.com/help/r/jfrog-distribution-documentation/create-release-bundles-v1) and [REST API](https://jfrog.com/help/r/jfrog-rest-apis/release-bundles-v1).\n\n" +
			"~>User must have matching Release Bundle writer permissions.",
	}
}

func (r *ReleaseBundleV1Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}

func (r *ReleaseBundleV1Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	go util.SendUsageResourceCreate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan ReleaseBundleV1ResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var releaseBundle ReleaseBundleV1APIModel
	resp.Diagnostics.Append(plan.toAPIModel(ctx, &releaseBundle)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result ReleaseBundleV1PostResponseAPIModel

	request := r.ProviderData.Client.R()

	if !plan.GPGPassphase.IsNull() {
		request.SetHeader("X-GPG-PASSPHRASE", plan.GPGPassphase.ValueString())
	}

	response, err := request.
		SetBody(releaseBundle).
		SetResult(&result).
		Post(ReleaseBundlesV1Endpoint)

	if err != nil {
		utilfw.UnableToCreateResourceError(resp, err.Error())
		return
	}

	if response.IsError() {
		utilfw.UnableToCreateResourceError(resp, response.String())
		return
	}

	resp.Diagnostics.Append(plan.fromPostAPIModel(ctx, result)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ReleaseBundleV1Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state ReleaseBundleV1ResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var releaseBundle ReleaseBundleV1GetAPIModel

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"name":    state.Name.ValueString(),
			"version": state.Version.ValueString(),
		}).
		SetQueryParam("format", "json").
		SetResult(&releaseBundle).
		Get(ReleaseBundleV1Endpoint)

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

	resp.Diagnostics.Append(state.fromGetAPIModel(ctx, releaseBundle)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ReleaseBundleV1Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	go util.SendUsageResourceUpdate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan ReleaseBundleV1ResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var releaseBundle ReleaseBundleV1APIModel
	resp.Diagnostics.Append(plan.toAPIModel(ctx, &releaseBundle)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result ReleaseBundleV1PostResponseAPIModel

	request := r.ProviderData.Client.R()

	if !plan.GPGPassphase.IsNull() {
		request.SetHeader("X-GPG-PASSPHRASE", plan.GPGPassphase.ValueString())
	}

	response, err := request.
		SetPathParams(map[string]string{
			"name":    plan.Name.ValueString(),
			"version": plan.Version.ValueString(),
		}).
		SetBody(releaseBundle).
		SetResult(&result).
		Put(ReleaseBundleV1Endpoint)

	if err != nil {
		utilfw.UnableToUpdateResourceError(resp, err.Error())
		return
	}

	if response.IsError() {
		utilfw.UnableToUpdateResourceError(resp, response.String())
		return
	}

	resp.Diagnostics.Append(plan.fromPostAPIModel(ctx, result)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ReleaseBundleV1Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	go util.SendUsageResourceDelete(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state ReleaseBundleV1ResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	response, err := r.ProviderData.Client.R().
		SetPathParams(map[string]string{
			"name":    state.Name.ValueString(),
			"version": state.Version.ValueString(),
		}).
		Delete(ReleaseBundleV1Endpoint)

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

	// If the logic reaches here, it implicitly succeeded and will remove
	// the resource from state if there are no other errors.
}

// ImportState imports the resource into the Terraform state.
func (r *ReleaseBundleV1Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			"Expected name:version",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("version"), parts[1])...)
}
