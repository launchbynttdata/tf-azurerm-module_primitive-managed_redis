locals {
  tags = merge(
    var.tags,
    {
      provisioner   = "terraform"
      resource_name = var.name
    }
  )
}
