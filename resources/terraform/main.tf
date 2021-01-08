resource "azurerm_resource_group" "rg" {
  name     = "${var.name}-rg"
  location = var.location
}

module "vm_group" {
  source   = "./modules/vm_group"
  count = length(var.vm_groups)

  name           = var.name
  rg_name        = azurerm_resource_group.rg.name
  vnet_id        = azurerm_virtual_network.vnet.id
  vnet_name      = azurerm_virtual_network.vnet.name
  location       = var.location
  admin_key_path = var.rsa_pub_path
  vm_groups       = var.vm_groups
  vm_group_number = count.index

  depends_on = [azurerm_subnet.subnets, azurerm_network_security_group.nsg]
}

/*module "vms" {
  source = "./modules/vm_group"

  instances     = var.vms_count
  name          = var.name
  rg_name       = azurerm_resource_group.rg.name
  vnet_id       = azurerm_virtual_network.vnet.id
  location      = var.location
  service       = azurerm_subnet.subnets[0].name
  use_public_ip = var.use_public_ip
  subnet_id     = azurerm_subnet.subnets[0].id

  tf_key_path = var.rsa_pub_path
}
*/
