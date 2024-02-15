package functions_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccProviderFunction_resource_name_from_self_link(t *testing.T) {
	t.Parallel()
	acctest.SkipIfVcr(t) // Need to determine if compatible with VCR, as functions are implemented in PF provider

	resourceName := "tf-test-my-resource"
	resourceNameRegex := regexp.MustCompile(fmt.Sprintf("^%s$", resourceName))

	zoneName := "us-central1-c"
	zoneNameRegex := regexp.MustCompile(fmt.Sprintf("^%s$", zoneName))

	validSelfLink := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-c/instances/%s", resourceName)
	truncatedValidSelfLink := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/my-project/zones/%s", zoneName)
	validId := fmt.Sprintf("projects/my-project/zones/us-central1-c/instances/%s", resourceName)
	repetitiveInput := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-c/instances/not-this-1/instances/not-this-2/instances/%s", resourceName)
	invalidInput := resourceName

	context := map[string]interface{}{
		"function_name": "resource_name_from_self_link",
		"output_name":   "resource_name",
		"resource_name": resourceName,
		"self_link":     "", // overridden in test cases
	}

	acctest.VcrTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Given valid resource self_link input, the output value matches the expected value
				Config: testProviderFunction_generic_element_from_self_link(context, validSelfLink),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), resourceNameRegex),
				),
			},
			{
				// Given a truncated version of a resource's self_link input, the output will be the last element in the path
				// This test case is included to show that retrieving resource name is highly dependent on the input, and
				// we cannot accurately determing good versus bad input strings because we'd need information about the
				// user's usecase
				Config: testProviderFunction_generic_element_from_self_link(context, truncatedValidSelfLink),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), zoneNameRegex),
				),
			},
			{
				// Given valid resource id input, the output value matches the expected value
				Config: testProviderFunction_generic_element_from_self_link(context, validId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), resourceNameRegex),
				),
			},
			{
				// Given repetitive input, the output value is the element at the end of the string
				Config: testProviderFunction_generic_element_from_self_link(context, repetitiveInput),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), resourceNameRegex),
				),
			},
			{
				// Given invalid input, an error occurs
				// NOTE: for resource name, function will pull last value of any path. Invalid means empty string or strings without `/`.
				Config:      testProviderFunction_generic_element_from_self_link(context, invalidInput),
				ExpectError: regexp.MustCompile("Error in function call"), // ExpectError doesn't inspect the specific error messages
			},
			{
				// Can get the resource name from a resource's id in one step
				// Uses google_pubsub_topic resource's id attribute with format projects/{{project}}/topics/{{name}}
				Config: testProviderFunction_get_resource_name_from_resource_id(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), resourceNameRegex),
				),
			},
			{
				// Can get the resource name from a resource's self_link in one step
				// Uses google_compute_subnetwork resource's self_link attribute
				Config: testProviderFunction_get_resource_name_from_resource_self_link(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchOutput(context["output_name"].(string), resourceNameRegex),
				),
			},
		},
	})
}

func testProviderFunction_get_resource_name_from_resource_id(context map[string]interface{}) string {
	return acctest.Nprintf(`
# terraform block required for provider function to be found
terraform {
	required_providers {
		google = {
			source = "hashicorp/google"
		}
	}
}

resource "google_pubsub_topic" "example" {
  name = "%{resource_name}"
}

output "%{output_name}" {
	value = provider::google::%{function_name}(google_pubsub_topic.example.id)
}
`, context)
}

func testProviderFunction_get_resource_name_from_resource_self_link(context map[string]interface{}) string {
	return acctest.Nprintf(`
# terraform block required for provider function to be found
terraform {
	required_providers {
		google = {
			source = "hashicorp/google"
		}
	}
}

data "google_compute_network" "default" {
  name = "default"
}

resource "google_compute_subnetwork" "subnet" {
  name          = "%{resource_name}"
  ip_cidr_range = "10.2.0.0/16"
  network        = data.google_compute_network.default.id
}

output "%{output_name}" {
	value = provider::google::%{function_name}(google_compute_subnetwork.subnet.self_link)
}
`, context)
}
