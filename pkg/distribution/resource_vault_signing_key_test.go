package distribution_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jfrog/terraform-provider-shared/testutil"
	"github.com/jfrog/terraform-provider-shared/util"
)

// To execute these tests successfully, you'll need:
// 1. Self-hosted Artifactory instance running with TLS enabled (see `security.tls` in `access.config.patch.yml`), OR
// 1.b A SaaS instance
// 2. Vault server running and configured with role (app role ID and secret), and permission
// 3. GPG public and private keys stored in Vault kv secret
// 4. Artifactory configured with Vault
// 5. Set env var JFROG_VAULT_ID=<Vault integration name in Artifactory>
//
// See https://github.com/jfrog/terraform-provider-distribution/wiki/How-to-setup-environment-for-testing-signing-key-from-Vault
func TestAccVaultSigningKey_full(t *testing.T) {
	vaultID := os.Getenv("JFROG_VAULT_ID")
	if vaultID == "" {
		t.Skipf("env var JFROG_VAULT_ID is not set.")
	}

	_, fqrn, resourceName := testutil.MkNames("test-vault-signing-key", "distribution_vault_signing_key")

	const template = `
	resource "distribution_vault_signing_key" "{{ .name }}" {
		protocol = "gpg"
		vault_id = "{{ .vault_id }}"

		public_key = {
			path = "{{ .vault_secret_path }}"
			key = "public"
		}

		private_key = {
			path = "{{ .vault_secret_path }}"
			key = "private"
		}
		
		propagate_to_edge_nodes = true
		fail_on_propagation_failure = true
		set_as_default = true
	}`

	testData := map[string]string{
		"name":              resourceName,
		"vault_id":          vaultID,
		"vault_secret_path": "secret/signing_key",
	}

	config := util.ExecuteTemplate("TestAccVaultSigningKey_gpg_full", template, testData)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fqrn, "protocol", "gpg"),
					resource.TestCheckResourceAttrSet(fqrn, "alias"),
					resource.TestCheckResourceAttr(fqrn, "public_key.path", testData["vault_secret_path"]),
					resource.TestCheckResourceAttr(fqrn, "public_key.key", "public"),
					resource.TestCheckResourceAttr(fqrn, "private_key.path", testData["vault_secret_path"]),
					resource.TestCheckResourceAttr(fqrn, "private_key.key", "private"),
					resource.TestCheckResourceAttr(fqrn, "propagate_to_edge_nodes", "true"),
					resource.TestCheckResourceAttr(fqrn, "fail_on_propagation_failure", "true"),
					resource.TestCheckResourceAttr(fqrn, "set_as_default", "true"),
				),
			},
		},
	})
}
