package distribution

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBundleBundleId_basic(t *testing.T) {
	// TODO: Implement acceptance test
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { /* testAccPreCheck(t) */ },
		ProtoV6ProviderFactories: nil, // TODO: set provider factories
		Steps: []resource.TestStep{
			{
				Config: testAccBundleBundleIdConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					// TODO: Add checks
				),
			},
		},
	})
}

func testAccBundleBundleIdConfig_basic() string {
	return `
resource "distribution_bundle_id" "test" {
}
`
}
