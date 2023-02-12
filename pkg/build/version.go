package build

import (
	"fmt"
	"runtime"
)

var buildVersion = "dev"
var buildCommit string
var buildDate string
var builtBy string

// Version returns build information added during the build process.
func Version() string {
	result := buildVersion
	if buildCommit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, buildCommit)
	}
	if buildDate != "" {
		result = fmt.Sprintf("%s\nbuild date: %s", result, buildDate)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	result = fmt.Sprintf("%s\ngoos: %s\ngoarch: %s", result, runtime.GOOS, runtime.GOARCH)

	return result
}
