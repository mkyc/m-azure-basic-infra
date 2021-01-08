variable "name" {
  type = string
}

variable "rg_name" {
  type = string
}

variable "vnet_id" {
  type = string
}

variable "vnet_name" {
  description = "Virtual network name that is passed to use subnet data source"
  type        = string
}

variable "location" {
  type = string
}

variable "admin_username" {
  type    = string
  default = "operations"
}

variable "admin_key_path" {
  type = string
}

variable vm_group {
  type = object({
    name          = string
    vm_count      = number
    vm_size       = string
    use_public_ip = bool
    subnet_names  = list(string)
    nsg_names     = list(string)
    image         = object({
      publisher = string
      offer     = string
      sku       = string
      version   = string
    })
  })
}
