variable "aws_region" {
  description = "The AWS region to deploy the appliance in."
}

variable "vpc_id" {
  description = "Specify the VPC ID where the appliance will be deployed in."
}

variable "subnet_id" {
  description = "Specify the existing subnet ID where the appliance will be deployed in."
}

variable "instance_type" {
  description = "Scan Target EC2 instance type."
  default     = "t2.micro"
}

variable "target_number" {
  description = "Number of scan target instances to be deployed set by the Autoscaling group."
  default     = "1"
}
