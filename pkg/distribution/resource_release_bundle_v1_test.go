package distribution_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jfrog/terraform-provider-shared/testutil"
	"github.com/jfrog/terraform-provider-shared/util"
)

func TestAccReleaseBundleV1_full(t *testing.T) {
	_, fqrn, resourceName := testutil.MkNames("test-release-bundle-v1", "distribution_release_bundle_v1")

	const template = `
	resource "distribution_release_bundle_v1" "{{ .name }}" {
		name = "{{ .name }}"
		version = "{{ .version }}"
		sign_immediately = {{ .sign_immediately }}
		description = "Test description"

		release_notes = {
			syntax = "plain_text"
			content = "test release notes"
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
					key = "test-key"
					values = ["test-value"]
				}]
				
				exclude_props_patterns = [
					"test-patterns"
				]
			}]
		}
	}`

	testData := map[string]string{
		"name":             resourceName,
		"version":          "1.0.0",
		"sign_immediately": "false",
	}

	config := util.ExecuteTemplate("TestAccReleaseBundleV1_full", template, testData)

	const updatedTemplate = `
	resource "distribution_release_bundle_v1" "{{ .name }}" {
		name = "{{ .name }}"
		version = "{{ .version }}"
		sign_immediately = {{ .sign_immediately }}
		description = "Test description"

		release_notes = {
			syntax = "plain_text"
			content = "test release notes"
		}

		spec = {
			queries = [{
				aql = "items.find({ \"repo\" : \"example-repo-local\" })"
				query_name: "query-1"

				mappings = [{
					input = "original_repository/(.*)"
					output = "new_repository/$1"
				}, {
					input = "(.*)/(.*)"
					output = "$1/new_folder/$2"
				}]

				added_props = [{
					key = "test-key"
					values = ["test-value"]
				}, {
					key = "test-key-2"
					values = ["test-value-2"]
				}]
				
				exclude_props_patterns = [
					"test-patterns",
					"test-patterns-2",
				]
			}]
		}
	}`

	updatedConfig := util.ExecuteTemplate("TestAccReleaseBundleV1_full", updatedTemplate, testData)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fqrn, "name", testData["name"]),
					resource.TestCheckResourceAttr(fqrn, "version", testData["version"]),
					resource.TestCheckResourceAttr(fqrn, "sign_immediately", testData["sign_immediately"]),
					resource.TestCheckResourceAttr(fqrn, "description", "Test description"),
					resource.TestCheckResourceAttr(fqrn, "release_notes.%", "2"),
					resource.TestCheckResourceAttr(fqrn, "release_notes.syntax", "plain_text"),
					resource.TestCheckResourceAttr(fqrn, "release_notes.content", "test release notes"),
					resource.TestCheckResourceAttr(fqrn, "spec.%", "1"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.#", "1"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.aql", "items.find({ \"repo\" : \"example-repo-local\" })"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.query_name", "query-1"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.mappings.#", "1"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.mappings.0.input", "original_repository/(.*)"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.mappings.0.output", "new_repository/$1"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.added_props.#", "1"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.added_props.0.key", "test-key"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.added_props.0.values.#", "1"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.added_props.0.values.0", "test-value"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.exclude_props_patterns.#", "1"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.exclude_props_patterns.0", "test-patterns"),
					resource.TestCheckResourceAttrSet(fqrn, "state"),
					resource.TestCheckResourceAttrSet(fqrn, "created"),
					resource.TestCheckResourceAttrSet(fqrn, "created_by"),
					resource.TestCheckResourceAttr(fqrn, "artifacts.#", "8"),
					resource.TestCheckResourceAttrSet(fqrn, "artifacts_size"),
					resource.TestCheckResourceAttrSet(fqrn, "archived"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fqrn, "name", testData["name"]),
					resource.TestCheckResourceAttr(fqrn, "version", testData["version"]),
					resource.TestCheckResourceAttr(fqrn, "sign_immediately", testData["sign_immediately"]),
					resource.TestCheckResourceAttr(fqrn, "description", "Test description"),
					resource.TestCheckResourceAttr(fqrn, "release_notes.%", "2"),
					resource.TestCheckResourceAttr(fqrn, "release_notes.syntax", "plain_text"),
					resource.TestCheckResourceAttr(fqrn, "release_notes.content", "test release notes"),
					resource.TestCheckResourceAttr(fqrn, "spec.%", "1"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.#", "1"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.aql", "items.find({ \"repo\" : \"example-repo-local\" })"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.query_name", "query-1"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.mappings.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(fqrn, "spec.queries.0.mappings.*", map[string]string{
						"input":  "original_repository/(.*)",
						"output": "new_repository/$1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(fqrn, "spec.queries.0.mappings.*", map[string]string{
						"input":  "(.*)/(.*)",
						"output": "$1/new_folder/$2",
					}),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.added_props.#", "2"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.added_props.0.key", "test-key"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.added_props.0.values.#", "1"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.added_props.0.values.0", "test-value"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.added_props.1.key", "test-key-2"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.added_props.1.values.#", "1"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.added_props.1.values.0", "test-value-2"),
					resource.TestCheckResourceAttr(fqrn, "spec.queries.0.exclude_props_patterns.#", "2"),
					resource.TestCheckTypeSetElemAttr(fqrn, "spec.queries.0.exclude_props_patterns.*", "test-patterns"),
					resource.TestCheckTypeSetElemAttr(fqrn, "spec.queries.0.exclude_props_patterns.*", "test-patterns-2"),
					resource.TestCheckResourceAttrSet(fqrn, "state"),
					resource.TestCheckResourceAttrSet(fqrn, "created"),
					resource.TestCheckResourceAttrSet(fqrn, "created_by"),
					resource.TestCheckResourceAttr(fqrn, "artifacts.#", "8"),
					resource.TestCheckResourceAttrSet(fqrn, "artifacts_size"),
					resource.TestCheckResourceAttrSet(fqrn, "archived"),
				),
			},
			{
				ResourceName:                         fqrn,
				ImportState:                          true,
				ImportStateId:                        fmt.Sprintf("%s:%s", testData["name"], testData["version"]),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"dry_run", "sign_immediately"},
			},
		},
	})
}

