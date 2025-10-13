## 1.3.0 (October 13, 2025).

FEATURES:

**New Resource:**
* `distribution_permission_target`

PR: [#29](https://github.com/jfrog/terraform-provider-distribution/pull/29)

## 1.2.0 (October 17, 2024). Tested on Artifactory 7.95.0 with Terraform 1.9.8 and OpenTofu 1.8.3

* provider: Add `tfc_credential_tag_name` configuration attribute to support use of different/[multiple Workload Identity Token in Terraform Cloud Platform](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/dynamic-provider-credentials/manual-generation#generating-multiple-tokens). Issue: [#68](https://github.com/jfrog/terraform-provider-shared/issues/68) PR: [#11](https://github.com/jfrog/terraform-provider-distribution/pull/11)

## 1.1.0 (September 4, 2024). Tested on Artifactory 7.94.1 with Terraform 1.9.5 and OpenTofu 1.8.1

FEATURES:

**New Resource:**
* `distribution_release_bundle_v1`

PR: [#6](https://github.com/jfrog/terraform-provider-distribution/pull/6)

## 1.0.0 (August 16, 2024). Tested on Artifactory 7.92.1 with Terraform 1.9.4 and OpenTofu 1.8.1

FEATURES:

**New Resource:**
* `distribution_signing_key`
* `distribution_vault_signing_key`

PR: [#2](https://github.com/jfrog/terraform-provider-distribution/pull/2)