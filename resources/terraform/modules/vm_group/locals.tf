locals {
  nic_vm_subnet_association = setproduct(range(var.vm_group.vm_count), data.azurerm_subnet.subnets.*.id)
  nic_nsg_vm_association = setproduct(data.azurerm_network_security_group.nsg.*.id, range(var.vm_group.vm_count))
}