func TestAccReleaseBundleV1_invalid_name(t *testing.T) {
	testCases := []struct {
		name       string
		errorRegex string
	}{
		{name: "@invalid", errorRegex: `.*must begin with a letter or digit and consist only of letters,\n.*digits, underscores, periods, hyphens, and colons.*`},
		{name: "invalid@", errorRegex: `.*must begin with a letter or digit and consist only of letters,\n.*digits, underscores, periods, hyphens, and colons.*`},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, _, resourceName := testutil.MkNames("test-release-bundle-v1", "distribution_release_bundle_v1")

			const template = `
			resource "distribution_release_bundle_v1" "{{ .resource_name }}" {
				name = "{{ .name }}"
				version = "{{ .version }}"
				sign_immediately = {{ .sign_immediately }}
				description = "Test description"
		
				release_notes = {
					syntax = "plain_text"
					content = "test release notes"
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
							key = "test-key"
							values = ["test-value"]
						}]
						
						exclude_props_patterns = [
							"test-patterns"
						]
					}]
				}
			}`

			testData := map[string]string{
				"resource_name":    resourceName,
				"name":             testCase.name,
				"version":          "1.0.0",
				"sign_immediately": "false",
			}

			config := util.ExecuteTemplate("TestAccReleaseBundleV1_full", template, testData)

			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProviders(),
				Steps: []resource.TestStep{
					{
						Config:      config,
						ExpectError: regexp.MustCompile(testCase.errorRegex),
					},
				},
			})
		})
	}
}

func TestAccReleaseBundleV1_invalid_version(t *testing.T) {
	testCases := []struct {
		version    string
		errorRegex string
	}{
		{version: "@invalid", errorRegex: `.*must begin with a letter or digit and consist only of\n.*letters, digits, underscores, periods, hyphens, and colons.*`},
		{version: "invalid@", errorRegex: `.*must begin with a letter or digit and consist only of\n.*letters, digits, underscores, periods, hyphens, and colons.*`},
		{version: "LATEST", errorRegex: `.*value must be none of: \["LATEST"\].*`},
	}
	for _, testCase := range testCases {
		t.Run(testCase.version, func(t *testing.T) {
			_, _, resourceName := testutil.MkNames("test-release-bundle-v1", "distribution_release_bundle_v1")

			const template = `
			resource "distribution_release_bundle_v1" "{{ .resource_name }}" {
				name = "test-name"
				version = "{{ .version }}"
				sign_immediately = {{ .sign_immediately }}
				description = "Test description"
		
				release_notes = {
					syntax = "plain_text"
					content = "test release notes"
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
							key = "test-key"
							values = ["test-value"]
						}]
						
						exclude_props_patterns = [
							"test-patterns"
						]
					}]
				}
			}`

			testData := map[string]string{
				"resource_name":    resourceName,
				"version":          testCase.version,
				"sign_immediately": "false",
			}

			config := util.ExecuteTemplate("TestAccReleaseBundleV1_full", template, testData)

			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProviders(),
				Steps: []resource.TestStep{
					{
						Config:      config,
						ExpectError: regexp.MustCompile(testCase.errorRegex),
					},
				},
			})
		})
	}
}

func TestAccReleaseBundleV1_invalid_query_name(t *testing.T) {
	testCases := []struct {
		queryName  string
		errorRegex string
	}{
		{queryName: "1invalid", errorRegex: `.*must start with alphabetic character followed by an alphanumeric or '_-.:'.*`},
		{queryName: "invalid@", errorRegex: `.*must start with alphabetic character followed by an alphanumeric or '_-.:'.*`},
		{queryName: "i", errorRegex: `.*must start with alphabetic character followed by an alphanumeric or '_-.:'.*`},
	}
	for _, testCase := range testCases {
		t.Run(testCase.queryName, func(t *testing.T) {
			_, _, resourceName := testutil.MkNames("test-release-bundle-v1", "distribution_release_bundle_v1")

			const template = `
			resource "distribution_release_bundle_v1" "{{ .resource_name }}" {
				name = "test-name"
				version = "1.0.0"
				sign_immediately = {{ .sign_immediately }}
				description = "Test description"
		
				release_notes = {
					syntax = "plain_text"
					content = "test release notes"
				}
		
				spec = {
					queries = [{
						aql = "items.find({ \"repo\" : \"example-repo-local\" })"
						query_name: "{{ .query_name }}"
		
						mappings = [{
							input = "original_repository/(.*)"
							output = "new_repository/$1"
						}]
		
						added_props = [{
							key = "test-key"
							values = ["test-value"]
						}]
						
						exclude_props_patterns = [
							"test-patterns"
						]
					}]
				}
			}`

			testData := map[string]string{
				"resource_name":    resourceName,
				"query_name":       testCase.queryName,
				"sign_immediately": "false",
			}

			config := util.ExecuteTemplate("TestAccReleaseBundleV1_full", template, testData)

			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProviders(),
				Steps: []resource.TestStep{
					{
						Config:      config,
						ExpectError: regexp.MustCompile(testCase.errorRegex),
					},
				},
			})
		})
	}
}
