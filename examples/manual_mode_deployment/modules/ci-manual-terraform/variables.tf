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
  default     = "Int"
}

variable "vpc_id" {
  description = "Specify the VPC ID where the appliance will be deployed in."
}

variable "subnet_id" {
  description = "Specify the existing subnet ID where the appliance will be deployed in."
}

variable "subnet_type" {
  description = "Select if the subnet is a public or private subnet. Enter Public or Private"
  default     = "Private"
}

variable "instance_type" {
  description = "AlertLogic Security Appliance EC2 instance type. Enter m3.medium, m3.large, m3.xlarge or m3.2xlarge"
  default     = "m3.medium"
}

variable "appliance_number" {
  description = "Number of appliances to be deployed set by the Autoscaling group."
  default     = "1"
}

// the latest AMI is provided by Alert Logic and should have been previously shared with the AWS account deploying the security appliance
variable "aws_amis" {
  type = map(string)
  default = {
    "ap-northeast-1" = "ami-09d23853b4f18fe07",
    "ap-northeast-2" = "ami-0e288c2ea1b3f86d8",
    "ap-south-1" = "ami-0fac5e4fbf52ac20b",
    "ap-southeast-1" = "ami-08f8cd467ffba598b",
    "ap-southeast-2" = "ami-02df21a1cf973219a",
    "ca-central-1" = "ami-001a6190891eee1c8",
    "eu-central-1" = "ami-0121f5d4e5dbc8903",
    "eu-west-1" = "ami-0403bb40be204ad94",
    "eu-west-2" = "ami-0dcded5935564e6d9",
    "eu-west-3" = "ami-0966a2a4d1c1d2158",
    "sa-east-1" = "ami-083f735459efd496b",
    "us-east-1" = "ami-09fcef270e9c5bd2f",
    "us-east-2" = "ami-0e6289d269c1cb4a6",
    "us-west-1" = "ami-02a9c201895d86126",
    "us-west-2" = "ami-0353cfadcb897df91"
  }
}

variable "stack_vaporator" {
  type = map(string)
  default = {
    US-host = "agentapi.cloudinsight.alertlogic.com"
    US-port = "443"
    UK-host = "agentapi.cloudinsight.alertlogic.co.uk"
    UK-port = "443"
    Int-host = "agentapi.product.dev.alertlogic.com"
    Int-port = "443"
  }
}

