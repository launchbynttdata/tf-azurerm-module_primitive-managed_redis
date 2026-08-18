# managed_redis default_database example

This example is equivalent to `examples/azure_managed`, but it explicitly sets a non-default `default_database` value to exercise dynamic block behavior in integration tests.

## Notes

- Uses `public_network_access = "Disabled"` by default.
- Sets `default_database.access_keys_authentication_enabled = true`.

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.5.0 |
| <a name="requirement_azurerm"></a> [azurerm](#requirement\_azurerm) | >= 4.0, < 5.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_managed_redis"></a> [managed\_redis](#module\_managed\_redis) | ../.. | n/a |
| <a name="module_resource_group"></a> [resource\_group](#module\_resource\_group) | terraform.registry.launch.nttdata.com/module_primitive/resource_group/azurerm | ~> 1.1 |
| <a name="module_resource_names"></a> [resource\_names](#module\_resource\_names) | terraform.registry.launch.nttdata.com/module_library/resource_name/launch | ~> 2.4 |

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_class_env"></a> [class\_env](#input\_class\_env) | Deployment environment used for name generation. | `string` | `"dev"` | no |
| <a name="input_customer_managed_key"></a> [customer\_managed\_key](#input\_customer\_managed\_key) | Optional CMK encryption settings. | <pre>object({<br/>    key_vault_key_id          = string<br/>    user_assigned_identity_id = string<br/>  })</pre> | `null` | no |
| <a name="input_default_database"></a> [default\_database](#input\_default\_database) | Default database configuration for managed redis. | <pre>object({<br/>    access_keys_authentication_enabled            = optional(bool, false)<br/>    client_protocol                               = optional(string, "Encrypted")<br/>    clustering_policy                             = optional(string, "OSSCluster")<br/>    eviction_policy                               = optional(string, "VolatileLRU")<br/>    geo_replication_group_name                    = optional(string)<br/>    persistence_append_only_file_backup_frequency = optional(string)<br/>    persistence_redis_database_backup_frequency   = optional(string)<br/>    modules = optional(list(object({<br/>      name = string<br/>      args = optional(string)<br/>    })), [])<br/>  })</pre> | `{}` | no |
| <a name="input_high_availability_enabled"></a> [high\_availability\_enabled](#input\_high\_availability\_enabled) | Whether high availability is enabled for the managed redis instance. | `bool` | `true` | no |
| <a name="input_identity"></a> [identity](#input\_identity) | Optional managed identity settings. | <pre>object({<br/>    type         = string<br/>    identity_ids = optional(list(string))<br/>  })</pre> | `null` | no |
| <a name="input_instance_env"></a> [instance\_env](#input\_instance\_env) | Environment instance number used for name generation. | `number` | `0` | no |
| <a name="input_instance_resource"></a> [instance\_resource](#input\_instance\_resource) | Resource instance number used for name generation. | `number` | `0` | no |
| <a name="input_location"></a> [location](#input\_location) | Azure region for the example deployment. | `string` | `"eastus"` | no |
| <a name="input_logical_product_family"></a> [logical\_product\_family](#input\_logical\_product\_family) | Logical product family used for name generation. | `string` | `"launch"` | no |
| <a name="input_logical_product_service"></a> [logical\_product\_service](#input\_logical\_product\_service) | Logical product service used for name generation. | `string` | `"redis"` | no |
| <a name="input_public_network_access"></a> [public\_network\_access](#input\_public\_network\_access) | Public network access mode for Managed Redis. | `string` | `"Disabled"` | no |
| <a name="input_resource_names_map"></a> [resource\_names\_map](#input\_resource\_names\_map) | Map used by the Launch resource naming module. | <pre>map(object({<br/>    name       = string<br/>    max_length = optional(number, 60)<br/>  }))</pre> | <pre>{<br/>  "managed_redis": {<br/>    "max_length": 60,<br/>    "name": "redis"<br/>  },<br/>  "resource_group": {<br/>    "max_length": 80,<br/>    "name": "rg"<br/>  }<br/>}</pre> | no |
| <a name="input_sku_name"></a> [sku\_name](#input\_sku\_name) | Managed Redis SKU name. | `string` | `"Balanced_B1"` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | Tags to apply to resources in the example. | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_managed_redis_hostname"></a> [managed\_redis\_hostname](#output\_managed\_redis\_hostname) | Hostname of the managed redis instance created by the example. |
| <a name="output_managed_redis_id"></a> [managed\_redis\_id](#output\_managed\_redis\_id) | ID of the managed redis instance created by the example. |
| <a name="output_managed_redis_name"></a> [managed\_redis\_name](#output\_managed\_redis\_name) | Name of the managed redis instance created by the example. |
| <a name="output_managed_redis_sku_name"></a> [managed\_redis\_sku\_name](#output\_managed\_redis\_sku\_name) | SKU name requested for the managed redis instance created by the example. |
| <a name="output_resource_group_name"></a> [resource\_group\_name](#output\_resource\_group\_name) | Name of the resource group created by the example. |
<!-- END_TF_DOCS -->
