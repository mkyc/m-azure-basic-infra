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

resource "azurerm_network_security_group" "nsg" {
  count               = length(var.network_security_groups)
  name                = var.network_security_groups[count.index].name
  location            = var.location
  resource_group_name = azurerm_resource_group.rg.name

  security_rule {
    name                       = var.network_security_groups[count.index].security_rule.name
    priority                   = var.network_security_groups[count.index].security_rule.priority
    direction                  = var.network_security_groups[count.index].security_rule.direction
    access                     = var.network_security_groups[count.index].security_rule.access
    protocol                   = var.network_security_groups[count.index].security_rule.protocol
    source_port_range          = var.network_security_groups[count.index].security_rule.source_port_range
    destination_port_range     = var.network_security_groups[count.index].security_rule.destination_port_range
    source_address_prefix      = var.network_security_groups[count.index].security_rule.source_address_prefix
    destination_address_prefix = var.network_security_groups[count.index].security_rule.destination_address_prefix
  }
}
