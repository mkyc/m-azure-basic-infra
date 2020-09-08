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

//TODO count over subnets
module "vms" {
  source = "./modules/vms"

  instances     = var.size
  name          = var.name
  rg_name       = azurerm_resource_group.main_rg.name
  vnet_id       = module.vnet.vnet_id
  location      = var.location
  service       = "all"
  use_public_ip = var.use_public_ip
  subnet_id     = module.vnet.vnet_subnets[0] //TODO from count

  tf_key_path = "/shared/azure_rsa.pub" //TODO template it
}