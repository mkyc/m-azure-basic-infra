variable "vms_count" {
  type = number
}

variable "use_public_ip" {
  type = bool
}

variable "name" {
  type = string
}

variable "location" {
  type = string
}

variable "address_space" {
  type = list(string)
}

variable "rsa_pub_path" {
  type = string
}

variable "subnets" {
  type    = list(object({
    name             = string
    address_prefixes = list(string)
  }))
  default = [
    {
      name             = "some-subnet1"
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
