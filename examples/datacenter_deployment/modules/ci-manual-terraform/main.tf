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


// initial setup of the security appliance, a bash bootstrap script is executed using AWS userdata functionality
data "template_file" "userdata" {
  template = file("${path.module}/userdata.tpl")
  vars = {
    stack_host    = var.stack_vaporator["${var.stack}-host"]
    stack_port    = var.stack_vaporator["${var.stack}-port"]
    account_id    = var.account_id
    deployment_id = var.deployment_id
  }
}

// create launch configuration for the security appliances to be created
resource "aws_launch_configuration" "appliance_lc" {
  name            = "AlertLogic Security Launch Configuration ${var.account_id}/${var.deployment_id}/${var.vpc_id}"
  image_id        = var.aws_amis[var.aws_region]
  security_groups = [aws_security_group.appliance_sg.id]
  instance_type   = var.instance_type

  associate_public_ip_address = var.subnet_type == "Public" ? true : false
  user_data                   = data.template_file.userdata.rendered

  lifecycle {
    create_before_destroy = true
  }
}

// create ASG to have the specified amount of security appliances up and running using the created launch configuration
resource "aws_autoscaling_group" "appliance_asg" {
  name                 = "AlertLogic Security Autoscaling Group ${var.account_id}/${var.deployment_id}/${var.vpc_id}"
  max_size             = var.appliance_number
  min_size             = var.appliance_number
  desired_capacity     = var.appliance_number
  force_delete         = true
  launch_configuration = aws_launch_configuration.appliance_lc.name
  vpc_zone_identifier  = [var.subnet_id]

  tags = [
    {
      key = "Name"
      value = "AlertLogic Security Appliance"
      propagate_at_launch = "true"
    },
    {
      key = "AlertLogic-AccountID"
      value = var.account_id
      propagate_at_launch = "true"
    },
    {
      key = "AlertLogic-EnvironmentID"
      value = var.deployment_id
      propagate_at_launch = "true"
    },
    {
      key = "AlertLogic"
      value = "Security"
      propagate_at_launch = "true"
    }
  ]
}

// create security group to allow security appliance traffic to flow outbound to any destination IP. In general, it will have no rules, which basically allows all traffic outbound but nothing inbound
resource "aws_security_group" "appliance_sg" {
  name        = "AlertLogic Security Group ${var.account_id}_${var.deployment_id}_${var.vpc_id}_${var.subnet_type}"
  description = "Alert Logic Security Group ${var.account_id}_${var.deployment_id}_${var.vpc_id}_${var.subnet_type}"
  vpc_id      = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "AlertLogic Security Group"
    AlertLogic-AccountID = var.account_id
    AlertLogic-EnvironmentID = var.deployment_id
    AlertLogic = "Security"
  }
}

