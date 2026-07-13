package testimpl

import (
	"os"
	"strings"
)

// SetTerraformInitArgsForTests ensures terraform init honors the lock file in tests.
func SetTerraformInitArgsForTests() {
	const lockfileReadonlyArg = "-lockfile=readonly"

	current, exists := os.LookupEnv("TF_CLI_ARGS_init")
	if !exists {
		_ = os.Setenv("TF_CLI_ARGS_init", lockfileReadonlyArg)
		return
	}

	if strings.Contains(current, "-lockfile=") {
		return
	}

	_ = os.Setenv("TF_CLI_ARGS_init", strings.TrimSpace(current+" "+lockfileReadonlyArg))
}
