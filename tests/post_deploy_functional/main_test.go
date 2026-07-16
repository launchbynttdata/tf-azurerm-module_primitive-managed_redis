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

// checkAMRFeatureRegistered verifies that AmrAugust2025Preview is in Registered state.
//
// In CI (LOCAL_RUN not set): any state other than Registered is a hard failure — the
// subscription must be bootstrapped with the preview feature before running tests.
//
// In local runs (LOCAL_RUN=true): if the feature is NotRegistered it is registered
// automatically, and if it is still Pending after a short wait the test is skipped
// rather than failing, so developers are not blocked during initial subscription setup.
func checkAMRFeatureRegistered(t *testing.T) {
	t.Helper()

	const (
		featureNamespace = "Microsoft.Cache"
		featureName      = "AmrAugust2025Preview"
		localMaxWait     = 2 * time.Minute
		pollInterval     = 30 * time.Second
	)

	isLocalRun := os.Getenv("LOCAL_RUN") == "true"

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

	if state == "Registered" {
		t.Logf("Feature %s/%s is Registered; waiting 60s for propagation...", featureNamespace, featureName)
		time.Sleep(60 * time.Second)
		t.Logf("Feature %s/%s propagation complete; proceeding", featureNamespace, featureName)
		return
	}

	if !isLocalRun {
		t.Fatalf("Feature %s/%s is %q — subscription prerequisite not met.\n"+
			"Register the preview feature at the subscription level before running CI:\n"+
			"  az feature register --namespace %s --name %s\n"+
			"  az provider register -n %s\n"+
			"Then wait for state to reach 'Registered' before re-running.",
			featureNamespace, featureName, state,
			featureNamespace, featureName, featureNamespace)
	}

	// LOCAL_RUN=true: attempt registration and wait briefly before skipping.
	if state == "NotRegistered" {
		t.Logf("LOCAL_RUN: Feature %s/%s is NotRegistered; registering...", featureNamespace, featureName)
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

	startTime := time.Now()
	for {
		state = getFeatureState()
		if state == "Registered" {
			t.Logf("LOCAL_RUN: Feature %s/%s is now Registered; waiting 60s for propagation...", featureNamespace, featureName)
			time.Sleep(60 * time.Second)
			t.Logf("LOCAL_RUN: Feature %s/%s propagation complete; proceeding", featureNamespace, featureName)
			return
		}
		t.Logf("LOCAL_RUN: Feature %s/%s is %q; waiting...", featureNamespace, featureName, state)
		if time.Since(startTime) >= localMaxWait {
			t.Skipf("LOCAL_RUN: Feature %s/%s is still %q after %v — skipping. "+
				"Rerun once the feature reaches Registered state.",
				featureNamespace, featureName, state, localMaxWait)
		}
		time.Sleep(pollInterval)
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
