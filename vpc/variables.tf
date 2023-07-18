variable "name" {
  description = "Name of the VPC"
}

variable "cidr_block" {
  description = "CIDR block for the VPC"
}

variable "subnet_id" {
  description = "ID of the subnet to associate with the route table"
}