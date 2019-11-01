provider "aws" {
    version = "2.7.0"
    region = "eu-west-2"
    profile = "personal"
}

provider "alertlogic" {
  endpoint = "https://api.product.dev.alertlogic.com"
//  endpoint = "https://api.cloudinsight.alertlogic.com"
}

// Used to retrieve the aws account_id of the authenticated provider
data "aws_caller_identity" "current" {}

//data "alertlogic_account_details" "Test" {}

resource "aws_iam_role" "automatic_mode" {
  assume_role_policy = data.aws_iam_policy_document.automatic-mode.json
}

resource "aws_iam_role_policy" "automatic_mode" {
  policy      = file("automatic_mode_policy.json")
  role       = aws_iam_role.automatic_mode.id
}

data "aws_iam_policy_document" "automatic-mode" {
  statement {
    actions = ["sts:AssumeRole"]
    condition {
      test     = "StringEquals"
      variable = "sts:ExternalId"
      values = ["134235891"]
    }
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::948063967832:root"]   // Integration
//      identifiers = ["arn:aws:iam::733251395267:root"] // US Prod TODO hardcoded, need to grab this themis
    }
    effect = "Allow"
  }
  version = "2012-10-17"
}

resource "alertlogic_credential" "test_credential_2" {
  account_id  = "134235891"
  name        = "dgreening-terraform-credential1"
  secret_type = "aws_iam_role"
  secret_arn  = aws_iam_role.automatic_mode.arn
  depends_on = ["aws_iam_role_policy.automatic_mode"]
}
