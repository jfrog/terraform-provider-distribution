package distribution

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccReleaseBundleStatus_basic(t *testing.T) {
	// TODO: Implement acceptance test
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { /* testAccPreCheck(t) */ },
		ProtoV6ProviderFactories: nil, // TODO: set provider factories
		Steps: []resource.TestStep{
			{
				Config: testAccReleaseBundleStatusConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					// TODO: Add checks
				),
			},
		},
	})
}

func testAccReleaseBundleStatusConfig_basic() string {
	return `
resource "distribution_status" "test" {
  name = "test-value"
  version = "test-value"
}
`
}
