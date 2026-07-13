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

// TestManagedRedis verifies managed redis is deployed and listed in Azure.
// Functional test runs both ID verification and list verification.
func TestManagedRedis(t *testing.T, ctx types.TestContext) {
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		t.Fatal("ARM_SUBSCRIPTION_ID environment variable is not set")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Unable to get credentials: %v\n", err)
	}

	t.Run("TestManagedRedisID", func(t *testing.T) {
		checkManagedRedisID(t, ctx, subscriptionID, cred)
	})

	// Managed Redis does not expose a safe idempotent write operation.
	// Functional coverage is strengthened by an additional list verification
	// that is not part of the readonly path.
	t.Run("TestManagedRedisListedInResourceGroup", func(t *testing.T) {
		checkManagedRedisListedInResourceGroup(t, ctx, subscriptionID, cred)
	})
}

// TestComposableManagedRedisReadOnly verifies managed redis ID using read operations only.
func TestComposableManagedRedisReadOnly(t *testing.T, ctx types.TestContext) {
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		t.Fatal("ARM_SUBSCRIPTION_ID environment variable is not set")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Unable to get credentials: %v\n", err)
	}

	t.Run("TestManagedRedisID", func(t *testing.T) {
		checkManagedRedisID(t, ctx, subscriptionID, cred)
	})
}

func checkManagedRedisID(t *testing.T, ctx types.TestContext, subscriptionID string, cred *azidentity.DefaultAzureCredential) {
	client := NewManagedRedisClient(t, subscriptionID, cred)

	resourceGroupName := terraform.Output(t, ctx.TerratestTerraformOptions(), "resource_group_name")
	managedRedisName := terraform.Output(t, ctx.TerratestTerraformOptions(), "managed_redis_name")
	expectedID := terraform.Output(t, ctx.TerratestTerraformOptions(), "managed_redis_id")
	expectedHostname := terraform.Output(t, ctx.TerratestTerraformOptions(), "managed_redis_hostname")
	expectedSKUName := terraform.Output(t, ctx.TerratestTerraformOptions(), "managed_redis_sku_name")

	managedRedis, err := client.Get(context.TODO(), resourceGroupName, managedRedisName, nil)
	if err != nil {
		t.Fatalf("failed to get Managed Redis instance: %v", err)
	}

	expectedIDLower := strings.ToLower(expectedID)
	actualIDLower := strings.ToLower(*managedRedis.ID)

	assert.Equal(t, expectedIDLower, actualIDLower, "Managed Redis ID doesn't match")

	if assert.NotNil(t, managedRedis.SKU, "Managed Redis SKU must be returned") && assert.NotNil(t, managedRedis.SKU.Name, "Managed Redis SKU name must be returned") {
		assert.Equal(t, expectedSKUName, string(*managedRedis.SKU.Name), "Managed Redis SKU doesn't match")
	}

	if assert.NotNil(t, managedRedis.Properties, "Managed Redis properties must be returned") && assert.NotNil(t, managedRedis.Properties.HostName, "Managed Redis hostname must be returned") {
		assert.Equal(t, strings.ToLower(expectedHostname), strings.ToLower(*managedRedis.Properties.HostName), "Managed Redis hostname doesn't match")
	}

	if assert.NotNil(t, managedRedis.Tags, "Managed Redis tags must be returned") {
		provisionerTag, hasProvisionerTag := managedRedis.Tags["provisioner"]
		if assert.True(t, hasProvisionerTag, "Managed Redis tags must include provisioner") && assert.NotNil(t, provisionerTag, "Managed Redis provisioner tag must have a value") {
			assert.Equal(t, "terraform", strings.ToLower(*provisionerTag), "Managed Redis provisioner tag doesn't match")
		}

		resourceNameTag, hasResourceNameTag := managedRedis.Tags["resource_name"]
		if assert.True(t, hasResourceNameTag, "Managed Redis tags must include resource_name") && assert.NotNil(t, resourceNameTag, "Managed Redis resource_name tag must have a value") {
			assert.Equal(t, managedRedisName, *resourceNameTag, "Managed Redis resource_name tag doesn't match")
		}
	}
}

func checkManagedRedisListedInResourceGroup(t *testing.T, ctx types.TestContext, subscriptionID string, cred *azidentity.DefaultAzureCredential) {
	client := NewManagedRedisClient(t, subscriptionID, cred)

	resourceGroupName := terraform.Output(t, ctx.TerratestTerraformOptions(), "resource_group_name")
	managedRedisName := terraform.Output(t, ctx.TerratestTerraformOptions(), "managed_redis_name")

	pager := client.NewListByResourceGroupPager(resourceGroupName, nil)
	foundManagedRedis := false

	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			t.Fatalf("failed to list Managed Redis instances in resource group: %v", err)
		}

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

	assert.True(t, foundManagedRedis, "Managed Redis instance was not found in resource group list")
}

func NewManagedRedisClient(t *testing.T, subscriptionID string, cred *azidentity.DefaultAzureCredential) *armredisenterprise.Client {
	client, err := armredisenterprise.NewClient(subscriptionID, cred, nil)
	if err != nil {
		t.Fatalf("failed to create Managed Redis client: %v", err)
	}
	return client
}
