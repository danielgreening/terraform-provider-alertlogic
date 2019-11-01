
variable aws_region {
  description = "The AWS region to deploy the appliance in."
  default = "ca-central-1"
}

variable al_account_id {
  description = "The Alert Logic account ID to use."
  default = "134279603"   // CDS US Prod Account
//  default = "134235891"   // CDS Integration Account
}
