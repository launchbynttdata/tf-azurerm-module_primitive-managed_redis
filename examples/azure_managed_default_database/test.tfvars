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
instance_resource       = 1
logical_product_family  = "launch"
logical_product_service = "redis"
class_env               = "gotest"

# Explicitly set non-default default_database values so test runs exercise
# dynamic block paths with user-provided configuration.
default_database = {
  access_keys_authentication_enabled = true
}
