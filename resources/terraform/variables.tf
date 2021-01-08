variable "name" {
  type = string
  default = "atsikham-azbi"
}

variable "location" {
  type = string
  default = "northeurope"
}

variable "address_space" {
  type = list(string)
  default = ["10.0.0.0/16"]
}

variable "rsa_pub_path" {
  type = string
  default = "/Users/anatolytikhomirov/repos/m-azure-basic-infrastructure/vms_rsa.pub"
}

variable "subnets" {
  type = list(object({
    name             = string
    address_prefixes = list(string)
  }))
  default = [
    {
      name             = "subnet1"
      address_prefixes = [
        "10.0.1.0/24"
      ]
    },
    {
      name             = "subnet2"
      address_prefixes = [
        "10.0.2.0/24"
      ]
    },
    {
      name             = "subnet3"
      address_prefixes = [
        "10.0.3.0/24"
      ]
    }
  ]
  validation {
    condition     = length(var.subnets) > 0
    error_message = "Subnets list needs to have at least one element."
  }
}

variable "network_security_groups" {
  description = "List of the network security groups that can be shared between VM groups"
  type        = list(object({
    name          = string
    security_rule = object({
      name                       = string
      priority                   = number
      direction                  = string
      access                     = string
      protocol                   = string
      source_port_range          = string
      destination_port_range     = string
      source_address_prefix      = string
      destination_address_prefix = string
    })
  }))
  default     = [{
    name          = "default-ssh"
    security_rule = {
      name                       = "SSH"
      priority                   = 100
      direction                  = "Inbound"
      access                     = "Allow"
      protocol                   = "Tcp"
      source_port_range          = "*"
      destination_port_range     = "22"
      source_address_prefix      = "*"
      destination_address_prefix = "*"
    }
    },
    {
      name          = "default-http"
      security_rule = {
        name                       = "HTTP"
        priority                   = 200
        direction                  = "Inbound"
        access                     = "Allow"
        protocol                   = "Tcp"
        source_port_range          = "*"
        destination_port_range     = "80"
        source_address_prefix      = "*"
        destination_address_prefix = "*"
      }
    }
  ]
}

variable vm_groups {
  type = list(object({
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
  }))
  default = [
    {
      name          = "vm-group1"
      vm_count      = 2
      vm_size       = "Standard_DS2_v2"
      use_public_ip = true
      subnet_names  = ["subnet1"]
      nsg_names     = ["default-ssh", "default-http"]
      image         = {
        publisher = "Canonical"
        offer     = "UbuntuServer"
        sku       = "18.04-LTS"
        version   = "18.04.202006101"
      }
    },
    {
      name          = "vm-group2"
      vm_count      = 2
      vm_size       = "Standard_DS2_v2"
      use_public_ip = true
      subnet_names  = ["subnet2", "subnet3"]
      nsg_names     = ["default-ssh", "default-http"]
      image         = {
        publisher = "Canonical"
        offer     = "UbuntuServer"
        sku       = "18.04-LTS"
        version   = "18.04.202006101"
      }
    }
  ]
}
