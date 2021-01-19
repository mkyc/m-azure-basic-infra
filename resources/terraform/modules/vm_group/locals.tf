locals {
  nic_vm_subnet_association = setproduct(range(var.vm_group.vm_count), data.azurerm_subnet.subnets.*.id)
  vms_data_disks_product    = setproduct(range(var.vm_group.vm_count), range(length(var.vm_group.data_disks)))
}
