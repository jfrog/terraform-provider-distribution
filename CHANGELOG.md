## 1.3.1 (August 20, 2026). Tested on Artifactory 7.164.0 with Terraform 1.15.9 and OpenTofu 1.12.6

SECURITY:

* Remediate CVE-2026-39821.
* Remediate CVE-2026-56865.
* Remediate CVE-2026-56864.
* Remediate CVE-2026-33818.
* Remediate CVE-2026-46600.
* Remediate CVE-2026-56862.
* Remediate CVE-2026-56859.
* Remediate CVE-2026-56860.
* Remediate CVE-2026-56858.
* Remediate CVE-2026-56853.
* Remediate CVE-2026-25680.
* Remediate CVE-2026-42506.
* Remediate CVE-2026-42502.
* Remediate CVE-2026-25681.
* Remediate CVE-2026-27136.
* Remediate CVE-2026-46595.
* Remediate CVE-2026-42508.
* Remediate CVE-2026-39834.
* Remediate CVE-2026-39833.
* Remediate CVE-2026-39832.
* Remediate CVE-2026-39831.
* Remediate CVE-2026-39830.
* Remediate CVE-2026-39829.
* Remediate CVE-2026-46597.
* Remediate CVE-2026-39828.
* Remediate CVE-2026-39827.
* Remediate CVE-2026-39835.
* Remediate CVE-2026-46598.
* Remediate CVE-2025-47914.
* Remediate CVE-2025-58181.
* Remediate CVE-2026-1229.

## 1.3.0 (October 13, 2025). Tested on Artifactory 7.124.1 with Terraform 1.13.3 and OpenTofu 1.10.6

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