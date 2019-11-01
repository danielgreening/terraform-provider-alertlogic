
provider aws {
    version = "2.7.0"
    region = var.aws_region
}

provider alertlogic {
  endpoint = "https://api.product.dev.alertlogic.com"
//  endpoint = "https://api.cloudinsight.alertlogic.com"
//  access_key = "bfe43a21c037810b"
//  secret_key = "b5aaad12e61cdbab1761eb04b2694827f2c6c235c27734b52c245bbb106c07d7"
}

locals{
  default_tags = {
    Name = "dgreening_dc_network"
    Terraform = "true"
  }
}

// Used to retrieve the aws account_id of the authenticated provider
data "aws_caller_identity" "current" {}

resource "aws_vpc" "default_network" {
  cidr_block = "10.0.0.0/24"

  tags = local.default_tags
}

resource "aws_internet_gateway" "default_gw" {
  vpc_id = aws_vpc.default_network.id
  tags = local.default_tags
}

resource "aws_route_table" "default_route_table"{
  vpc_id = aws_vpc.default_network.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.default_gw.id
  }
}

resource "aws_subnet" "default_subnet" {
  cidr_block = aws_vpc.default_network.cidr_block
  vpc_id = aws_vpc.default_network.id

  tags = local.default_tags
}

resource "aws_route_table_association" "a" {
  subnet_id      = aws_subnet.default_subnet.id
  route_table_id = aws_route_table.default_route_table.id
}

resource "alertlogic_deployment" "dc_deployment" {
  account_id    = var.al_account_id
  name          = "dgreening-terraform-1"
  platform_type = "datacenter"
  mode          = "manual"
  lifecycle {
    ignore_changes = ["scope_include"]
  }
}

module scan-target-ec2 {
  source = "./modules/scan-target"
  aws_region = var.aws_region
  subnet_id = aws_subnet.default_subnet.id
  vpc_id = aws_vpc.default_network.id
  target_number = "2"
}

module ids-manual-ec2 {
  source = "./modules/ids-manual-ec2"
  account_id = var.al_account_id
  aws_region = var.aws_region
  deployment_id = alertlogic_deployment.dc_deployment.id
  subnet_id = aws_subnet.default_subnet.id
  vpc_id = aws_vpc.default_network.id
  network_claim_key = "e0785ef996bd2f6d1f4a16d04f18b566a777f4e70014ac23d8"
  appliance_number = "1"
}
