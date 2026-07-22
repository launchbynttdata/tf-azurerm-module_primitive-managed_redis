// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/launchbynttdata/lcaf-component-terratest/lib"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/launchbynttdata/tf-azurerm-module_primitive-managed_redis/tests/testimpl"
)

const (
	testConfigsExamplesFolderDefault = "../../examples"
	infraTFVarFileNameDefault        = "test.tfvars"
)

func setTerraformInitArgsForTests() {
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

func requireAMRFeatureRegistered(t *testing.T) {
	t.Helper()

	cmd := exec.Command(
		"az",
		"feature",
		"show",
		"--namespace", "Microsoft.Cache",
		"--name", "AmrAugust2025Preview",
		"--query", "properties.state",
		"-o", "tsv",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf(
			"Could not determine feature state for Microsoft.Cache/AmrAugust2025Preview: %v. Ensure az is installed and authenticated. Output: %s",
			err,
			strings.TrimSpace(string(out)),
		)
	}

	state := strings.TrimSpace(string(out))
	if state != "Registered" {
		t.Fatalf(
			"Prerequisite not met: Microsoft.Cache/AmrAugust2025Preview must be Registered before running functional tests (current state: %q). Run: az feature register --namespace Microsoft.Cache --name AmrAugust2025Preview && az provider register -n Microsoft.Cache",
			state,
		)
	}
}

func TestManagedRedisModule(t *testing.T) {
	setTerraformInitArgsForTests()
	requireAMRFeatureRegistered(t)

	ctx := types.CreateTestContextBuilder().
		SetTestConfig(&testimpl.ThisTFModuleConfig{}).
		SetTestConfigFolderName(testConfigsExamplesFolderDefault).
		SetTestConfigFileName(infraTFVarFileNameDefault).
		Build()

	lib.RunSetupTestTeardown(t, *ctx, testimpl.TestComposableManagedRedis)
}
