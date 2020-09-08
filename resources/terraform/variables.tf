variable "size" {
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

variable "subnet_cidrs" {
  type = list(string)
}

variable "subnet_names" {
  type = list(string)
}
