---
layout: ""
page_title: "JFrog Distribution Provider"
description: |-
  The JFrog Distribution provider provides resources to interact with features from JFrog distribution.
---

# JFrog Distribution Provider

The [JFrog](https://jfrog.com/) Distribution provider is used to interact with the features from [JFrog Distribution REST API](https://jfrog.com/help/r/jfrog-rest-apis/distribution-rest-apis). The provider needs to be configured with the proper credentials before it can be used.

Links to documentation for specific resources can be found in the table of contents to the left.

## Example Usage

```terraform
terraform {
  required_providers {
    distribution = {
      source  = "jfrog/distribution"
      version = "1.0.0"
    }
  }
}

variable "jfrog_url" {
  type = string
  default = "http://localhost:8081"
}

provider "distribution" {
  url = "${var.jfrog_url}"
  // supply JFROG_ACCESS_TOKEN as env var
}

resource "distribution_signing_key" "my-gpg-signing-key" {
  protocol = "gpg"
  alias = "my-gpg-signing-key"

  private_key = <<EOF
-----BEGIN PGP PRIVATE KEY BLOCK-----

lQVYBGa6hqoBDADAM1mIF2+ibcES+nP/gA6lHyRSGQ9JThQgIe18I/hQUkkM+Uji
dmJJ0uNmmc5hk+1FpR2NmmPnjNEgiV3Yu79Y+duX2QQtbclF6Nx//Z4/9cUTx2Us
...
nXOyvPCOk/4h817dmp0JqJi8XIABA9v0Jm0F209h09acd5baNaCszwn0adRWwSdU
ahFaWeXyMrgXl7+aVfwrBQ6G9tSP3Di6SiKOAlw=
=/yQZ
-----END PGP PRIVATE KEY BLOCK-----
EOF
  public_key = <<EOF
-----BEGIN PGP PUBLIC KEY BLOCK-----

mQGNBGa6hqoBDADAM1mIF2+ibcES+nP/gA6lHyRSGQ9JThQgIe18I/hQUkkM+Uji
dmJJ0uNmmc5hk+1FpR2NmmPnjNEgiV3Yu79Y+duX2QQtbclF6Nx//Z4/9cUTx2Us
...
2/QmbQXbT2HT1px3lto1oKzPCfRp1FbBJ1RqEVpZ5fIyuBeXv5pV/CsFDob21I/c
OLpKIo4CXA==
=ieBG
-----END PGP PUBLIC KEY BLOCK-----
EOF

  passphrase = "my-secret-passphrase"
  propagate_to_edge_nodes = true
  fail_on_propagation_failure = true
  set_as_default = true
}
```

## Authentication

The JFrog Distribution provider supports for the following types of authentication:
* Scoped token
* Terraform Cloud OIDC provider

### Scoped Token

JFrog scoped tokens may be used via the HTTP Authorization header by providing the `access_token` field to the provider block. Getting this value from the environment is supported with the `JFROG_ACCESS_TOKEN` environment variable.

Usage:
```terraform
provider "distribution" {
  url = "myinstance.jfrog.io"
  access_token = "abc...xy"
}
```

### Terraform Cloud OIDC Provider

If you are using this provider on Terraform Cloud and wish to use dynamic credentials instead of static access token for authentication with JFrog platform, you can leverage Terraform as the OIDC provider.

To setup dynamic credentials, follow these steps:
1. Configure Terraform Cloud as a generic OIDC provider
2. Set environment variable in your Terraform Workspace
3. Setup Terraform Cloud in your configuration

During the provider start up, if it finds env var `TFC_WORKLOAD_IDENTITY_TOKEN` it will use this token with your JFrog instance to exchange for a short-live access token. If that is successful, the provider will use the access token for all subsequent API requests with the JFrog instance.

#### Configure Terraform Cloud as generic OIDC provider

Follow [confgure an OIDC integration](https://jfrog.com/help/r/jfrog-platform-administration-documentation/configure-an-oidc-integration). Enter a name for the provider, e.g. `terraform-cloud`. Use `https://app.terraform.io` for "Provider URL". Choose your own value for "Audience", e.g. `jfrog-terraform-cloud`.

Then [configure an identity mapping](https://jfrog.com/help/r/jfrog-platform-administration-documentation/configure-identity-mappings) with appropriate "Claims JSON" (e.g. `aud`, `sub` at minimum. See [Terraform Workload Identity - Configuring Trust with your Cloud Platform](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/dynamic-provider-credentials/workload-identity-tokens#configuring-trust-with-your-cloud-platform)), and select the "Token scope", "User", and "Service" as desired.

#### Set environment variable in your Terraform Workspace

In your workspace, add an environment variable `TFC_WORKLOAD_IDENTITY_AUDIENCE` with audience value (e.g. `jfrog-terraform-cloud`) from JFrog OIDC integration above. See [Manually Generating Workload Identity Tokens](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/dynamic-provider-credentials/manual-generation) for more details.

When a run starts on Terraform Cloud, it will create a workload identity token with the specified audience and assigns it to the environment variable `TFC_WORKLOAD_IDENTITY_TOKEN` for the provider to consume.

See [Generating Multiple Tokens](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/dynamic-provider-credentials/manual-generation#generating-multiple-tokens) on HCP Terraform for more details on using different tokens.

#### Setup Terraform Cloud in your configuration

Add `cloud` block to `terraform` block, and add `oidc_provider_name` attribute (from JFrog OIDC integration) to provider block:

```terraform
terraform {
  cloud {
    organization = "my-org"
    workspaces {
      name = "my-workspace"
    }
  }

  required_providers {
    platform = {
      source  = "jfrog/distribution"
      version = "1.0.0"
    }
  }
}

provider "platform" {
  url = "https://myinstance.jfrog.io"
  oidc_provider_name = "terraform-cloud"
  tfc_credential_tag_name = "JFROG"
}
```

**Note:** Ensure `access_token` attribute is not set

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `access_token` (String, Sensitive) This is a access token that can be given to you by your admin under `Platform Configuration -> User Management -> Access Tokens`. This can also be sourced from the `JFROG_ACCESS_TOKEN` environment variable.
- `oidc_provider_name` (String) OIDC provider name. See [Configure an OIDC Integration](https://jfrog.com/help/r/jfrog-platform-administration-documentation/configure-an-oidc-integration) for more details.
- `tfc_credential_tag_name` (String) Terraform Cloud Workload Identity Token tag name. Use for generating multiple TFC workload identity tokens. When set, the provider will attempt to use env var with this tag name as suffix. **Note:** this is case sensitive, so if set to `JFROG`, then env var `TFC_WORKLOAD_IDENTITY_TOKEN_JFROG` is used instead of `TFC_WORKLOAD_IDENTITY_TOKEN`. See [Generating Multiple Tokens](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/dynamic-provider-credentials/manual-generation#generating-multiple-tokens) on HCP Terraform for more details.
- `url` (String) JFrog Platform URL. This can also be sourced from the `JFROG_URL` environment variable.