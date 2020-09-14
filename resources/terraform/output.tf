output "private_ips" {
  value = module.vms.private_ips
}

output "public_ips" {
  value = module.vms.public_ips
}

output "vm_names" {
  value = module.vms.vm_names
}

output "subnet_id" {
  value = azurerm_subnet.subnet.id
}
