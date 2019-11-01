provider "alertlogic" {
  endpoint = "https://api.product.dev.alertlogic.com"
//  endpoint = "https://api.cloudinsight.alertlogic.com"
}

variable "manual_discovery_role_arn" {}

resource "alertlogic_credential" "test_credential_2" {
  account_id  = "134235891"
  name        = "dgreening-terraform-credential1"
  secret_type = "aws_iam_role"
  secret_arn  = var.manual_discovery_role_arn
}
