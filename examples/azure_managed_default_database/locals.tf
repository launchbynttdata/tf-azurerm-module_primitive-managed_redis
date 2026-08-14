locals {
  managed_redis_name  = module.resource_names["managed_redis"].minimal_random_suffix
  resource_group_name = module.resource_names["resource_group"].minimal_random_suffix
}
