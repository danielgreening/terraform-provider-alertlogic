package alertlogic

import (
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/service/credentials"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/alertlogic"
	"algithub.pd.alertlogic.net/daniel-greening/terraform-provider-alertlogic/al-client-base"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlertlogicCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlertlogicCredentialCreate,
		Read:   resourceAlertlogicCredentialRead,
		Delete: resourceAlertlogicCredentialDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				accountId, credentialId, err := resourceAlertlogicCredentialParseImportId(d.Id())

				if err != nil {
					return nil, err
				}

				d.Set("account_id", accountId)
				d.SetId(credentialId)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secret_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secret_arn": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"secret_external_id": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"secret_ad_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secret_client_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secret_client_secret": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlertlogicCredentialParseImportId(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected account_id:credential_id", id)
	}

	return parts[0], parts[1], nil
}

func resourceAlertlogicCredentialCreate(d *schema.ResourceData, meta interface{}) error {
	credentialsconn := meta.(*ALClient).credentialsconn

	var accountId string
	if confAccountId, ok := d.GetOk("account_id"); ok {
		accountId = confAccountId.(string)
	} else {
		accountId = meta.(*ALClient).accountid
	}

	if accountId == "" {
		return errors.New("account ID not set in credential configuration or found in session details")
	}

	secrets := &credentials.Secrets{
		Type: alertlogic.String(d.Get("secret_type").(string)),
	}

	if secretArn := d.Get("secret_arn"); secretArn != "" {
		secrets.Arn = alertlogic.String(secretArn.(string))
	}

	if secretAdID := d.Get("secret_ad_id"); secretAdID != "" {
		secrets.AdID = alertlogic.String(secretAdID.(string))
	}

	if secretClientID := d.Get("secret_client_id"); secretClientID != "" {
		secrets.ClientID = alertlogic.String(secretClientID.(string))
	}

	if secretClientSecret := d.Get("secret_client_secret"); secretClientSecret != "" {
		secrets.ClientSecret = alertlogic.String(secretClientSecret.(string))
	}

	params := &credentials.CreateCredentialInput{
		AccountId: alertlogic.String(accountId),
		Name:      alertlogic.String(d.Get("name").(string)),
		Secrets:   secrets,
	}

	log.Printf("[DEBUG] Creating Credential: %#v", params)
	resp, err := credentialsconn.CreateCredential(params)

	if err != nil {
		return fmt.Errorf("error creating credential: %s\n resp: %s", err, resp)
	}

	d.SetId(*resp.Credential.Id)

	return resourceAlertlogicCredentialRead(d, meta)
}

func resourceAlertlogicCredentialRead(d *schema.ResourceData, meta interface{}) error {
	credentialsconn := meta.(*ALClient).credentialsconn

	var accountId string
	if confAccountId, ok := d.GetOk("account_id"); ok {
		accountId = confAccountId.(string)
	} else {
		accountId = meta.(*ALClient).accountid
	}

	params := &credentials.GetCredentialInput{
		AccountId: alertlogic.String(accountId),
		Id:        alertlogic.String(d.Id()),
	}

	return resource.Retry(time.Duration(1)*time.Minute, func() *resource.RetryError {
		resp, err := credentialsconn.GetCredential(params)

		if err != nil {
			if al_client_base.IsALErr(err, "404", "Not Found") {
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("error retrieving credential: %s", err))
		}

		d.Set("account_id", accountId)
		d.Set("name", resp.Credential.Name)
		d.Set("secret_type", resp.Credential.Secrets.Type)
		d.Set("secret_arn", resp.Credential.Secrets.Arn)
		d.Set("secret_external_id", resp.Credential.Secrets.ExternalId)
		d.Set("secret_ad_id", resp.Credential.Secrets.AdID)
		d.Set("secret_client_id", resp.Credential.Secrets.ClientID)

		return nil
	})
}

func resourceAlertlogicCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	credentialsconn := meta.(*ALClient).credentialsconn

	log.Printf("[INFO] Deleting Credential: %s", d.Id())

	params := &credentials.DeleteCredentialInput{
		AccountId: alertlogic.String(d.Get("account_id").(string)),
		Id:        alertlogic.String(d.Id()),
	}

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := credentialsconn.DeleteCredential(params)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil && !al_client_base.IsALErr(err, "404", "Not Found") {
		return fmt.Errorf("error deleting credential: %s", err)
	}

	return nil
}
