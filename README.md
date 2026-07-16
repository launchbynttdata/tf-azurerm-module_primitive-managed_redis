# tf-azurerm-module_primitive-managed_redis

Terraform primitive module for Azure Managed Redis using `azurerm` provider v4.

## Overview

This module provisions a single `azurerm_managed_redis` resource and exposes core connection and database outputs.

## Requirements

- Terraform `>= 1.0`
- Provider `hashicorp/azurerm` `>= 4.0, < 5.0`

## Features

- Managed Redis deployment
- Optional managed identity configuration
- Optional customer-managed key configuration
- Configurable default database settings
- Standardized tagging (`provisioner`, `resource_name`)

## Usage

```hcl
module "managed_redis" {
  source = "terraform.registry.launch.nttdata.com/module_primitive/managed_redis/azurerm"

  name                = "example-managed-redis"
  location            = "eastus"
  resource_group_name = "example-rg"

  sku_name                  = "Balanced_B1"
  high_availability_enabled = true
  public_network_access     = "Disabled"

  default_database = {
    access_keys_authentication_enabled = true
    client_protocol                    = "Encrypted"
    clustering_policy                  = "OSSCluster"
    eviction_policy                    = "VolatileLRU"
  }

  tags = {
    environment = "dev"
  }
}
```

## Outputs

- `id`
- `name`
- `hostname`
- `database_id`
- `database_port`
- `primary_access_key` (sensitive)
- `secondary_access_key` (sensitive)

## Example

See [examples/azure_managed](examples/azure_managed) for the baseline runnable example.
See [examples/azure_managed_default_database](examples/azure_managed_default_database) for a variant that explicitly sets non-default `default_database` values to exercise dynamic configuration paths.

## Test Dependency Baseline

The current `go.mod` test dependencies are intentionally pinned to align with the validated LCAF terratest baseline used by this repository:

- `github.com/gruntwork-io/terratest v0.43.12`
- `github.com/launchbynttdata/lcaf-component-terratest v1.0.3`
- `github.com/stretchr/testify v1.9.0`

These pins are deliberate for compatibility with the existing CI/runtime expectations and should be upgraded together as a coordinated change.

## Azure Feature Registration Prerequisite

This module requires the **`Microsoft.Cache/AmrAugust2025Preview`** preview feature. The test framework handles registration automatically:

### Test Behavior (Same for Local and CI)

