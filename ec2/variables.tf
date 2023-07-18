variable "name" {
  description = "Name of the EC2 instance"
}

variable "ami" {
  description = "ID of the AMI to use"
}

variable "instance_type" {
  description = "Type of EC2 instance"
}

variable "subnet_id" {
  description = "ID of the subnet to launch the instance in"
}

variable "key_name" {
  description = "Name of the key pair to use for SSH access"
}

variable "security_group_id" {
  description = "ID of the security group to use"
}