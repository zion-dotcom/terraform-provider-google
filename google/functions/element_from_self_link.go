package functions

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// TODO name stuff well
func GetElement(ctx context.Context, input string, regex *regexp.Regexp, template string, pattern string, req function.RunRequest, resp *function.RunResponse) string {
	submatches := regex.FindAllStringSubmatchIndex(input, -1)

	// Zero matches means unusable input; error returned
	if len(submatches) == 0 {
		resp.Diagnostics.AddArgumentError(
			0,
			"No project id is present in the input string",
			fmt.Sprintf("The input string \"%s\" doesn't contain the expected pattern \"%s\".", input, pattern),
		)
		resp.Diagnostics.Append(resp.Result.Set(ctx, "")...)
		return ""
	}

	// >1 matches means input usable but not ideal; issue warning
	if len(submatches) > 1 {
		resp.Diagnostics.AddArgumentWarning(
			0,
			"Ambiguous input string could contain more than one project id",
			fmt.Sprintf("The input string \"%s\" contains more than one match for the pattern \"%s\". Terraform will use the first found match.", input, pattern),
		)
	}

	// Return found element
	submatch := submatches[0] // Take the only / left-most submatch
	result := []byte{}
	result = regex.ExpandString(result, template, input, submatch)
	return string(result)
}
