provider "aws" {
    version = "2.7.0"
    region = var.aws_region
}

// create launch configuration for the security appliances to be created
resource "aws_launch_configuration" "appliance_lc" {
  name            = "EVIL LAUNCH CONFIG"
  image_id        = "ami-0b85d4ff00de6a225"
  instance_type   = "c5.large"
}

resource "aws_autoscaling_group" "appliance_asg" {
  count = 200
  name            = "EVIL ASG CONFIG ${count.index}"
  max_size = 0
  min_size = 0
  desired_capacity = 0
  launch_configuration = aws_launch_configuration.appliance_lc.name
  vpc_zone_identifier  = ["subnet-a1f2fed9"]
}