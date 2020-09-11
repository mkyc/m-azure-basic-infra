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
  value = length(module.vnet.vnet_subnets) > 0 ? module.vnet.vnet_subnets[0] : null
}
