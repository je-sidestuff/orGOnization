package integration

import (
	"fmt"
	"os"
	"testing"
)

// SkipStageEnvVarPrefix is the prefix used for skipping stage environment variables, as inspired by terratest.
const SkipStageEnvVarPrefix = "SKIP_"

// RunTestStage is inspired by the `RunTestStage` function in terratest. It executes the given test stage (e.g., setup,
// teardown, validation) if an environment variable of the name `SKIP_<stageName>` (e.g., SKIP_teardown) is not set. If the
// stage is skipped the 'otherwise' func is run instead.
func RunTestStage(t *testing.T, stageName string, stage func(), otherwise func()) {
	envVarName := fmt.Sprintf("%s%s", SkipStageEnvVarPrefix, stageName)
	if os.Getenv(envVarName) == "" {
		t.Logf("The '%s' environment variable is not set, so executing stage '%s'.", envVarName, stageName)
		stage()
	} else {
		t.Logf("The '%s' environment variable is set, so skipping stage '%s'.", envVarName, stageName)
		otherwise()
	}
}
