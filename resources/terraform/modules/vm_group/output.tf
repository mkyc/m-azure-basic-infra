output "private_ips" {
  value = azurerm_network_interface.nic.*.private_ip_address
}

output "public_ips" {
  value = azurerm_public_ip.pubip.*.ip_address
}

output "vm_names" {
  value = azurerm_linux_virtual_machine.vm.*.name
}
