package distribution

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExportReleaseBundle_basic(t *testing.T) {
	// TODO: Implement acceptance test
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { /* testAccPreCheck(t) */ },
		ProtoV6ProviderFactories: nil, // TODO: set provider factories
		Steps: []resource.TestStep{
			{
				Config: testAccExportReleaseBundleConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					// TODO: Add checks
				),
			},
		},
	})
}

func testAccExportReleaseBundleConfig_basic() string {
	return `
resource "distribution_release_bundle" "test" {
  name = "test-value"
  version = "test-value"
}
`
}
