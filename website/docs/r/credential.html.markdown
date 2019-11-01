---
layout: "alertlogic"
page_title: "Alert Logic: alertlogic_credential"
description: |-
  Provides an Alert Logic Credential.
---

# Resource: alertlogic_credential

Provides an Alert Logic Credential.

## Example Usage

```hcl
resource "alertlogic_credential" "test_role" {
  account_id  = "123456"
  name        = "test-aws-iam-role"
  secret_type = "aws_iam_role"
  secret_arn  = "arn:aws:iam::123456789012:root"
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional, Forces new resource) The ID of the Alert Logic account to create the credential under.
* `name` - (Required, Forces new resource) The name of the credential.
* `secret_type` (Required, Forces new resource) Type of credential to be stored, e.g. `aws_iam_role`, See [Credentials API Documentation](https://console.cloudinsight.alertlogic.com/api/credentials/index.html#api-Management-CreateCredential) for full list of supported credential types.

* `secret_arn` (Optional, Forces new resource) The ARN of the AWS IAM role, **Note** required for credentials of type `aws_iam_role`

* `secret_ad_id` (Optional, Forces new resource) The Active Directory Identifier of the credential, **Note** required for credentials of type `azure_ad_client_`
* `secret_client_id` (Optional, Forces new resource) The Client Identifier of the credential, **Note** required for credentials of type `azure_ad_client`
* `secret_client_secret` (Optional, Forces new resource) The Secret password associated with the credential, **Note** required for credentials of type `azure_ad_client`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The stable and unique string identifying the credential.
* `external_id` - The External ID associated with the provided AWS IAM Role, inferred by the Credentials service.


## Example of AWS IAM Role Credential

```hcl
resource "alertlogic_credential" "test_role" {
  account_id  = "123456"
  name        = "test-aws-iam-role"
  secret_type = "aws_iam_role"
  secret_arn  = "arn:aws:iam::123456789012:root"
}
```

## Example of Azure Active Directory Client Credential

```hcl
resource "alertlogic_credential" "test_azure_client_credential" {
  account_id  			= "123456"
  name        			= "test-azure-credential"
  secret_type 			= "azure_ad_client"
  secret_ad_id  		= "active_directory_id"
  secret_client_id      = "test_client_app_id"
  secret_client_secret  = "strong_password"
}
```

## Import

Credentials can be imported using a composite of the Alert Logic Account ID and Credential ID, separated by a `:` e.g.

```
$ terraform import alertlogic_credential.test 123456:ABCDE-FGHIKLMNOPQR-STUV
```
