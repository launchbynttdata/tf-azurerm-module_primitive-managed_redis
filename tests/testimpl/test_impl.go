package testimpl

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestManagedRedis verifies managed redis is deployed correctly and exists in Azure.
func TestManagedRedis(t *testing.T, ctx types.TestContext) {
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		t.Fatal("ARM_SUBSCRIPTION_ID environment variable is not set")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Unable to get credentials: %v\n", err)
	}

	t.Run("TestManagedRedisMatchesTerraformOutputs", func(t *testing.T) {
		checkManagedRedisMatchesTerraformOutputs(t, ctx, subscriptionID, cred)
	})

	// Managed Redis does not expose a safe, idempotent write operation for functional tests
	// in this module context. Functional coverage is strengthened with an additional list
	// verification that is not part of the readonly path.
	t.Run("TestManagedRedisListedInResourceGroup", func(t *testing.T) {
		checkManagedRedisListedInResourceGroup(t, ctx, subscriptionID, cred)
	})
}

// TestComposableManagedRedisReadOnly verifies managed redis exists using read operations only.
func TestComposableManagedRedisReadOnly(t *testing.T, ctx types.TestContext) {
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		t.Fatal("ARM_SUBSCRIPTION_ID environment variable is not set")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Unable to get credentials: %v\n", err)
	}

	t.Run("TestManagedRedisMatchesTerraformOutputs", func(t *testing.T) {
		checkManagedRedisMatchesTerraformOutputs(t, ctx, subscriptionID, cred)
	})
}

func checkManagedRedisMatchesTerraformOutputs(t *testing.T, ctx types.TestContext, subscriptionID string, cred *azidentity.DefaultAzureCredential) {
	client := NewManagedRedisClient(t, subscriptionID, cred)

	resourceGroupName := terraform.Output(t, ctx.TerratestTerraformOptions(), "resource_group_name")
	managedRedisName := terraform.Output(t, ctx.TerratestTerraformOptions(), "managed_redis_name")
	expectedID := terraform.Output(t, ctx.TerratestTerraformOptions(), "managed_redis_id")
	expectedHostName := terraform.Output(t, ctx.TerratestTerraformOptions(), "managed_redis_hostname")

	managedRedis, err := client.Get(context.Background(), resourceGroupName, managedRedisName, nil)
	require.NoError(t, err, "failed to get Managed Redis instance")

	require.NotNil(t, managedRedis.ID, "Managed Redis ID should be set")
	require.NotNil(t, managedRedis.Name, "Managed Redis name should be set")
	require.NotNil(t, managedRedis.Properties, "Managed Redis properties should be set")
	require.NotNil(t, managedRedis.Properties.HostName, "Managed Redis hostname should be set")

	expectedIDLower := strings.ToLower(expectedID)
	actualIDLower := strings.ToLower(*managedRedis.ID)

	assert.Equal(t, expectedIDLower, actualIDLower, "Managed Redis ID doesn't match")
	assert.Equal(t, managedRedisName, *managedRedis.Name, "Managed Redis name doesn't match")
	assert.Equal(t, expectedHostName, *managedRedis.Properties.HostName, "Managed Redis hostname doesn't match")
}

func checkManagedRedisListedInResourceGroup(t *testing.T, ctx types.TestContext, subscriptionID string, cred *azidentity.DefaultAzureCredential) {
	client := NewManagedRedisClient(t, subscriptionID, cred)

	resourceGroupName := terraform.Output(t, ctx.TerratestTerraformOptions(), "resource_group_name")
	managedRedisName := terraform.Output(t, ctx.TerratestTerraformOptions(), "managed_redis_name")

	pager := client.NewListByResourceGroupPager(resourceGroupName, nil)
	foundManagedRedis := false

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err, "failed to list Managed Redis instances in resource group")

		for _, managedRedis := range page.Value {
			if managedRedis.Name != nil && strings.EqualFold(*managedRedis.Name, managedRedisName) {
				foundManagedRedis = true
				break
			}
		}

		if foundManagedRedis {
			break
		}
	}

	assert.True(t, foundManagedRedis, "Managed Redis instance was not found in ListByResourceGroup results")
}

func NewManagedRedisClient(t *testing.T, subscriptionID string, cred *azidentity.DefaultAzureCredential) *armredisenterprise.Client {
	client, err := armredisenterprise.NewClient(subscriptionID, cred, nil)
	if err != nil {
		t.Fatalf("failed to create Managed Redis client: %v", err)
	}
	return client
}
