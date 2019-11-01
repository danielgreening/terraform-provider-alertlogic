variable "aws_region" {
  description = "The AWS region to deploy the appliance in."
}

variable "account_id" {
  description = "AlertLogic Account ID."
}

variable "deployment_id" {
  description = "AlertLogic cloudinsight Deployment ID."
}

variable "stack" {
  description = "AlertLogic DataCenter where the appliance will be deployed in. Enter US or UK"
  default     = "US"
}

variable "vpc_id" {
  description = "Specify the VPC ID where the appliance will be deployed in."
}

variable "subnet_id" {
  description = "Specify the existing subnet ID where the appliance will be deployed in."
}

variable "subnet_type" {
  description = "Select if the subnet is a public or private subnet. Enter Public or Private"
  default     = "Public"
}

variable "instance_type" {
  description = "AlertLogic IDS Appliance EC2 instance type. Enter m3.medium, m3.large, m3.xlarge or m3.2xlarge"
  default     = "c5.large"
}

variable "appliance_number" {
  description = "Number of appliances to be deployed set by the Autoscaling group."
  default     = "1"
}

// the latest AMI is provided by Alert Logic and should have been previously shared with the AWS account deploying the security appliance
variable "aws_amis" {
  type = map(string)
  default = {
    ca-central-1   = "ami-01fe6df823a1bac12"
  }
}


