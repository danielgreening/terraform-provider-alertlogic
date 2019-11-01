Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

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

Provider Installation
----------------------

To use a custom-built provider in your Terraform environment (e.g. the provider binary from the build instructions above), follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

If you are running the `darwin_amd64` architecture, a pre-compiled plugin is provided with each release on the [Releases page](https://algithub.pd.alertlogic.net/daniel-greening/terraform-provider-alertlogic/releases). 

Provider Documentation
---------------------------------
The provider documentation is presented as an add-on to the terraform website.  

To deploy the docs locally you must perform the following.

```bash
# clone the terraform website repository into your GOPATH:
git clone clone https://github.com/hashicorp/terraform-website $GOPATH/src/github.com/hashicorp/terraform-website

# export this provider repo name
export PROVIDER_REPO=alertlogic

# link the provider repo
pushd "$GOPATH/src/github.com/hashicorp/terraform-website/ext/providers"
ln -sf "$GOPATH/src/algithub.alertlogic.pd.net/terraform-provider-$PROVIDER_REPO" "$PROVIDER_REPO"
popd

# link the layout file
pushd "$GOPATH/src/github.com/hashicorp/terraform-website/content/source/layouts"
ln -sf "../../../ext/providers/$PROVIDER_REPO/website/$PROVIDER_REPO.erb" "$PROVIDER_REPO.erb"
popd

# link the content
pushd "$GOPATH/src/github.com/hashicorp/terraform-website/content/source/docs/providers"
ln -sf "../../../../ext/providers/$PROVIDER_REPO/website/docs" "$PROVIDER_REPO"
popd

# start middleman
cd "$GOPATH/src/github.com/terraform-provider-$PROVIDER_REPO"
make website

```


$(GOPATH)/src/$(WEBSITE_REPO

You can run this locally by running `make website`. [Docker](https://www.docker.com/) is required and the steps outlined in the [terraform website repo](https://github.com/hashicorp/terraform-website#new-provider-repositories) must to be followed before the docs can be created.

Once the site is up and running, the documentation the provider specific configuration options can be found on the [provider's doc page](https://localhost:4567/docs/providers/alertlogic/index.html). 

Testing the Provider
---------------------------

In order to test the provider, you can run `make test`.

*Note:* Make sure no `ALERTLOGIC_ACCESS_KEY_ID` or `ALERTLOGIC_SECRET_ACCESS_KEY` variables are set. For more info on setting up access key authentication, see the [AIMS API docs](https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Resources-CreateAccessKey)
```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run. Please read [Running an Acceptance Test](https://github.com/terraform-providers/terraform-provider-aws/blob/master/.github/CONTRIBUTING.md#running-an-acceptance-test) in the contribution guidelines for more information on usage.

```sh
$ make testacc
```
