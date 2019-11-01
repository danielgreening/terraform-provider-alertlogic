package alertlogic

import (
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/service/credentials"
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/service/deployments"
	"github.com/hashicorp/terraform/helper/logging"
	"log"

	albase "algithub.pd.alertlogic.net/daniel-greening/terraform-provider-alertlogic/al-client-base"
	"github.com/hashicorp/terraform/terraform"
)

type Config struct {
	AccessKey     string
	SecretKey     string
	CredsFilename string
	Profile       string
	Token         string
	Region        string
	MaxRetries    int

	Endpoint string
	Insecure bool

	SkipCredsValidation bool
}

type ALClient struct {
	accountid       string
	deploymentsconn *deployments.Deployments
	credentialsconn *credentials.Credentials
}

// Client configures and returns a fully initialized ALClient
func (c *Config) Client() (interface{}, error) {
	log.Println("[INFO] Building Alert Logic auth structure")
	albaseConfig := &albase.Config{
		AccessKey:           c.AccessKey,
		CredsFilename:       c.CredsFilename,
		DebugLogging:        logging.IsDebugOrHigher(),
		Insecure:            c.Insecure,
		MaxRetries:          c.MaxRetries,
		SecretKey:           c.SecretKey,
		SkipCredsValidation: c.SkipCredsValidation,
		YarpEndpoint:        c.Endpoint,
		Token:               c.Token,
		UserAgentProducts: []*albase.UserAgentProduct{
			{Name: "APN", Version: "1.0"},
			{Name: "AlertLogic", Version: "v0.0.1"},
			{Name: "Terraform", Version: terraform.VersionString()},
		},
	}

	sess, accountID, err := albase.GetSessionWithAccountID(albaseConfig)
	if err != nil {
		return nil, err
	}

	if accountID == "" {
		log.Printf("[WARN] Alert Logic account ID not found for alertlogic.")
	}

	client := &ALClient{
		accountid:       accountID,
		deploymentsconn: deployments.New(sess.Copy()),
		credentialsconn: credentials.New(sess.Copy()),
	}

	return client, nil
}
