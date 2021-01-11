output "vm_group" {
  value = {
    vm_group: var.vm_group.name,
    vms: [
    for vm in azurerm_linux_virtual_machine.vm:
    {
      vm_name: vm.name
      public_ip: vm.public_ip_address
      private_ips: vm.private_ip_addresses
    }
    ]
  }
}
