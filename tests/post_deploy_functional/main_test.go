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
	"time"

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

// checkAMRFeatureRegistered verifies that AmrAugust2025Preview reaches Registered state.
//
// Auto-registers the feature if NotRegistered, then waits (polling every 30s) until
// the feature is Registered. Once Registered, waits 60s for Azure internal propagation
// before allowing the test to proceed. This ensures the feature is fully ready across
// all Azure endpoints before attempting to provision the managed_redis resource.
func checkAMRFeatureRegistered(t *testing.T) {
	t.Helper()

	const (
		featureNamespace = "Microsoft.Cache"
		featureName      = "AmrAugust2025Preview"
		pollInterval     = 30 * time.Second
	)

	getFeatureState := func() string {
		cmd := exec.Command("az", "feature", "show", "--namespace", featureNamespace,
			"--name", featureName, "--query", "properties.state", "-o", "tsv")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Could not determine state of %s/%s (is 'az' installed and authenticated?): %v",
				featureNamespace, featureName, err)
		}
		return strings.TrimSpace(string(output))
	}

	state := getFeatureState()

	// Auto-register if NotRegistered
	if state == "NotRegistered" {
		t.Logf("Feature %s/%s is NotRegistered; registering...", featureNamespace, featureName)
		regCmd := exec.Command("az", "feature", "register", "--namespace", featureNamespace, "--name", featureName)
		if out, err := regCmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed to register feature %s/%s: %v\n%s", featureNamespace, featureName, err, out)
		}
		provCmd := exec.Command("az", "provider", "register", "-n", featureNamespace)
		if out, err := provCmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed to re-register provider %s: %v\n%s", featureNamespace, err, out)
		}
		time.Sleep(5 * time.Second)
	}

	// Wait indefinitely for Registered state (poll every 30s)
	for {
		state = getFeatureState()
		if state == "Registered" {
			t.Logf("Feature %s/%s is Registered; waiting 60s for propagation...", featureNamespace, featureName)
			time.Sleep(60 * time.Second)
			t.Logf("Feature %s/%s propagation complete; proceeding with test", featureNamespace, featureName)
			return
		}
		t.Logf("Feature %s/%s is %q; checking again in %v...", featureNamespace, featureName, state, pollInterval)
		time.Sleep(pollInterval)
	}
}
}

func TestManagedRedisModule(t *testing.T) {
	setTerraformInitArgsForTests()
	checkAMRFeatureRegistered(t)

	ctx := types.CreateTestContextBuilder().
		SetTestConfig(&testimpl.ThisTFModuleConfig{}).
		SetTestConfigFolderName(testConfigsExamplesFolderDefault).
		SetTestConfigFileName(infraTFVarFileNameDefault).
		Build()

	lib.RunSetupTestTeardown(t, *ctx, testimpl.TestManagedRedis)
}
