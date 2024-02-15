package functions

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = EchoFunction{}

func NewResourceNameFromSelfLinkFunction() function.Function {
	return &ResourceNameFromSelfLinkFunction{}
}

type ResourceNameFromSelfLinkFunction struct{}

func (f ResourceNameFromSelfLinkFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "resource_name_from_self_link"
}

func (f ResourceNameFromSelfLinkFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Returns the resource name within the resource self link or id provided as an argument.",
		Description: "Takes a single string argument, which should be a self link or id of a resource. This function will either return the resource's short name from the input string or raise an error. The function returns the last element in that path before the end of the input string, e.g. when the function is passed the self link \"https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-c/instances/my-instance\" as an argument it will return \"my-instance\".",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "self_link",
				Description: "A self link of a resouce, or an id. For example, both \"https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-c/instances/my-instance\" and \"projects/my-project/zones/us-central1-c/instances/my-instance\" are valid inputs",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f ResourceNameFromSelfLinkFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {

	// Load arguments from function call
	var arg0 string
	resp.Diagnostics.Append(req.Arguments.GetArgument(ctx, 0, &arg0)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare how we'll identify resource name from input string
	regex := regexp.MustCompile("/(?P<ResourceName>[^/]+)$") // Should match the pattern below
	template := "$ResourceName"                              // Should match the submatch identifier in the regex
	pattern := "resourceType/{name}$"                        // Human-readable pseudo-regex pattern used in errors and warnings

	// Get and return element from input string
	resourceName := GetElement(ctx, arg0, regex, template, pattern, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.Result.Set(ctx, resourceName)...)
}
