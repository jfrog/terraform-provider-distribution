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
