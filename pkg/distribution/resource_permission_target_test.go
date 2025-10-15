package distribution_test

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jfrog/terraform-provider-shared/testutil"
	"github.com/jfrog/terraform-provider-shared/util"
)

// generateRandomName creates a safe random name for testing
func generateRandomName(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s-%d", prefix, rand.Intn(100000))
}

func TestAccPermissionTarget_basic(t *testing.T) {
	_, fqrn, resourceName := testutil.MkNames("test-permission-target", "distribution_permission_target")
	userName := generateRandomName("test-user")
	groupName := generateRandomName("test-group")

	const template = `

	resource "artifactory_managed_user" "test-user" {
		name     = "{{ .userName }}"
		password = "Password1!"
		email    = "test@tempurl.org"
	}

	resource "platform_group" "testgroup" {
		name                       = "{{ .groupName }}"
		description 	           = "Test group"
		auto_join                  = true
		admin_privileges           = false
		use_group_members_resource = false
	}

	resource "distribution_permission_target" "{{ .name }}" {
		name        = "{{ .name }}"
		resource_type = "destination"
		distribution_destinations = [{
			site_name     = "*"
			city_name     = "*"
			country_codes = ["*"]
		}]
		principals = {
			users = {
				"{{ .userName }}" = ["d", "x"]
			}
			groups = {
				"{{ .groupName }}" = ["x"]
			}
		}
		depends_on = [
			artifactory_managed_user.test-user,
			platform_group.testgroup
		]			
	}`

	testData := map[string]string{
		"name":      resourceName,
		"userName":  userName,
		"groupName": groupName,
	}

	config := util.ExecuteTemplate("TestAccPermissionTarget_basic", template, testData)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviders(),
		ExternalProviders: map[string]resource.ExternalProvider{
			"platform": {
				Source: "jfrog/platform",
			},
			"artifactory": {
				Source: "jfrog/artifactory",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fqrn, "name", testData["name"]),
					resource.TestCheckResourceAttr(fqrn, "resource_type", "destination"),
					resource.TestCheckResourceAttr(fqrn, "distribution_destinations.0.site_name", "*"),
					resource.TestCheckResourceAttr(fqrn, "distribution_destinations.0.city_name", "*"),
					resource.TestCheckResourceAttr(fqrn, "distribution_destinations.0.country_codes.0", "*"),
					resource.TestCheckResourceAttr(fqrn, "principals.users."+userName+".0", "d"),
					resource.TestCheckResourceAttr(fqrn, "principals.users."+userName+".1", "x"),
					resource.TestCheckResourceAttr(fqrn, "principals.groups."+groupName+".0", "x"),
				),
			},
			{
				ResourceName:                         fqrn,
				ImportState:                          true,
				ImportStateId:                        testData["name"],
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

// Test cases for validation rules

func TestAccPermissionTarget_InvalidResourceType(t *testing.T) {
	_, _, resourceName := testutil.MkNames("test-permission-target", "distribution_permission_target")

	const template = `
	resource "distribution_permission_target" "{{ .name }}" {
		name        = "{{ .name }}"
		resource_type = "invalid-type"
		distribution_destinations = [{
			site_name     = "*"
			city_name     = "*"
			country_codes = ["*"]
		}]
		principals = {
			users = {
				"test-user" = ["x", "d"]
			}
		}
	}`

	testData := map[string]string{
		"name": resourceName,
	}

	config := util.ExecuteTemplate("TestAccPermissionTarget_InvalidResourceType", template, testData)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`Attribute resource_type value must be one of: \["destination"\], got:`),
			},
		},
	})
}

func TestAccPermissionTarget_MissingPrincipals(t *testing.T) {
	_, _, resourceName := testutil.MkNames("test-permission-target", "distribution_permission_target")

	const template = `
	resource "distribution_permission_target" "{{ .name }}" {
		name        = "{{ .name }}"
		resource_type = "destination"
		distribution_destinations = [{
			site_name     = "*"
			city_name     = "*"
			country_codes = ["*"]
		}]
		principals = {}
	}`

	testData := map[string]string{
		"name": resourceName,
	}

	config := util.ExecuteTemplate("TestAccPermissionTarget_MissingPrincipals", template, testData)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`lists.*users, groups`),
			},
		},
	})
}

func TestAccPermissionTarget_NoPrincipals(t *testing.T) {
	_, _, resourceName := testutil.MkNames("test-permission-target", "distribution_permission_target")

	const template = `
	resource "distribution_permission_target" "{{ .name }}" {
		name        = "{{ .name }}"
		resource_type = "destination"
		distribution_destinations = [{
			site_name     = "*"
			city_name     = "*"
			country_codes = ["*"]
		}]
	}`

	testData := map[string]string{
		"name": resourceName,
	}

	config := util.ExecuteTemplate("TestAccPermissionTarget_NoPrincipals", template, testData)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`The argument "principals" is required`),
			},
		},
	})
}

func TestAccPermissionTarget_MissingDistributionDestinations(t *testing.T) {
	_, _, resourceName := testutil.MkNames("test-permission-target", "distribution_permission_target")

	const template = `
	resource "distribution_permission_target" "{{ .name }}" {
		name        = "{{ .name }}"
		resource_type = "destination"
		principals = {
			users = {
				"test-user" = ["x", "d"]
			}
		}
	}`

	testData := map[string]string{
		"name": resourceName,
	}

	config := util.ExecuteTemplate("TestAccPermissionTarget_MissingDistributionDestinations", template, testData)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`The argument "distribution_destinations" is required`),
			},
		},
	})
}

