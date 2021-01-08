output "private_ips" {
  value = module.vm_group.*.private_ips
}

output "public_ips" {
  value = module.vm_group.*.public_ips
}

output "vm_names" {
  value = module.vm_group.*.vm_names
}

output rg_name {
  value = azurerm_resource_group.rg.name
}

output vnet_name {
  value = azurerm_virtual_network.vnet.name
}
