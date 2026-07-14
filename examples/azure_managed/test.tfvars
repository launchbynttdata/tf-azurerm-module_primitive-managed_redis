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

# Keep CI inputs minimal; only pin values that reduce Azure provisioning risk.
high_availability_enabled = false

# Remaining defaults in use:
# - location                  = "eastus"
# - sku_name                  = "Balanced_B1"
# - public_network_access     = "Disabled"
# - default_database          = {}
