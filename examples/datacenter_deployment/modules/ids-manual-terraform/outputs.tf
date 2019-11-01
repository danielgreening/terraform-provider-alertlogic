output "sg-id" {
  value = aws_security_group.appliance_sg.id
}

output "subnet-id" {
  value = "${var.subnet_id}/${var.subnet_type}"
}

