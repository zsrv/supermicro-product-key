package build

import (
	"fmt"
	"runtime"
)

var buildVersion = "1.1.0"
var buildCommit string
var buildDate = "unknown"
var builtBy string

// Version returns build information added during the build process.
func Version() string {
	result := buildVersion
	if buildCommit != "" {
		result = fmt.Sprintf("%s-%s", result, buildCommit)
	}

	result = fmt.Sprintf("%s %s/%s", result, runtime.GOOS, runtime.GOARCH)

	if buildDate != "" {
		result = fmt.Sprintf("%s BuildDate=%s", result, buildDate)
	}

	if builtBy != "" {
		result = fmt.Sprintf("%s BuiltBy=%s", result, builtBy)
	}

	return result
}
