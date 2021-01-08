# Get subnet IDs by names
data "azurerm_subnet" "subnets" {
  count                = length(var.vm_groups[var.vm_group_number].subnet_names)
  name                 = var.vm_groups[var.vm_group_number].subnet_names[count.index]
  virtual_network_name = var.vnet_name
  resource_group_name  = var.rg_name
}

# Get nsg IDs by names
data "azurerm_network_security_group" "nsg" {
  count               = length(var.vm_groups[var.vm_group_number].nsg_names)
  name                = var.vm_groups[var.vm_group_number].nsg_names[count.index]
  resource_group_name = var.rg_name
}

# Allocate public IPs for each VM
resource "azurerm_public_ip" "pubip" {
  count                   = var.vm_groups[var.vm_group_number].use_public_ip != true ? 0 : var.vm_groups[var.vm_group_number].vm_count
  name                    = "${var.name}-${var.vm_groups[var.vm_group_number].name}-${count.index}-pubip"
  location                = var.location
  resource_group_name     = var.rg_name
  allocation_method       = "Static"
  idle_timeout_in_minutes = "30"
  sku                     = "Standard"
}

# Create NICs for each VM in each subnet
resource "azurerm_network_interface" "nic" {
#  count                         = length(local.nic_subnet_vm_association)
  count                         = length(local.nic_vm_subnet_association)
  name                          = "${var.name}-${var.vm_groups[var.vm_group_number].name}-${count.index}-nic"
  location                      = var.location
  resource_group_name           = var.rg_name
  enable_accelerated_networking = "false"

  ip_configuration {
    name                          = "vm-ipconf-0"
#    subnet_id                     = local.nic_subnet_vm_association[count.index][0]
    subnet_id                     = local.nic_vm_subnet_association[count.index][1]
    private_ip_address_allocation = "Dynamic"
    # Assign public IPs only to the NICs in the first subnet
#    public_ip_address_id          = length(azurerm_public_ip.pubip.*.id) > 0 && count.index < var.vm_group.vm_count ? azurerm_public_ip.pubip[count.index].id : ""
    public_ip_address_id          = length(azurerm_public_ip.pubip.*.id) > 0 && count.index % var.vm_groups[var.vm_group_number].vm_count == 0 ? azurerm_public_ip.pubip[floor(count.index / var.vm_groups[var.vm_group_number].vm_count)].id : ""
  }
}

# Associate NICs that have public IPs assigned with network security groups
resource "azurerm_network_interface_security_group_association" "nic-nsg-assoc" {
  count                     = var.vm_groups[var.vm_group_number].use_public_ip != true ? 0 : length(local.nic_nsg_vm_association)
  network_interface_id      = azurerm_network_interface.nic[local.nic_nsg_vm_association[count.index][1]].id
  network_security_group_id = local.nic_nsg_vm_association[count.index][0]
}

# Create VMs
resource "azurerm_linux_virtual_machine" "vm" {
  count                 = var.vm_groups[var.vm_group_number].vm_count
  name                  = "${var.name}-${var.vm_groups[var.vm_group_number].name}-${count.index}"
  location              = var.location
  resource_group_name   = var.rg_name
  size                  = var.vm_groups[var.vm_group_number].vm_size
  network_interface_ids = slice(azurerm_network_interface.nic.*.id, count.index * length(var.vm_groups[var.vm_group_number].subnet_names), count.index * (length(var.vm_groups[var.vm_group_number].subnet_names) + 1))

  disable_password_authentication = true

  admin_username = var.admin_username

  admin_ssh_key {
    username   = var.admin_username
    public_key = file(var.admin_key_path)
  }

  source_image_reference {
    publisher = var.vm_groups[var.vm_group_number].image.publisher
    offer     = var.vm_groups[var.vm_group_number].image.offer
    sku       = var.vm_groups[var.vm_group_number].image.sku
    version   = var.vm_groups[var.vm_group_number].image.version
  }

  os_disk {
    name                 = "${var.name}-${var.vm_groups[var.vm_group_number].name}-${count.index}-disk"
    caching              = "ReadWrite"
    disk_size_gb         = "32"
    storage_account_type = "Premium_LRS"
  }
}
