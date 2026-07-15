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

// ensureAMRFeatureRegistered ensures AmrAugust2025Preview feature is Registered.
// If NotRegistered/Pending, it registers and waits briefly. If still not Registered after
// a short wait, test is skipped (for first CI pipeline run when feature is propagating).
func ensureAMRFeatureRegistered(t *testing.T) {
	const (
		featureNamespace = "Microsoft.Cache"
		featureName      = "AmrAugust2025Preview"
		maxWait          = 2 * time.Minute  // Short wait for local re-runs
		pollInterval     = 30 * time.Second
	)

	// Helper to get current feature state
	getFeatureState := func() (string, error) {
		cmd := exec.Command("az", "feature", "show", "--namespace", featureNamespace,
			"--name", featureName, "--query", "properties.state", "-o", "tsv")
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(output)), nil
	}

	// Check current state
	state, err := getFeatureState()
	if err != nil {
		t.Logf("Warning: Could not check %s feature state: %v", featureName, err)
		return
	}

	// If already Registered, proceed
	if state == "Registered" {
		t.Logf("Feature %s is Registered; proceeding with test", featureName)
		return
	}

	// If NotRegistered, register it
	if state == "NotRegistered" {
		t.Logf("Feature %s is NotRegistered; registering...", featureName)
		cmd := exec.Command("az", "feature", "register", "--namespace", featureNamespace, "--name", featureName)
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Logf("Warning: Feature registration attempt returned: %v", err)
		}
		// Re-register provider to propagate
		cmd = exec.Command("az", "provider", "register", "-n", featureNamespace)
		if _, err := cmd.CombinedOutput(); err != nil {
			t.Logf("Warning: Provider re-registration returned: %v", err)
		}
		time.Sleep(5 * time.Second)
	}

	// Poll briefly for Registered state (don't wait 15 min on first run)
	startTime := time.Now()
	for {
		state, err := getFeatureState()
		if err != nil {
			t.Logf("Poll error: %v; retrying...", err)
		} else if state == "Registered" {
			t.Logf("Feature %s is now Registered; proceeding with test", featureName)
			return
		} else {
			t.Logf("Feature %s is %s; waiting for propagation...", featureName, state)
		}

		if time.Since(startTime) > maxWait {
			// Skip gracefully on first run when feature is still Pending/propagating
			t.Skipf("Feature %s is %s (still propagating after %v). Skipping provisioning test. "+
				"This is expected on first CI pipeline run. Rerun after feature propagates.",
				featureName, state, maxWait)
		}

		time.Sleep(pollInterval)
	}
}

func TestManagedRedisModule(t *testing.T) {
	setTerraformInitArgsForTests()
	ensureAMRFeatureRegistered(t)

	ctx := types.CreateTestContextBuilder().
		SetTestConfig(&testimpl.ThisTFModuleConfig{}).
		SetTestConfigFolderName(testConfigsExamplesFolderDefault).
		SetTestConfigFileName(infraTFVarFileNameDefault).
		Build()

	lib.RunSetupTestTeardown(t, *ctx, testimpl.TestManagedRedis)
}
