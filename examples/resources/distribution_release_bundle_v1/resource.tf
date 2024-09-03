resource "distribution_release_bundle_v1" "my-release-bundle-v1" {
  name = "my-release-bundle-v1"
  version = "1.0.0"
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