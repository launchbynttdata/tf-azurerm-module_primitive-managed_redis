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

output "id" {
  description = "The ID of the Managed Redis instance."
  value       = azurerm_managed_redis.redis.id
}

output "name" {
  description = "The name of the Managed Redis instance."
  value       = azurerm_managed_redis.redis.name
}

output "hostname" {
  description = "Hostname of the Managed Redis endpoint."
  value       = azurerm_managed_redis.redis.hostname
}

output "database_id" {
  description = "The ID of the default Managed Redis database."
  value       = azurerm_managed_redis.redis.default_database[0].id
}

output "database_port" {
  description = "The port of the default Managed Redis database."
  value       = azurerm_managed_redis.redis.default_database[0].port
}

output "primary_access_key" {
  description = "Primary access key for the default database when access key authentication is enabled."
  value       = azurerm_managed_redis.redis.default_database[0].primary_access_key
  sensitive   = true
}

output "secondary_access_key" {
  description = "Secondary access key for the default database when access key authentication is enabled."
  value       = azurerm_managed_redis.redis.default_database[0].secondary_access_key
  sensitive   = true
}
