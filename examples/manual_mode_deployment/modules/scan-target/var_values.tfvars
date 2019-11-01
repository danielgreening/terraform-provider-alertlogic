/*
# -------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# -------------------------------------------------------------------------------------------------------------------
*/

aws_region = "xx-xxxx-x" // The AWS region to deploy the appliance in
vpc_id = "vpc-xxxxxxxx" // Specify the VPC ID where the appliance will be deployed in
subnet_id = "subnet-xxxxxxxx" // Specify the existing subnet ID where the appliance will be deployed in
instance_type = "m3.medium" // AlertLogic Security Appliance EC2 instance type. Enter m3.medium, m3.large, m3.xlarge or m3.2xlarge
target_number = "1" // Number of appliances to be deployed set by the Autoscaling group
