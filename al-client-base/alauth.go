package al_client_base

import (
	"errors"
	"fmt"
	"log"
	"time"

	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/alertlogic"
	alCredentials "algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/alertlogic/credentials"
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/service/aims"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-multierror"
)

func GetAccountID(aimsconn *aims.AIMS) (string, error) {
	var accountID string
	var err, errors error

	accountID, err = GetAccountIDAFromAIMSGetTokenInfo(aimsconn)
	if accountID != "" {
		return accountID, nil
	}
	errors = multierror.Append(errors, err)

	return accountID, errors
}

func GetAccountIDAFromAIMSGetTokenInfo(aimsconn *aims.AIMS) (string, error) {
	log.Println("[DEBUG] Trying to get account information via aims:GetTokenInfo")

	output, err := aimsconn.GetTokenInfo(&aims.GetTokenInfoInput{})
	if err != nil {
		return "", fmt.Errorf("error calling aims:GetTokenInfo: %s", err)
	}

	if output == nil || output.Account == nil {
		err = errors.New("empty aims:GetTokenInfo response")
		log.Printf("[DEBUG] %s", err)
		return "", err
	}

	return alertlogic.StringValue(output.Account.Id), nil
}

// This function is responsible for reading credentials from the
// environment in the case that they're not explicitly specified
// in the Terraform configuration.
func GetCredentials(c *Config) (*alCredentials.Credentials, error) {
	// build a chain provider, lazy-evaluated by al_client
	providers := []alCredentials.Provider{
		&alCredentials.StaticProvider{Value: alCredentials.Value{
			AccessKeyID: c.AccessKey,
			SecretKey:   c.SecretKey,
			AuthToken:   c.Token,
		}},
		&alCredentials.EnvProvider{},
		&alCredentials.SharedCredentialsProvider{
			Filename: c.CredsFilename,
		},
	}

	// Build isolated HTTP client to avoid issues with globally-shared settings
	client := cleanhttp.DefaultClient()

	client.Timeout = 100 * time.Millisecond

	return alCredentials.NewChainCredentials(providers), nil
}
