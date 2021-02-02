locals {
  create_public_access_nsg = length(
    [
      for vm_group in var.vm_groups:
      vm_group.use_public_ip
      if vm_group.use_public_ip == true
    ]
  ) > 0
}
