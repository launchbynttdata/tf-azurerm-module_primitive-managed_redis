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

resource "azurerm_managed_redis" "redis" {
  name                = var.name
  location            = var.location
  resource_group_name = var.resource_group_name
  sku_name            = var.sku_name

  high_availability_enabled = var.high_availability_enabled
  public_network_access     = var.public_network_access

  dynamic "identity" {
    for_each = var.identity != null ? [var.identity] : []
    content {
      type         = identity.value.type
      identity_ids = identity.value.identity_ids
    }
  }

  dynamic "customer_managed_key" {
    for_each = var.customer_managed_key != null ? [var.customer_managed_key] : []
    content {
      key_vault_key_id          = customer_managed_key.value.key_vault_key_id
      user_assigned_identity_id = customer_managed_key.value.user_assigned_identity_id
    }
  }

  dynamic "default_database" {
    for_each = [var.default_database]
    content {
      access_keys_authentication_enabled            = default_database.value.access_keys_authentication_enabled
      client_protocol                               = default_database.value.client_protocol
      clustering_policy                             = default_database.value.clustering_policy
      eviction_policy                               = default_database.value.eviction_policy
      geo_replication_group_name                    = default_database.value.geo_replication_group_name
      persistence_append_only_file_backup_frequency = default_database.value.persistence_append_only_file_backup_frequency
      persistence_redis_database_backup_frequency   = default_database.value.persistence_redis_database_backup_frequency

      dynamic "module" {
        for_each = default_database.value.modules
        content {
          name = module.value.name
          args = module.value.args
        }
      }
    }
  }

  tags = local.tags
}
