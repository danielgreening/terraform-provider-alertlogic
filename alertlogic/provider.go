package alertlogic

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	homedir "github.com/mitchellh/go-homedir"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALERTLOGIC_ACCESS_KEY", ""),
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALERTLOGIC_SECRET_KEY", ""),
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALERTLOGIC_ENDPOINT", ""),
			},
			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ALERTLOGIC_SHARED_CREDENTIALS_FILE", ""),
			},
			"max_retries": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"alertlogic_deployment": resourceAlertLogicDeployment(),
			"alertlogic_credential": resourceAlertlogicCredential(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessKey:  d.Get("access_key_id").(string),
		SecretKey:  d.Get("secret_access_key").(string),
		Endpoint:   d.Get("endpoint").(string),
		MaxRetries: d.Get("max_retries").(int),
		Insecure:   d.Get("insecure").(bool),
	}

	// Set CredsFilename, expanding home directory
	credsPath, err := homedir.Expand(d.Get("shared_credentials_file").(string))
	if err != nil {
		return nil, err
	}

	config.CredsFilename = credsPath

	return config.Client()
}
