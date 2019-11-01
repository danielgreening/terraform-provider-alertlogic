package main

import (
	"algithub.pd.alertlogic.net/daniel-greening/terraform-provider-alertlogic/alertlogic"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: alertlogic.Provider,
	})
}
