variable "name" {
  type        = string
  description = "resources base name"
}

variable "location" {
  type        = string
  description = "resource group location"
}

variable "address_space" {
  type        = list(string)
  description = "vnet address space"
}

variable "subnet_cidrs" {
  type        = list(string)
  description = "subnets cidrs"
}

variable "subnet_names" {
  type        = list(string)
  description = "subnets names"
}