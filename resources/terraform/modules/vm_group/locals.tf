locals {
  nic_vm_subnet_association = setproduct(range(var.vm_group.vm_count), [for name in var.vm_group.subnet_names: var.subnets_available[name] ])
  vms_data_disks_product    = setproduct(range(var.vm_group.vm_count), range(length(var.vm_group.data_disks)))
}
