package distribution_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jfrog/terraform-provider-shared/testutil"
	"github.com/jfrog/terraform-provider-shared/util"
)

func TestAccSigningKey_gpg_full(t *testing.T) {
	_, fqrn, resourceName := testutil.MkNames("test-signing-key", "distribution_signing_key")

	const template = `
	resource "distribution_signing_key" "{{ .name }}" {
		protocol = "gpg"
		alias = "{{ .alias }}"

		private_key = <<EOF
-----BEGIN PGP PRIVATE KEY BLOCK-----

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
ahFaWeXyMrgXl7+aVfwrBQ6G9tSP3Di6SiKOAlw=
=/yQZ
-----END PGP PRIVATE KEY BLOCK-----
EOF
		public_key = <<EOF
-----BEGIN PGP PUBLIC KEY BLOCK-----

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
OLpKIo4CXA==
=ieBG
-----END PGP PUBLIC KEY BLOCK-----
EOF

		passphrase = "{{ .passphrase }}"
		
		propagate_to_edge_nodes = true
		fail_on_propagation_failure = true
		set_as_default = true
	}`

	testData := map[string]string{
		"name":       resourceName,
		"alias":      resourceName,
		"passphrase": "password",
	}

	config := util.ExecuteTemplate("TestAccSigningKey_gpg_full", template, testData)

	updatedTestData := map[string]string{
		"name":       resourceName,
		"alias":      fmt.Sprintf("%s-2", resourceName),
		"passphrase": "password",
	}

	updatedConfig := util.ExecuteTemplate("TestAccSigningKey_gpg_full", template, updatedTestData)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fqrn, "protocol", "gpg"),
					resource.TestCheckResourceAttr(fqrn, "alias", testData["alias"]),
					resource.TestCheckResourceAttrSet(fqrn, "public_key"),
					resource.TestCheckResourceAttrSet(fqrn, "private_key"),
					resource.TestCheckResourceAttr(fqrn, "propagate_to_edge_nodes", "true"),
					resource.TestCheckResourceAttr(fqrn, "fail_on_propagation_failure", "true"),
					resource.TestCheckResourceAttr(fqrn, "set_as_default", "true"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fqrn, "alias", updatedTestData["alias"]),
				),
			},
		},
	})
}

