resource "azurerm_resource_group" "main_rg" {
  name     = "${var.name}-rg"
  location = var.location
}

module "vnet" {
  source  = "Azure/vnet/azurerm"
  version = "2.1.0"

  resource_group_name = azurerm_resource_group.main_rg.name
  address_space       = var.address_space
  subnet_prefixes     = var.subnet_cidrs
  subnet_names        = var.subnet_names
  vnet_name           = "${var.name}-vnet"

  tags = {}
}
