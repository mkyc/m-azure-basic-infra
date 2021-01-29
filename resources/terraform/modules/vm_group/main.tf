# Get subnet IDs by names
data "azurerm_subnet" "subnets" {
  count                = length(var.vm_group.subnet_names)
  name                 = var.vm_group.subnet_names[count.index]
  virtual_network_name = var.vnet_name
  resource_group_name  = var.rg_name
}

# Allocate public IPs, 1 per VM
resource "azurerm_public_ip" "pubip" {
  count                   = var.vm_group.use_public_ip != true ? 0 : var.vm_group.vm_count
  name                    = "${var.name}-${var.vm_group.name}-${count.index}-pubip"
  location                = var.location
  resource_group_name     = var.rg_name
  allocation_method       = "Static"
  idle_timeout_in_minutes = "30"
  sku                     = "Standard"
}

/*
Create NICs for each VM in each subnet
NICs order example:
+----------------+----------------+
| vm0            | vm1            |
+================+================+
| net0 net1 net2 | net0 net1 net2 |
+----------------+----------------+
| nic0 nic1 nic2 | nic3 nic4 nic5 |
+----------------+----------------+
*/
resource "azurerm_network_interface" "nic" {
  count                         = length(local.nic_vm_subnet_association)
  name                          = "${var.name}-${var.vm_group.name}-${count.index}-nic"
  location                      = var.location
  resource_group_name           = var.rg_name
  enable_accelerated_networking = "false"

  ip_configuration {
    name                          = "${var.name}-${var.vm_group.name}-${count.index}-ipconf"
    subnet_id                     = local.nic_vm_subnet_association[count.index][1]
    private_ip_address_allocation = "Dynamic"
    # Assign public IPs only to the first NIC of the VM
    public_ip_address_id          = length(azurerm_public_ip.pubip) > 0 && count.index % length(var.vm_group.subnet_names) == 0 ? azurerm_public_ip.pubip[count.index / length(var.vm_group.subnet_names)].id : ""
  }
}

# Associate NICs that have public IPs assigned with network security group (1 per VM)
resource "azurerm_network_interface_security_group_association" "nic-nsg-assoc" {
  count                     = var.vm_group.use_public_ip == true ? var.vm_group.vm_count : 0
  network_interface_id      = azurerm_network_interface.nic[count.index * length(var.vm_group.subnet_names)].id
  network_security_group_id = var.security_group_id
}

# Create VMs
resource "azurerm_linux_virtual_machine" "vm" {
  count                           = var.vm_group.vm_count
  name                            = "${var.name}-${var.vm_group.name}-${count.index}"
  location                        = var.location
  resource_group_name             = var.rg_name
  size                            = var.vm_group.vm_size
  network_interface_ids           = slice(azurerm_network_interface.nic.*.id, count.index * length(var.vm_group.subnet_names), count.index * length(var.vm_group.subnet_names) + length(var.vm_group.subnet_names))
  disable_password_authentication = true
  admin_username                  = var.admin_username

  admin_ssh_key {
    username   = var.admin_username
    public_key = file(var.admin_key_path)
  }

  source_image_reference {
    publisher = var.vm_group.vm_image.publisher
    offer     = var.vm_group.vm_image.offer
    sku       = var.vm_group.vm_image.sku
    version   = var.vm_group.vm_image.version
  }

  os_disk {
    name                 = "${var.name}-${var.vm_group.name}-${count.index}-os-disk"
    caching              = "ReadWrite"
    storage_account_type = "Premium_LRS"
  }
}

#Create managed data disks
resource "azurerm_managed_disk" "data_disks" {
  count                = length(local.vms_data_disks_product)
  name                 = "${var.name}-${var.vm_group.name}-${local.vms_data_disks_product[count.index][0]}-data-disk-${local.vms_data_disks_product[count.index][1]}"
  location             = var.location
  resource_group_name  = var.rg_name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = var.vm_group.data_disks[local.vms_data_disks_product[count.index][1]].disk_size_gb
}

#Attach data disk to vm
resource "azurerm_virtual_machine_data_disk_attachment" "vms-dds-attachment" {
  count              = length(local.vms_data_disks_product)
  caching            = "ReadWrite"
  lun                = 10 + local.vms_data_disks_product[count.index][1]
  managed_disk_id    = azurerm_managed_disk.data_disks[count.index].id
  virtual_machine_id = azurerm_linux_virtual_machine.vm[local.vms_data_disks_product[count.index][0]].id
}
