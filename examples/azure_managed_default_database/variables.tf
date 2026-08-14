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

variable "resource_names_map" {
  description = "Map used by the Launch resource naming module."
  type = map(object({
    name       = string
    max_length = optional(number, 60)
  }))

  default = {
    resource_group = {
      name       = "rg"
      max_length = 80
    }
    managed_redis = {
      name       = "redis"
      max_length = 60
    }
  }
}

variable "instance_env" {
  description = "Environment instance number used for name generation."
  type        = number
  default     = 0
}

variable "instance_resource" {
  description = "Resource instance number used for name generation."
  type        = number
  default     = 0
}

variable "logical_product_family" {
  description = "Logical product family used for name generation."
  type        = string
  default     = "launch"
}

variable "logical_product_service" {
  description = "Logical product service used for name generation."
  type        = string
  default     = "redis"
}

variable "class_env" {
  description = "Deployment environment used for name generation."
  type        = string
  default     = "dev"
}

variable "location" {
  description = "Azure region for the example deployment."
  type        = string
  default     = "eastus"
}

variable "sku_name" {
  description = "Managed Redis SKU name."
  type        = string
  default     = "Balanced_B1"
}

variable "high_availability_enabled" {
  description = "Whether high availability is enabled for the managed redis instance."
  type        = bool
  default     = true
}

variable "public_network_access" {
  description = "Public network access mode for Managed Redis."
  type        = string
  default     = "Disabled"
}

variable "identity" {
  description = "Optional managed identity settings."
  type = object({
    type         = string
    identity_ids = optional(list(string))
  })
  default = null
}

variable "customer_managed_key" {
  description = "Optional CMK encryption settings."
  type = object({
    key_vault_key_id          = string
    user_assigned_identity_id = string
  })
  default = null
}

variable "default_database" {
  description = "Default database configuration for managed redis."
  type = object({
    access_keys_authentication_enabled            = optional(bool, false)
    client_protocol                               = optional(string, "Encrypted")
    clustering_policy                             = optional(string, "OSSCluster")
    eviction_policy                               = optional(string, "VolatileLRU")
    geo_replication_group_name                    = optional(string)
    persistence_append_only_file_backup_frequency = optional(string)
    persistence_redis_database_backup_frequency   = optional(string)
    modules = optional(list(object({
      name = string
      args = optional(string)
    })), [])
  })
  default = {}
}

variable "tags" {
  description = "Tags to apply to resources in the example."
  type        = map(string)
  default     = {}
}
