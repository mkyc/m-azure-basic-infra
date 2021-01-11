output rg_name {
  value = azurerm_resource_group.rg.name
}

output vnet_name {
  value = azurerm_virtual_network.vnet.name
}

output "vm_groups" {
  value = module.vm_group.*.vm_group
}
