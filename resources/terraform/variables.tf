variable "name" {
  description = "String value to use as resources name prefix"
  type        = string
}

variable "location" {
  description = "Azure location to create resources in"
  type        = string
  default     = "northeurope"
}

variable "address_space" {
  description = "VNet address space"
  type        = list(string)
  default     = ["10.0.0.0/16"]
}

variable "rsa_pub_path" {
  description = "The filesystem path to SSH public key"
  type        = string
}

variable "subnets" {
  description = "The list of subnets definition objects"
  type        = list(object({
    name             = string
    address_prefixes = list(string)
  }))
  default     = [
    {
      name             = "subnet0"
      address_prefixes = [
        "10.0.1.0/24"
      ]
    }
  ]
  validation {
    condition     = length(var.subnets) > 0
    error_message = "Subnets list needs to have at least one element."
  }
}

variable vm_groups {
  description = "The list of VM group definition objects"
  type        = list(object({
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
    data_disks = list(object({
      disk_size_gb = number
    }))
  }))
  default     = [
    {
      name          = "vm-group0"
      vm_count      = 1
      vm_size       = "Standard_DS2_v2"
      use_public_ip = true
      subnet_names  = ["subnet0"]
      vm_image      = {
        publisher = "Canonical"
        offer     = "UbuntuServer"
        sku       = "18.04-LTS"
        version   = "18.04.202006101"
      }
      data_disks = [
        {
          disk_size_gb = 10
        }
      ]
    }
  ]
}
