output "vm_group" {
  value = {
    vm_group_name: var.vm_group.name,
    vms: [
    for vm in azurerm_linux_virtual_machine.vm:
    {
      vm_name: vm.name
      public_ip: vm.public_ip_address
      private_ips: vm.private_ip_addresses
      id: vm.id
    }
    ]
    data_disks: [
    for md in azurerm_managed_disk.data_disks:
    {
      id: md.id
      name: md.name
      size: md.disk_size_gb
    }
    ]
    dd_attachments: [
    for dda in azurerm_virtual_machine_data_disk_attachment.vms-dds-attachment:
    {
      managed_disk_id: dda.managed_disk_id
      virtual_machine_id: dda.virtual_machine_id
      lun: dda.lun
    }
    ]
  }
}