func TestAccPermissionTarget_EmptyDistributionDestinations(t *testing.T) {
	_, _, resourceName := testutil.MkNames("test-permission-target", "distribution_permission_target")

	const template = `
	resource "distribution_permission_target" "{{ .name }}" {
		name        = "{{ .name }}"
		resource_type = "destination"
		distribution_destinations = []
		principals = {
			users = {
				"test-user" = ["x", "d"]
			}
		}
	}`

	testData := map[string]string{
		"name": resourceName,
	}

	config := util.ExecuteTemplate("TestAccPermissionTarget_EmptyDistributionDestinations", template, testData)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`list must contain at least 1 elements`),
			},
		},
	})
}

func TestAccPermissionTarget_MissingDestinationFields(t *testing.T) {
	_, _, resourceName := testutil.MkNames("test-permission-target", "distribution_permission_target")

	const template = `
	resource "distribution_permission_target" "{{ .name }}" {
		name        = "{{ .name }}"
		resource_type = "destination"
		distribution_destinations = [{
			site_name = "*"
			# Missing city_name and country_codes
		}]
		principals = {
			users = {
				"test-user" = ["x", "d"]
			}
		}
	}`

	testData := map[string]string{
		"name": resourceName,
	}

	config := util.ExecuteTemplate("TestAccPermissionTarget_MissingDestinationFields", template, testData)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`attributes "city_name" and "country_codes" are required`),
			},
		},
	})
}

func TestAccPermissionTarget_MissingCountryCodes(t *testing.T) {
	_, _, resourceName := testutil.MkNames("test-permission-target", "distribution_permission_target")

	const template = `
	resource "distribution_permission_target" "{{ .name }}" {
		name        = "{{ .name }}"
		resource_type = "destination"
		distribution_destinations = [{
			site_name = "*"
			city_name = "*"
			# Missing country_codes
		}]
		principals = {
			users = {
				"test-user" = ["x", "d"]
			}
		}
	}`

	testData := map[string]string{
		"name": resourceName,
	}

	config := util.ExecuteTemplate("TestAccPermissionTarget_MissingCountryCodes", template, testData)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`attribute "country_codes" is required`),
			},
		},
	})
}

func TestAccPermissionTarget_ValidWithOnlyUsers(t *testing.T) {
	_, fqrn, resourceName := testutil.MkNames("test-permission-target", "distribution_permission_target")
	userName := generateRandomName("test-user")

	const template = `
	resource "artifactory_managed_user" "test-user" {
		name     = "{{ .userName }}"
		password = "Password1!"
		email    = "test@tempurl.org"
		lifecycle {
			ignore_changes = [name]
		}
	}

	resource "distribution_permission_target" "{{ .name }}" {
		name        = "{{ .name }}"
		resource_type = "destination"
		distribution_destinations = [{
			site_name     = "*"
			city_name     = "*"
			country_codes = ["*"]
		}]
		principals = {
			users = {
				"{{ .userName }}" = ["d", "x"]
			}
		}
		depends_on = [
			artifactory_managed_user.test-user
		]
	}`

	testData := map[string]string{
		"name":     resourceName,
		"userName": userName,
	}

	config := util.ExecuteTemplate("TestAccPermissionTarget_ValidWithOnlyUsers", template, testData)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviders(),
		ExternalProviders: map[string]resource.ExternalProvider{
			"artifactory": {
				Source: "jfrog/artifactory",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fqrn, "name", testData["name"]),
					resource.TestCheckResourceAttr(fqrn, "resource_type", "destination"),
					resource.TestCheckResourceAttr(fqrn, "principals.users."+userName+".0", "d"),
					resource.TestCheckResourceAttr(fqrn, "principals.users."+userName+".1", "x"),
				),
			},
		},
	})
}

func TestAccPermissionTarget_ValidWithOnlyGroups(t *testing.T) {
	_, fqrn, resourceName := testutil.MkNames("test-permission-target", "distribution_permission_target")
	groupName := generateRandomName("test-group")

	const template = `
	resource "platform_group" "testgroup" {
		name                       = "{{ .groupName }}"
		description 	           = "Test group"
		auto_join                  = true
		admin_privileges           = false
		use_group_members_resource = false
		lifecycle {
			ignore_changes = [name]
		}
	}

	resource "distribution_permission_target" "{{ .name }}" {
		name        = "{{ .name }}"
		resource_type = "destination"
		distribution_destinations = [{
			site_name     = "*"
			city_name     = "*"
			country_codes = ["*"]
		}]
		principals = {
			groups = {
				"{{ .groupName }}" = ["d", "x"]
			}
		}
		depends_on = [
			platform_group.testgroup
		]
	}`

	testData := map[string]string{
		"name":      resourceName,
		"groupName": groupName,
	}

	config := util.ExecuteTemplate("TestAccPermissionTarget_ValidWithOnlyGroups", template, testData)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviders(),
		ExternalProviders: map[string]resource.ExternalProvider{
			"platform": {
				Source: "jfrog/platform",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fqrn, "name", testData["name"]),
					resource.TestCheckResourceAttr(fqrn, "resource_type", "destination"),
					resource.TestCheckResourceAttr(fqrn, "principals.groups."+groupName+".0", "d"),
					resource.TestCheckResourceAttr(fqrn, "principals.groups."+groupName+".1", "x"),
				),
			},
		},
	})
}
