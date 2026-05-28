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

See [examples/with_cake](examples/with_cake) for a complete runnable example that includes resource name generation and resource group creation.

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
