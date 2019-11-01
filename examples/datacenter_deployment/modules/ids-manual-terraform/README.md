AWS Terraform Template for Deploying Cloud Insight Appliance into Existing Subnet
=================================================================================
Manual Mode Deployment of Cloud Insight via Terraform

Use Case
--------
Customer who wishes to integrate the deployment of Cloud Insight security appliance into an existing subnet using Terraform instead of CFT.

Architecture Overview
---------------------
Alert Logic recommends deploying Cloud Insight scanners in a dedicated subnet created by our automation and currently the default behavior when deploying Cloud Insight. The Terraform scripts in this article are currently provided as a proof-of-concept of deploying Cloud Insight manually and is a workaround solution for supporting the appliance to be deployed into an existing subnet. Please contact the Deployment Architecture team before using these templates for instructions on how to setup the Cloud Insight deployment into manual mode.

The resources created by the main template will reuse an existing subnet (public or private) within the protected VPC and won't make any attempts to create:

  * New Route Table
  * New Security Subnet
  * New Alert Logic NACL

But will create the following:

  * New Security Group
  * New Launch Configuration
  * New Auto Scaling Group to Launch the Appliance from a shared AMI

It is the customer's responsibility to properly configure their network access attached to the VPC and subnet (NACL, IGW, NAT, etc.).

Step-by-step guide
------------------
> These steps **only** apply to the deployment portion of the appliance into an existing subnet. This article assumes that you've set the Cloud Insight deployment into manual mode and applied the read-only IAM policy already. Please contact the Deployment Architecture team to assist with prior configuration of the Cloud Insight deployment.

1. Download the main template and other associated files (*outputs.tf*, *userdata.tpl*, *variables.tf*) inside a working directory, i.e. *~/aws-ci-security-appliance-manual*
2. Assign values to variables by declaring them under the *var_values.tfvars* file.
3. Set up the provider access for AWS. Credentials can be provided from a separate file (default file name is *credentials.tf*). Variables can also be loaded from separate file or passed as environment variables. A full list of access methods are available here: https://www.terraform.io/docs/providers/aws/index.html
4. Run terraform apply
```
$ terraform apply -var-file var_values.tfvars
```
