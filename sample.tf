terraform {
  required_providers {
    distribution = {
      source  = "jfrog/distribution"
    }
  }
}

provider "distribution" {
  url          = "artifactory.site.com"
  access_token = "abc..xy"
  // Also user can supply the following env vars:
  // JFROG_URL
  // JFROG_ACCESS_TOKEN
}

resource "distribution_permission_target" "example_permission4" {
  name        = "example-permission-4"
  resource_type = "destination"
  distribution_destinations = [
    {
      site_name     = "*"
      city_name     = "*"
      country_codes = ["*"]
    }
  ]
  
  principals = {
    users = {
      "test1" = ["x"]
      "test2" = ["d","x"]
    }
    groups = {
      "grp1" = ["x"]
      "grp2" = ["d","x"]
    }
  }
}

resource "distribution_release_bundle_v1" "my-release-bundle-v1" {
  name = "my-release-bundle-v1"
  version = "1.2.1"
  sign_immediately = true
  description = "My description"

  release_notes = {
    syntax = "plain_text"
    content = "My release notes"
  }

  spec = {
    queries = [{
      aql = "items.find({ \"repo\" : \"example-repo-local\" })"
      query_name: "query-1"

      mappings = [{
        input = "original_repository/(.*)"
        output = "new_repository/$1"
      }]

      added_props = [{
        key = "my-key"
        values = ["my-value"]
      }]

      exclude_props_patterns = [
        "my-prop-*"
      ]
    }]
  }
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

resource "distribution_vault_signing_key" "my-vault-gpg-signing-key" {
  protocol = "gpg"
  vault_id = "my-vault-integration"

  public_key = {
    path = "kv/public/path"
    key = "public"
  }

  private_key = {
    path = "kv/private/path"
    key = "private"
  }

  propagate_to_edge_nodes = true
  fail_on_propagation_failure = true
  set_as_default = true
}