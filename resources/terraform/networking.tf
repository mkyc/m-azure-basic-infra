resource "azurerm_virtual_network" "vnet" {
  name                = "${var.name}-vnet"
  address_space       = var.address_space
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
}

resource "azurerm_subnet" "subnets" {
  count                = length(var.subnets)
  name                 = var.subnets[count.index].name
  address_prefixes     = var.subnets[count.index].address_prefixes
  resource_group_name  = azurerm_resource_group.rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
}
