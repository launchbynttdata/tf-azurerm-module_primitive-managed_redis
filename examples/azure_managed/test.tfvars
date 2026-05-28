resource_names_map = {
  resource_group = {
    name       = "rg"
    max_length = 80
  }
  managed_redis = {
    name       = "redis"
    max_length = 24
  }
}
instance_env              = 0
instance_resource         = 0
logical_product_family    = "launch"
logical_product_service   = "redis"
class_env                 = "gotest"
location                  = "eastus"
sku_name                  = "Balanced_B1"
high_availability_enabled = false
public_network_access     = "Enabled"
default_database = {
  access_keys_authentication_enabled = true
  clustering_policy                  = "OSSCluster"
  eviction_policy                    = "VolatileLRU"
}
