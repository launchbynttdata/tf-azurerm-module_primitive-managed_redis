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

	t.Run("TestManagedRedisExists", func(t *testing.T) {
		checkManagedRedisExists(t, ctx, subscriptionID, cred)
	})
}

// TestManagedRedisReadOnly verifies managed redis exists using read operations only.
func TestManagedRedisReadOnly(t *testing.T, ctx types.TestContext) {
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		t.Fatal("ARM_SUBSCRIPTION_ID environment variable is not set")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Unable to get credentials: %v\n", err)
	}

	t.Run("TestManagedRedisExists", func(t *testing.T) {
		checkManagedRedisExists(t, ctx, subscriptionID, cred)
	})
}

func checkManagedRedisExists(t *testing.T, ctx types.TestContext, subscriptionID string, cred *azidentity.DefaultAzureCredential) {
	client := NewManagedRedisClient(t, subscriptionID, cred)

	resourceGroupName := terraform.Output(t, ctx.TerratestTerraformOptions(), "resource_group_name")
	managedRedisName := terraform.Output(t, ctx.TerratestTerraformOptions(), "managed_redis_name")
	expectedID := terraform.Output(t, ctx.TerratestTerraformOptions(), "managed_redis_id")

	managedRedis, err := client.Get(context.Background(), resourceGroupName, managedRedisName, nil)
	if err != nil {
		t.Fatalf("failed to get Managed Redis instance: %v", err)
	}

	expectedIDLower := strings.ToLower(expectedID)
	actualIDLower := strings.ToLower(*managedRedis.ID)

	assert.Equal(t, expectedIDLower, actualIDLower, "Managed Redis ID doesn't match")
	assert.NotEmpty(t, managedRedis.Name, "Managed Redis name should not be empty")
	assert.NotEmpty(t, managedRedis.Properties.HostName, "Managed Redis hostname should not be empty")
}

func NewManagedRedisClient(t *testing.T, subscriptionID string, cred *azidentity.DefaultAzureCredential) *armredisenterprise.Client {
	client, err := armredisenterprise.NewClient(subscriptionID, cred, nil)
	if err != nil {
		t.Fatalf("failed to create Managed Redis client: %v", err)
	}
	return client
}
