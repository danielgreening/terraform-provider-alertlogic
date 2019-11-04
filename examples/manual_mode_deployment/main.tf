provider "aws" {
    version = "2.7.0"
    region = var.aws_region
}

provider "alertlogic" {
  endpoint = "https://api.product.dev.alertlogic.com"
}

locals{
  default_tags = {
    Name = "dgreening_terraform_manual_mode"
    Terraform = "true"
  }
}

// Used to retrieve the aws account_id of the authenticated provider
data "aws_caller_identity" "current" {}

//data "alertlogic_account_details" "Test" {}

resource "aws_iam_role" "manual_mode" {
  assume_role_policy = data.aws_iam_policy_document.manual_mode.json
}

resource "aws_iam_role_policy" "manual_mode" {
  policy      = file("private/manual_policy.json")
  role       = aws_iam_role.manual_mode.id
}

data "aws_iam_policy_document" "manual_mode" {
  statement {
    actions = ["sts:AssumeRole"]
    condition {
      test     = "StringEquals"
      variable = "sts:ExternalId"
      values = [var.al_account_id]
    }
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::948063967832:root"]   // Integration
//      identifiers = ["733251395267"] // US Prod TODO hardcoded, need to grab this themis
    }
    effect = "Allow"
  }
  version = "2012-10-17"
}

resource "aws_vpc" "default_vpc" {
  cidr_block = "10.0.0.0/24"

  tags = local.default_tags
}

resource "aws_internet_gateway" "default_gw" {
  vpc_id = aws_vpc.default_vpc.id
  tags = local.default_tags
}

resource "aws_route_table" "default_route_table"{
  vpc_id = aws_vpc.default_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.default_gw.id
  }
}

resource "aws_subnet" "default_subnet" {
  cidr_block = aws_vpc.default_vpc.cidr_block
  vpc_id = aws_vpc.default_vpc.id

  tags = local.default_tags
}

resource "aws_route_table_association" "a" {
  subnet_id      = aws_subnet.default_subnet.id
  route_table_id = aws_route_table.default_route_table.id
}

resource "alertlogic_credential" "test_credential_2" {
  name        = "dgreening-terraform-credential1"
  secret_type = "aws_iam_role"
  secret_arn  = aws_iam_role.manual_mode.arn
}

resource "alertlogic_deployment" "test_deployment_1" {
  account_id    = var.al_account_id
  name          = "dgreening-terraform-1"
  platform_type = "aws"
  platform_id   = data.aws_caller_identity.current.account_id
  mode          = "manual"
  scope_include {
    type = "region"
    key  = var.aws_region
    policy_id = "D12D5E67-166C-474F-87AA-6F86FC9FB9BC"
  }
  credential {
      id = alertlogic_credential.test_credential_2.id
      purpose = "discover"
  }
  cloud_defender_location_id = "defender-us-denver"
  cloud_defender_enabled = false
  depends_on = ["aws_iam_role_policy.manual_mode"]
}

module "scan_appliance" {
  source = "./modules/ci-manual-terraform"
  account_id = var.al_account_id
  aws_region = var.aws_region
  deployment_id = alertlogic_deployment.test_deployment_1.id
  subnet_id = aws_subnet.default_subnet.id
  vpc_id = aws_vpc.default_vpc.id
  instance_type = "c5.large"
}

module "scan_targets" {
  source = "./modules/scan-target"
  aws_region = var.aws_region
  subnet_id = aws_subnet.default_subnet.id
  vpc_id = aws_vpc.default_vpc.id
  target_number = "2"
}