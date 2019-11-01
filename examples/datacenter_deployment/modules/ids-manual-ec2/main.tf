/*
TerraForm Template (v1.0): Creates AlertLogic security appliance for Cloud Insight and launch configuration inside an existing VPC/Subnet.
The artifacts created are: A Security Group, a Launch Configuration for the security appliances that uses the Security Group and an Autoscaling Group that uses the Launch Configuration inside an existing Subnet.

Usage:
1. copy this template and other associated files (outputs.tf, userdata.tpl, inputs.tf) to your working directory, i.e. ~/aws-ci-security-appliance-manual
2. add all the required variable values in a separate file under the same directory, i.e. var_values.tfvars
3. run TerraForm apply
   > terraform apply -var-file var_values.tfvars

Provider configuration:
 Credentials can be provided from separate file (default file name is credentials.tf)
 Variables can be loaded from separate file or passed as parameters below.
*/

// Specify the provider and access details below

// create IDS appliance instance
resource "aws_instance" "appliance_instance" {
  count = var.appliance_number
  ami = var.aws_amis[var.aws_region]
  instance_type   = var.instance_type
  associate_public_ip_address = var.subnet_type == "Public" ? true : false
  tags = {
    Name = "AlertLogic IDS Security Appliance"
    AlertLogic-AccountID = var.account_id
    AlertLogic-EnvironmentID = var.deployment_id
    AlertLogic = "Security"
  }
  vpc_security_group_ids = [aws_security_group.appliance_sg.id]
  subnet_id = var.subnet_id

}

// create security group to allow security appliance traffic to flow outbound to any destination IP. In general, it will have no rules, which basically allows all traffic outbound but nothing inbound
resource "aws_security_group" "appliance_sg" {
  name        = "AlertLogic IDS Security Group ${var.account_id}_${var.deployment_id}_${var.vpc_id}_${var.subnet_type}"
  description = "Alert Logic IDS Security Group ${var.account_id}_${var.deployment_id}_${var.vpc_id}_${var.subnet_type}"
  vpc_id      = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22223
    to_port     = 22223
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "AlertLogic IDS Security Group"
    AlertLogic-AccountID = var.account_id
    AlertLogic-EnvironmentID = var.deployment_id
    AlertLogic = "Security"
  }
}
