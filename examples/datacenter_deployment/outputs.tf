
output "al_account_id" {
  value = alertlogic_deployment.dc_deployment.account_id
}

output "al_deployment_id" {
  value = alertlogic_deployment.dc_deployment.id
}

output "appliance-private-ips" {
  value = module.ids-manual-ec2.appliance-private-ips
}

output "appliance-public-ips" {
  value = module.ids-manual-ec2.appliance-public-ips
}

output "appliance-instance-ids" {
  value = module.ids-manual-ec2.appliance-instance_ids
}

