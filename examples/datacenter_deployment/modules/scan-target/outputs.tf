output "target-private-ips" {
  value = aws_instance.target_instance.*.private_ip
}

output "target-public-ips" {
  value = aws_instance.target_instance.*.public_ip
}

output "target-private-dnsnames" {
  value = aws_instance.target_instance.*.private_dns
}

output "target-public-dnsnames" {
  value = aws_instance.target_instance.*.public_dns
}

