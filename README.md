Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin from source)
- [Docker](https://www.docker.com/products/docker-desktop) 19.0.3+ (to build provider documentation)

Building Documentation
---------------------------------
The provider documentation is presented as an add-on to the terraform website.  

You can preview the website from a local checkout of this repo as follows:

1. Install [Docker](https://docs.docker.com/install/) if you have not already done so.
2. Go to the top directory of this repo in your terminal, and run `make website`.
3. Open `http://localhost:4567` in your web browser.
4. When you're done with the preview, press ctrl-C in your terminal to stop the server.

Once the site is up and running, the documentation the provider specific configuration options can be found on the [provider's doc page](https://localhost:4567/docs/providers/alertlogic/index.html). 

Provider Installation
----------------------

To use a custom-built provider in your Terraform environment (e.g. the provider binary from the build instructions above), follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

If you are running the `darwin_amd64` architecture, a pre-compiled plugin is provided with each release on the [Releases page](https://algithub.pd.alertlogic.net/daniel-greening/terraform-provider-alertlogic/releases). You can begin using the plugin by simply by placing the binary in the third-party plugins directory, see the [Terraform docs](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) for more information on third-party plugins. 

Configuring the Provider
----------------------

Once installed, you must provide a means of configuring the provider before you may start creating resources. 

The following configuration options are supported:

- Endpoint - the Alert Logic API endpoint to call when creating resources
- Access Key ID - the ID of the Alert Logic access key ID  
- Secret Access Key - the secret key associated with the given Alert Logic access key  

Currently, only AIMS access keys are supported as a means of authentication. See the [AIMS API docs](https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Resources-CreateAccessKey) for information on creating access keys for your Alert Logic User.

Once you have an access key and corresponding secret key set up, there are two mechanisms for passing them to the provider: 

### Provider argument
The provider can be configured directly in Terraform HCL:

Example Terraform HCL configuration:
```hcl
provider alertlogic {
  access_key_id = "your-access-key"
  secret_access_key = "your-secret-key"
  endpoint = "https://api.cloudinsight.alertlogic.com" // US Production
}
```

### Environment Variables
If no arguments are found for the provider in the HCL, then environment variables can be used instead:

Example Terraform HCL configuration
```hcl
provider alertlogic {}
```

Corresponding variables:
```bash
export ALERTLOGIC_ACCESS_KEY_ID="your-access-key"
export ALERTLOGIC_SECRET_ACCESS_KEY="your-secret-key"
export ALERTLOGIC_ENDPOINT="https://api.cloudinsight.alertlogic.com"
```

Further information on configuring and workable examples can be found in the [Provider Documentation](#Building-Documentation)

Developing the Provider
---------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please check the [requirements](#Requirements) before proceeding).

*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (i.e `$HOME/development/terraform-providers/`).

Clone repository to: `$HOME/development/terraform-providers/`

```sh
$ mkdir -p $HOME/development/terraform-providers/; cd $HOME/development/terraform-providers/
$ git clone git@algithub.pd.alertlogic.net:daniel-greening/terraform-provider-alertlogic
...
```

Enter the provider directory and run `make tools`. This will install the needed tools for the provider.

```sh
$ make tools
```

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-alertlogic
...
```

Testing the Provider
---------------------------

In order to test the provider, you can run `make test`.

*Note:* Make sure no `ALERTLOGIC_ACCESS_KEY_ID` or `ALERTLOGIC_SECRET_ACCESS_KEY` variables are set. For more info on setting up access key authentication, see the [AIMS API docs](https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Resources-CreateAccessKey)
```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* `ALERTLOGIC_ACCESS_KEY_ID` or `ALERTLOGIC_SECRET_ACCESS_KEY` must be set for acceptance testing. For more info on setting up access key authentication, see the [AIMS API docs](https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Resources-CreateAccessKey)

```sh
$ make testacc
```
