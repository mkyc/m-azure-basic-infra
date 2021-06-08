variable "name" {
  description = "String value to use as resources name prefix"
  type        = string
}

variable "rg_name" {
  description = "Resource group name"
  type        = string
}

variable "vnet_id" {
  description = "Virtual network id"
  type        = string
}

variable "vnet_name" {
  description = "Virtual network name that is passed to use subnet data source"
  type        = string
}

variable "location" {
  description = "Azure location to create resources in"
  type        = string
}

variable "admin_username" {
  description = "Admin user name"
  type        = string
  default     = "operations"
}

variable "admin_key_path" {
  description = "The filesystem path to SSH public key"
  type        = string
}

variable "subnets_available" {
  description = "Subnets name => id mapping"
  type        = map(string)
}

variable vm_group {
  description = "VM group definition object"
  type        = object({
    name          = string
    vm_count      = number
    vm_size       = string
    use_public_ip = bool
    subnet_names  = list(string)
    vm_image      = object({
      publisher = string
      offer     = string
      sku       = string
      version   = string
    })
    data_disks    = list(object({
      disk_size_gb = number
      storage_type = string
    }))
  })
}

variable "security_group_id" {
  description = "Security group id for NICs with assigned public IP address"
  type        = string
}
