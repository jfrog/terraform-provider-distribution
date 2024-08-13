[![Terraform & OpenTofu Acceptance Tests](https://github.com/jfrog/terraform-provider-distribution/actions/workflows/acceptance-tests.yml/badge.svg)](https://github.com/jfrog/terraform-provider-distribution/actions/workflows/acceptance-tests.yml)

# Terraform Provider for JFrog Platform

## Quick Start

Create a new Terraform file with `distribution` resource.

### HCL Example

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

Initialize Terrform:
```sh
$ terraform init
```

Plan (or Apply):
```sh
$ terraform plan
```

Detailed documentation of the resource and attributes are on [Terraform Registry](https://registry.terraform.io/providers/jfrog/platform/latest/docs).

## Versioning

In general, this project follows [semver](https://semver.org/) as closely as we can for tagging releases of the package. We've adopted the following versioning policy:

* We increment the **major version** with any incompatible change to functionality, including changes to the exported Go API surface or behavior of the API.
* We increment the **minor version** with any backwards-compatible changes to functionality.
* We increment the **patch version** with any backwards-compatible bug fixes.

## Contributors

See the [contribution guide](CONTRIBUTIONS.md).

## License

Copyright (c) 2024 JFrog.

Apache 2.0 licensed, see [LICENSE][LICENSE] file.

[LICENSE]: ./LICENSE
