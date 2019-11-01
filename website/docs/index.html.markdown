---
layout: "alertlogic"
page_title: "Provider: Alert Logic"
---

# Alert Logic Provider

The Alert Logic provider is used to interact with the
many resources supported by Alert Logic. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Alert Logic Provider
provider "alertlogic" {
  endpoint  = "https://api.cloudinsight.alertlogic.com"
}

# Create a Deployment
resource "alertlogic_deployment" "example" {
  name          = "terraform-deployment"
  platform_type = "datacenter"
  mode          = "manual"
}
```

## Authentication

The Alert Logic provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Static credentials
- Environment variables

### Static credentials ###

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding an `access_key_id` and `secret_access_key`
in-line in the AWS provider block:

Usage:

```hcl
provider "alertlogic" {
  access_key_id = "my-access-key"
  secret_access_key = "my-secret-key"
}
```

### Environment variables

You can provide your credentials via the `ALERTLOGIC_ACCESS_KEY_ID` and
`ALERTLOGIC_SECRET_ACCESS_KEY`, environment variables, representing your Alert Logic
Access Key and Alert Logic Secret Key, respectively.  Note that setting your
Alert Logic credentials using either these environment variables
will override the use of `ALERTLOGIC_SHARED_CREDENTIALS_FILE` and `ALERTLOGIC_PROFILE`.
The `ALERTLOGIC_ENDPOINT` and `ALERTLOGIC_SESSION_TOKEN` environment variables
are also used, if applicable:

```hcl
provider "alertlogic" {}
```

Usage:

```sh
$ export ALERTLOGIC_ACCESS_KEY_ID="anaccesskey"
$ export ALERTLOGIC_SECRET_ACCESS_KEY="asecretkey"
$ export ALERTLOGIC_ENDPOINT="https://api.cloudinsight.alertlogic.com"
$ terraform plan
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the Alert Logic
 `provider` block:

* `access_key_id` - (Optional) This is the Alert Logic access key. It must be provided, but
  it can also be sourced from the `ALERTLOGIC_ACCESS_KEY_ID` environment variable, or via
  a shared credentials file if `profile` is specified.

* `secret_access_key` - (Optional) This is the Alert Logic secret key. It must be provided, but
  it can also be sourced from the `ALERTLOGIC_SECRET_ACCESS_KEY` environment variable, or
  via a shared credentials file if `profile` is specified.

* `endpoint` - (Optional) This is the Alert Logic API endpoint. It must be provided, but
  it can also be sourced from the `ALERTLOGIC_ENDPOINT` environment variable.

* `max_retries` - (Optional) This is the maximum number of times an API
  call is retried, in the case where requests are being throttled or
  experiencing transient failures. The delay between the subsequent API
  calls increases exponentially.
