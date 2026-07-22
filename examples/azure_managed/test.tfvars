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

instance_env            = 0
instance_resource       = 0
logical_product_family  = "launch"
logical_product_service = "redis"
class_env               = "gotest"

# Keep CI inputs minimal and rely on module/example defaults for deployability.
# Defaults used:
# - location                  = "eastus2"
# - sku_name                  = "Balanced_B1"
# - high_availability_enabled = true
# - public_network_access     = "Disabled"
# - default_database          = {}
