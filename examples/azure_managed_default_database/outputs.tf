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

output "resource_group_name" {
  description = "Name of the resource group created by the example."
  value       = module.resource_group.name
}

output "managed_redis_id" {
  description = "ID of the managed redis instance created by the example."
  value       = module.managed_redis.id
}

output "managed_redis_name" {
  description = "Name of the managed redis instance created by the example."
  value       = module.managed_redis.name
}

output "managed_redis_hostname" {
  description = "Hostname of the managed redis instance created by the example."
  value       = module.managed_redis.hostname
}

output "managed_redis_sku_name" {
  description = "SKU name requested for the managed redis instance created by the example."
  value       = var.sku_name
}
