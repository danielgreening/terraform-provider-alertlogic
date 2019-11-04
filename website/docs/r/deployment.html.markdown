---
layout: "alertlogic"
page_title: "Alert Logic: alertlogic_deployment"
description: |-
  Provides an Alert Logic Deployment.
---

# Resource: alertlogic_deployment

Provides an Alert Logic Deployment.

## Example Usage

```hcl
resource "alertlogic_deployment" "test_deployment" {
  account_id  	= "123456"
  name        	= "test-aws-deployment"
  provider_type = "aws" 
  provider_id 	= "7654321" 
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional, Forces new resource) The ID of the Alert Logic account to create the credential under.
* `name` - (Required, Forces new resource) The name of the credential.
* `provider_type` - (Required, Forces new resource) The type of deployment to create, currently supported are `aws`, `azure` or `datacenter`
* `provider_id` - (Optional) The Identifier of the external account associated with the deployment i.e. AWS account ID, or Azure subscription ID. Not required for Datacenter deployments
* `mode` - (Required) - The installation mode the deployment should use, `manual` is supported for all providers, while `aws` deployments may be `automatic`
* `scope_include` - (Optional) Configuration block for declaring an in-scope object
* `scope_exclude` - (Optional) Configuration block for declaring an out-of-scope object
* `credential` - (Optional) Configuration block for associating a credential with the deployment

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The stable and unique string identifying the credential.
* `external_id` - The External ID associated with the provided AWS IAM Role, inferred by the Credentials service.


## Example of an AWS Deployment with a discovery credential

```hcl
resource "alertlogic_deployment" "test_deployment" {
  account_id  = "123456"
  name        = "test-aws-deployment"
  platform_id = "platform_id"
  secret_arn  = "arn:aws:iam::123456789012:root"
  
  credential {
  	type = "discover"
  	id = alertlogic_credential.test.id
  }
}

resource "alertlogic_credential" "discovery_credential" {
  account_id  = "123456" 
  name		  = "discovery_credential"
  secret_type = "aws_iam_role"
  secret_arn  = "arn:aws:iam::123456789012:root"
}
```

## Example of an Azure Deployment with a VPC (VNET) in Professional scope

```hcl
resource "alertlogic_deployment" "test_azure_client_credential" {
  account_id  = "123456"
  name        = "test-aws-deployment"
  platform_id = "platform_id"
  secret_arn  = "arn:aws:iam::123456789012:root"
  
  scope_include {
  	type = "vpc"
  	key =  "/azure/region-name/vpc/vnet-id"
  	policy_id = "D12D5E67-166C-474F-87AA-6F86FC9FB9BC"
  }
}
```

## Import

Existing deployments can be imported using a composite of the Alert Logic Account ID and Deployment ID, separated by a `:` e.g.

```
$ terraform import alertlogic_deployment.test 123456:ABCDE-FGHIKLMNOPQR-STUV
```
