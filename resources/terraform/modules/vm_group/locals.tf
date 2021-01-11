locals {
  nic_vm_subnet_association = setproduct(range(var.vm_group.vm_count), data.azurerm_subnet.subnets.*.id)
}