1. **Check feature state** → If already `Registered`, proceed
2. **Auto-register** → If `NotRegistered`, register it automatically
3. **Wait for propagation** → Wait up to 2 minutes for `Registered` state
   - If reached → Wait additional 60 seconds for Azure internal propagation, then run test
   - If timeout → **Skip test** (don't fail) to avoid blocking pipelines during initial subscription setup
4. **Run provisioning + assertions** → Full test executes once feature is ready

### Expected Test Runtime

- **Feature already registered**: ~30 minutes (full provisioning + assertions)
- **Feature auto-registered first time**: ~2–15 minutes (register + wait for propagation + skip or test)
- **Subsequent runs**: ~30 minutes (feature already ready)

### Command to Pre-Register (Optional)

To avoid initial registration wait, pre-register the feature at the subscription level:

```bash
az feature register --namespace Microsoft.Cache --name AmrAugust2025Preview
az provider register -n Microsoft.Cache
# Wait 15–30 minutes for state to reach "Registered"
```

Once the feature is `Registered` at the subscription level, all test runs will proceed directly to provisioning (~30m).

## Module Development

### Pre-Requisites

The following commands should be available on your system:

- `asdf` or `mise`
- `make`
- `python3` (for pre-commit)

Additionally, your `git` user and email must be configured. Run the `make configure` command from the root of the repository to ensure that you meet these requirements.

### Pre-Commit Hooks

The [.pre-commit-config.yaml](.pre-commit-config.yaml) file defines certain `pre-commit` hooks that are relevant to Terraform and Golang, as well as common linting tasks. These will be configured for you when you run `make configure`.

### Local Validation

You should validate module changes locally before pushing a branch.

1. Ensure that you have run `make configure` successfully.
2. Ensure you are signed into the appropriate cloud provider for the module under test.
3. Run linting:

```bash
make lint
```

4. Run integration tests (provisions infrastructure, runs tests, and tears down):

```bash
make test
```

### Review & Merge Process

Open a pull request targeting `main` after local validation passes.

The pull request title determines the version bump and must follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#specification).

- breaking change -> major
- feature -> minor
- other changes -> patch

Ensure CI workflows are passing, address review feedback, and obtain required approvals from CODEOWNERS before merge.

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.0 |
| <a name="requirement_azurerm"></a> [azurerm](#requirement\_azurerm) | >= 4.0, < 5.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [azurerm_managed_redis.redis](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/managed_redis) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_customer_managed_key"></a> [customer\_managed\_key](#input\_customer\_managed\_key) | Optional customer-managed key configuration for encryption. | <pre>object({<br/>    key_vault_key_id          = string<br/>    user_assigned_identity_id = string<br/>  })</pre> | `null` | no |
| <a name="input_default_database"></a> [default\_database](#input\_default\_database) | Default database configuration for Managed Redis. | <pre>object({<br/>    access_keys_authentication_enabled            = optional(bool, false)<br/>    client_protocol                               = optional(string, "Encrypted")<br/>    clustering_policy                             = optional(string, "OSSCluster")<br/>    eviction_policy                               = optional(string, "VolatileLRU")<br/>    geo_replication_group_name                    = optional(string)<br/>    persistence_append_only_file_backup_frequency = optional(string)<br/>    persistence_redis_database_backup_frequency   = optional(string)<br/>    modules = optional(list(object({<br/>      name = string<br/>      args = optional(string)<br/>    })), [])<br/>  })</pre> | `{}` | no |
| <a name="input_high_availability_enabled"></a> [high\_availability\_enabled](#input\_high\_availability\_enabled) | Whether high availability is enabled for Managed Redis. | `bool` | `true` | no |
| <a name="input_identity"></a> [identity](#input\_identity) | Managed identity block. Allowed type values: SystemAssigned, UserAssigned, or SystemAssigned, UserAssigned. | <pre>object({<br/>    type         = string<br/>    identity_ids = optional(list(string))<br/>  })</pre> | `null` | no |
| <a name="input_location"></a> [location](#input\_location) | Azure region where the Managed Redis instance is deployed. | `string` | n/a | yes |
| <a name="input_name"></a> [name](#input\_name) | Name of the Managed Redis instance. | `string` | n/a | yes |
| <a name="input_public_network_access"></a> [public\_network\_access](#input\_public\_network\_access) | Public network access mode. Allowed values are Enabled and Disabled. | `string` | `"Disabled"` | no |
| <a name="input_resource_group_name"></a> [resource\_group\_name](#input\_resource\_group\_name) | Name of the resource group containing the Managed Redis instance. | `string` | n/a | yes |
| <a name="input_sku_name"></a> [sku\_name](#input\_sku\_name) | Managed Redis SKU name. Example values: Balanced\_B1, Balanced\_B3, ComputeOptimized\_X3, MemoryOptimized\_M10. | `string` | `"Balanced_B1"` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | Tags to apply to the Managed Redis instance. | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_database_id"></a> [database\_id](#output\_database\_id) | The ID of the default Managed Redis database. |
| <a name="output_database_port"></a> [database\_port](#output\_database\_port) | The port of the default Managed Redis database. |
| <a name="output_hostname"></a> [hostname](#output\_hostname) | Hostname of the Managed Redis endpoint. |
| <a name="output_id"></a> [id](#output\_id) | The ID of the Managed Redis instance. |
| <a name="output_name"></a> [name](#output\_name) | The name of the Managed Redis instance. |
| <a name="output_primary_access_key"></a> [primary\_access\_key](#output\_primary\_access\_key) | Primary access key for the default database when access key authentication is enabled. |
| <a name="output_secondary_access_key"></a> [secondary\_access\_key](#output\_secondary\_access\_key) | Secondary access key for the default database when access key authentication is enabled. |
<!-- END_TF_DOCS -->
