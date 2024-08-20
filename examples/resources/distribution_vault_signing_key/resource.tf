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