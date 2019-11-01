package alertlogic

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/alertlogic"
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/service/deployments"
	"algithub.pd.alertlogic.net/daniel-greening/terraform-provider-alertlogic/al-client-base"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlertLogicDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlertLogicDeploymentCreate,
		Read:   resourceAlertLogicDeploymentRead,
		Update: resourceAlertLogicDeploymentUpdate,
		Delete: resourceAlertLogicDeploymentDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				accountId, deploymentId, err := resourceAlertLogicDeploymentParseImportId(d.Id())

				if err != nil {
					return nil, err
				}

				d.Set("account_id", accountId)
				d.SetId(deploymentId)

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
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"platform_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}
					return false
				},
			},
			"scope_include": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"policy_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"scope_exclude": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_defender_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"cloud_defender_location_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"credential": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"purpose": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlertLogicDeploymentParseImportId(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected account_id:deployment_id", id)
	}

	return parts[0], parts[1], nil
}

func resourceAlertLogicDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	deploymentsconn := meta.(*ALClient).deploymentsconn

	var accountId string
	if confAccountId, ok := d.GetOk("account_id"); ok {
		accountId = confAccountId.(string)
	} else {
		accountId = meta.(*ALClient).accountid
	}

	if accountId == "" {
		return errors.New("account ID not set in deployment configuration or found in session details")
	}

	platform := &deployments.Platform{
		Type: alertlogic.String(d.Get("platform_type").(string)),
	}

	if platformId, ok := d.GetOk("platform_id"); ok {
		platform.Id = alertlogic.String(platformId.(string))
	}

	scope := &deployments.Scope{}

	if scopeInclude, ok := d.GetOk("scope_include"); ok {
		scope.Include = expandDeploymentScope(scopeInclude.(*schema.Set))
	}

	if scopeExclude, ok := d.GetOk("scope_exclude"); ok {
		scope.Exclude = expandDeploymentScope(scopeExclude.(*schema.Set))
	}

	params := &deployments.CreateDeploymentInput{
		AccountId: alertlogic.String(accountId),
		Name:      alertlogic.String(d.Get("name").(string)),
		Mode:      alertlogic.String(d.Get("mode").(string)),
		Platform:  platform,
		Scope:     scope,
	}

	if credentials, ok := d.GetOk("credential"); ok {
		params.Credentials = expandDeploymentCredentials(credentials.(*schema.Set))
	}

	cloudDefender := &deployments.CloudDefender{}

	if cloudDefenderEnabled, ok := d.GetOk("cloud_defender_enabled"); ok {
		cloudDefender.Enabled = alertlogic.Bool(cloudDefenderEnabled.(bool))
	}

	if cloudDefenderLocationId, ok := d.GetOk("cloud_defender_location_id"); ok {
		cloudDefender.LocationId = alertlogic.String(cloudDefenderLocationId.(string))
	}

	if cloudDefender.Enabled != nil || cloudDefender.LocationId != nil {
		params.CloudDefender = cloudDefender
	}

	log.Printf("[DEBUG] Creating Deployment: %#v", params)

	err := resource.Retry(time.Duration(1)*time.Minute, func() *resource.RetryError {
		var err error
		log.Println("Initial error: ", err)
		resp, err := deploymentsconn.CreateDeployment(params)
		if err != nil {
			if al_client_base.IsALErr(err, "400", "Bad Request") {
				log.Printf("[DEBUG] Received server err, retrying: %s", err)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.SetId(*resp.Deployment.Id)
		return nil
	})

	if err != nil {
		return fmt.Errorf("error creating deployment: %s", err)
	}

	return resourceAlertLogicDeploymentRead(d, meta)
}

func resourceAlertLogicDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	deploymentsconn := meta.(*ALClient).deploymentsconn

	var accountId string
	if confAccountId, ok := d.GetOk("account_id"); ok {
		accountId = confAccountId.(string)
	} else {
		accountId = meta.(*ALClient).accountid
	}

	params := &deployments.GetDeploymentInput{
		AccountId: alertlogic.String(accountId),
		Id:        alertlogic.String(d.Id()),
	}

	return resource.Retry(time.Duration(1)*time.Minute, func() *resource.RetryError {
		resp, err := deploymentsconn.GetDeployment(params)

		if err != nil {
			if al_client_base.IsALErr(err, "404", "Not Found") {
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("error retrieving deployment: %s", err))
		}
		log.Println("[DEBUG] Account received from deployment: ", alertlogic.StringValue(resp.Deployment.AccountId))

		d.Set("account_id", resp.Deployment.AccountId)
		d.Set("name", resp.Deployment.Name)
		d.Set("version", resp.Deployment.Version)
		d.Set("mode", resp.Deployment.Mode)
		d.Set("platform_type", resp.Deployment.Platform.Type)
		d.Set("platform_id", resp.Deployment.Platform.Id)

		d.Set("scope_include", flattenDeploymentScope(resp.Deployment.Scope.Include))
		d.Set("scope_exclude", flattenDeploymentScope(resp.Deployment.Scope.Exclude))

		d.Set("credential", flattenDeploymentCredentials(resp.Deployment.Credentials))

		d.Set("status", resp.Deployment.Status.Status)
		d.Set("status_message", resp.Deployment.Status.Message)
		d.Set("status_updated", resp.Deployment.Status.Updated.String())

		d.Set("cloud_defender_enabled", resp.Deployment.CloudDefender.Enabled)
		d.Set("cloud_defender_location_id", resp.Deployment.CloudDefender.LocationId)

		return nil
	})
}

