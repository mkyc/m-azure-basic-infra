locals {
  subnet_ids = zipmap(var.vm_groups[var.vm_group_number].subnet_names, data.azurerm_subnet.subnets.*.id)
  nic_vm_subnet_association = setproduct(range(var.vm_groups[var.vm_group_number].vm_count), data.azurerm_subnet.subnets.*.id)
  nic_nsg_vm_association = setproduct(data.azurerm_network_security_group.nsg.*.id, range(var.vm_groups[var.vm_group_number].vm_count))
}