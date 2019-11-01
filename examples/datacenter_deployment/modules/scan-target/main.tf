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

// create Scan Target instance
resource "aws_instance" "target_instance" {
  count = var.target_number
  ami = "ami-0b85d4ff00de6a225"  // Redhat in ca-central-1 for now
  instance_type   = var.instance_type
  tags = {
    Name = "dgreening scan target test"
    Terraform = "true"
  }
  vpc_security_group_ids = [aws_security_group.allow_all.id]
  subnet_id = var.subnet_id
}

// create security group to allow security appliance traffic to flow outbound to any destination IP. In general, it will have no rules, which basically allows all traffic outbound but nothing inbound
resource "aws_security_group" "allow_all" {
  description = "dgreening scan target test sg"
  vpc_id      = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "Dgreening Scan Target SG"
    Terraform = "true"
  }
}