func TestAccSigningKey_pgp_full(t *testing.T) {
	_, fqrn, resourceName := testutil.MkNames("test-signing-key", "distribution_signing_key")

	const template = `
	resource "distribution_signing_key" "{{ .name }}" {
		protocol = "pgp"
		alias = "{{ .alias }}"

		private_key = <<EOF
-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: Keybase OpenPGP v1.0.0
Comment: https://keybase.io/crypto

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
zoQkekoL4H4GU8QB6GY=
=JQxa
-----END PGP PRIVATE KEY BLOCK-----
EOF
		public_key = <<EOF
-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: Keybase OpenPGP v1.0.0
Comment: https://keybase.io/crypto

xo0EYGrVNAEEAPD3YDt0qP8kSV8bnmqVP5XDPoN40gEpUGtDLjAn6d+cRMeNGaru
6H0bdgwQpND8Gz9Qx2pCNSxlWDZpY1fCvRQ174iGjvO/3527f148cgKNZtwLsKrZ
laW8z3tB2LuCM2e97ijX+lzRf7YJUXU3pOfoCFWpOPoRg1CHV0NyHl0VABEBAAHN
FmFsYW4gPGFsYW5uQGpmcm9nLmNvbT7CrQQTAQoAFwUCYGrVNAIbLwMLCQcDFQoI
Ah4BAheAAAoJENzR2QJlA6glZmsD/iqhnNFy1Elj3hGL0HaEzeb+KDpcSL/L5a/8
WIGCQFeLcEn9lC+68b/eERKGIoXJ7z8HfPDFNRTKvomKIdAqFiAeDAUUD0B82rsx
xDf8USnTwJlnd0bPe9nxgXYcrwioEYbPVYGl3jima/KQrbW8XlKyiypy4Nd66Wcn
TuM6PwRFzo0EYGrVNAEEANVNINyfCQ+y1haaaAJ0uCgx3dW52LwcZfvOP6i798WZ
dyGA+WSUCEcrklUwZ595E2dNkNKptksftwSeQ0+EH5S1ZlEaq2YUv8fCx32F1ckh
D3eHaCKRxTPx/zbb96q4ruEGKhOBXceid3o341HbtGVKi8VjBx3XNukskQ+EOvgt
ABEBAAHCwIMEGAEKAA8FAmBq1TQFCQ8JnAACGy4AqAkQ3NHZAmUDqCWdIAQZAQoA
BgUCYGrVNAAKCRBddQ63FhKl6NmDBACqxC4lAnsCQERjs02LYAEAwVDhDf0rXxD0
H+hKDyxQZc80M7WIpXaBHmbs8ekJRnY7JHcer7sizDMdfkR3xB62jNGhc0XiW6nc
mlwvtWt3+E6AkObmWnocRy5ztTQI0gye0B3cPs2txE2fCs+WD7yLRnM3HqIAh83W
Cccvh0+uG96dBADlPbZ8g8q6bkeeT72gOi3OCN0A+Y8lUPifhrpSiI9xMpP3aomM
beZJB6fWjEzNoblQ9jUr/E54bF9jMr6L3uE4OJH9SYJ/HvqcKJC+1TFeQ9lXR7g7
MTdfxEvhMDhcsd/pYIgrzvDry+B2+jANW10R1yejT/C8QdlWIndDsEsaKs6NBGBq
1TQBBADJD6gVGTMGb/WsfnSaL5yA3AEMczhPqxD4FDC0vzqGG3XzKgxtmW8cJXls
NCK80F61daxJ72/UAmfxHbNP1qmHSiosEe+h1YZ/Zo3pN/LODzg9JrOs9A2xwjqE
nU9mS0jDz5oEQtr9K4+YKOdJvmuaN85ueBizfQ1TYRNuDtmbnQARAQABwsCDBBgB
CgAPBQJgatU0BQkPCZwAAhsuAKgJENzR2QJlA6glnSAEGQEKAAYFAmBq1TQACgkQ
ehctgYvYtmh38gP6A9lnQaLuVnTElJLy2XSDTqwWOcy/5J842S/xdQEsWUMXh4I5
mlotkZwkrdvXp8E/F3P8X7GbxhNAVZX+Xcm95V3g/kmP+Pq7PeUmoZR5LD8ppBfO
7v6XgaUhraUPAZl6lx4L5pYNCX9JBNUtQAG9xIoap4slvksdz5SN/BwSgV6qqwQA
tr4YTDXvLyoWwMFB2FjWcw4zwV+7yHwGzogKfGCQy5qVlDoQyWdkwwF1awyk5RIe
ZxwPZ2SDaiznOmZ+4LjR2NPmjnT96d9RKRtgEjkfW+a19BofrvEalS9wh/jkboea
d8rDu8wMbLAl77dq1c6dpJDgzoQkekoL4H4GU8QB6GY=
=fot9
-----END PGP PUBLIC KEY BLOCK-----
EOF

		passphrase = "{{ .passphrase }}"
		
		propagate_to_edge_nodes = true
		fail_on_propagation_failure = true
		set_as_default = true
	}`

	testData := map[string]string{
		"name":       resourceName,
		"alias":      resourceName,
		"passphrase": "password",
	}

	config := util.ExecuteTemplate("TestAccSigningKey_pgp_full", template, testData)

	updatedTestData := map[string]string{
		"name":       resourceName,
		"alias":      fmt.Sprintf("%s-2", resourceName),
		"passphrase": "password",
	}

	updatedConfig := util.ExecuteTemplate("TestAccSigningKey_pgp_full", template, updatedTestData)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fqrn, "protocol", "pgp"),
					resource.TestCheckResourceAttr(fqrn, "alias", testData["alias"]),
					resource.TestCheckResourceAttrSet(fqrn, "public_key"),
					resource.TestCheckResourceAttrSet(fqrn, "private_key"),
					resource.TestCheckResourceAttr(fqrn, "propagate_to_edge_nodes", "true"),
					resource.TestCheckResourceAttr(fqrn, "fail_on_propagation_failure", "true"),
					resource.TestCheckResourceAttr(fqrn, "set_as_default", "true"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fqrn, "alias", updatedTestData["alias"]),
				),
			},
		},
	})
}
