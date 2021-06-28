resource "azurerm_resource_group" "rg" {
  // follow <prefix>-<resource_type>-<index> naming pattern
  name     = "${var.name}-rg-0"
  location = var.location
}

module "vm_group" {
  source = "./modules/vm_group"
  count  = length(var.vm_groups)

  // follow <prefix>-<resource_type>-<index> naming pattern
  name                 = "${var.name}-vmg-${count.index}"
  rg_name              = azurerm_resource_group.rg.name
  vnet_id              = azurerm_virtual_network.vnet.id
  vnet_name            = azurerm_virtual_network.vnet.name
  location             = var.location
  admin_key_path       = var.rsa_pub_path
  subnets_available    = {for subnet in azurerm_subnet.subnets: subnet.name => subnet.id}
  vm_group             = var.vm_groups[count.index]
  security_group_id    = length(azurerm_network_security_group.nsg) == 0 ? "" : azurerm_network_security_group.nsg[0].id
  environment_basename = var.name

  depends_on = [
    azurerm_subnet.subnets]
}