func resourceAlertLogicDeploymentUpdate(d *schema.ResourceData, meta interface{}) error {
	deploymentsconn := meta.(*ALClient).deploymentsconn

	var accountId string
	if confAccountId, ok := d.GetOk("account_id"); ok {
		accountId = confAccountId.(string)
	} else {
		accountId = meta.(*ALClient).accountid
	}

	platform := &deployments.Platform{}

	if platformId, ok := d.GetOk("platform_id"); ok {
		platform.Id = alertlogic.String(platformId.(string))
	}

	scope := &deployments.Scope{}

	if scopeInclude, ok := d.GetOk("scope_include"); ok {
		scope.Include = expandDeploymentScope(scopeInclude.(*schema.Set))
	}

	if scopeExclude, ok := d.GetOk("scope_exclude"); ok {
		scope.Exclude = expandDeploymentScope(scopeExclude.(*schema.Set))
	}

	params := &deployments.UpdateDeploymentInput{
		AccountId: alertlogic.String(accountId),
		Id:        alertlogic.String(d.Id()),
		Version:   alertlogic.Int64(int64(d.Get("version").(int))),
		Name:      alertlogic.String(d.Get("name").(string)),
		Mode:      alertlogic.String(d.Get("mode").(string)),
		Platform:  platform,
		Scope:     scope,
	}

	if credentials, ok := d.GetOk("credential"); ok {
		params.Credentials = expandDeploymentCredentials(credentials.(*schema.Set))
	}

	log.Printf("[DEBUG] Updating Deployment: %#v", params)

	err := resource.Retry(time.Duration(1)*time.Minute, func() *resource.RetryError {
		_, err := deploymentsconn.UpdateDeployment(params)

		if err != nil {
			if al_client_base.IsALErr(err, "500", "") {
				log.Printf("[DEBUG] Received server err, retrying: %s", err)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error updating deployment: %s", err)
	}

	return resourceAlertLogicDeploymentRead(d, meta)
}

func resourceAlertLogicDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	deploymentsconn := meta.(*ALClient).deploymentsconn

	var accountId string
	if confAccountId, ok := d.GetOk("account_id"); ok {
		accountId = confAccountId.(string)
	} else {
		accountId = meta.(*ALClient).accountid
	}

	log.Printf("[INFO] Deleting Deployment: %s", d.Id())

	params := &deployments.DeleteDeploymentInput{
		AccountId: alertlogic.String(accountId),
		Id:        alertlogic.String(d.Id()),
	}

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		_, err := deploymentsconn.DeleteDeployment(params)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil && !al_client_base.IsALErr(err, "404", "Not Found") {
		return fmt.Errorf("error deleting deployment: %s", err)
	}

	return nil
}
