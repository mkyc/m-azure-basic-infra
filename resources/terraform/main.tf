resource "azurerm_resource_group" "main_rg" {
  name     = var.rg-name
  location = var.location
}
