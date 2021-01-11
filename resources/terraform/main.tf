resource "azurerm_resource_group" "rg" {
  name     = "${var.name}-rg"
  location = var.location
}

module "vm_group" {
  source = "./modules/vm_group"
  count  = length(var.vm_groups)

  name              = var.name
  rg_name           = azurerm_resource_group.rg.name
  vnet_id           = azurerm_virtual_network.vnet.id
  vnet_name         = azurerm_virtual_network.vnet.name
  location          = var.location
  admin_key_path    = var.rsa_pub_path
  vm_group          = var.vm_groups[count.index]
  security_group_id = length(azurerm_network_security_group.nsg) == 0 ? "" : azurerm_network_security_group.nsg[0].id

  depends_on = [azurerm_subnet.subnets]
}
