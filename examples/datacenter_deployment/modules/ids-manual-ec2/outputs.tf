output "sg-id" {
  value = aws_security_group.appliance_sg.id
}

output "subnet-id" {
  value = "${var.subnet_id}/${var.subnet_type}"
}

output "appliance-private-ips" {
  value = aws_instance.appliance_instance.*.private_ip
}

output "appliance-public-ips" {
  value = aws_instance.appliance_instance.*.public_ip
}

output "appliance-private-dnsnames" {
  value = aws_instance.appliance_instance.*.private_dns
}

output "appliance-public-dnsnames" {
  value = aws_instance.appliance_instance.*.public_dns
}

output "appliance-instance_ids" {
  value = aws_instance.appliance_instance.*.id
}

