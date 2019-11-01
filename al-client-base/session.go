package al_client_base

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"log"
	"net/http"

	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/alertlogic"
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/alertlogic/credentials/aimscreds"
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/alertlogic/request"
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/alertlogic/session"
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/service/aims"
)

// GetSessionOptions attempts to return valid AlertLogic Go SDK session authentication
// options based on pre-existing credential provider, configured profile, or
// fallback to automatically a determined session via the ALertLogic Go SDK.
func GetSessionOptions(c *Config) (*session.Options, error) {
	options := &session.Options{
		Config: alertlogic.Config{
			HTTPClient: cleanhttp.DefaultClient(),
			MaxRetries: alertlogic.Int(c.MaxRetries),
			Endpoint:   alertlogic.String(c.YarpEndpoint),
		},
	}

	creds, err := GetCredentials(c)
	if err != nil {
		return nil, err
	}

	// Call Get to check for credential provider. If nothing found, we'll get an
	// error, and we can present it nicely to the user
	cp, err := creds.Get()
	if err != nil {
		if IsALErr(err, "NoCredentialProviders", "") {
			// If a profile wasn't specified, the session may still be able to resolve credentials from shared config.
			if c.Profile == "" {
				sess, err := session.NewSession()
				if err != nil {
					return nil, errors.New(`no valid credential sources found for AL Provider`)
				}
				_, err = sess.Config.Credentials.Get()
				if err != nil {
					return nil, errors.New(`no valid credential sources found for AL Provider`)
				}
				log.Printf("[INFO] Using session-derived AL Auth")
				options.Config.Credentials = sess.Config.Credentials
			}
		} else {
			return nil, fmt.Errorf("error loading credentials for AL Provider: %s", err)
		}
	} else {
		// add the validated credentials to the session options
		log.Printf("[INFO] AL Auth provider used: %q", cp.ProviderName)
		options.Config.Credentials = creds
	}

	if c.Insecure {
		transport := options.Config.HTTPClient.Transport.(*http.Transport)
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	if c.DebugLogging {
		options.Config.LogLevel = alertlogic.LogLevel(alertlogic.LogDebugWithHTTPBody | alertlogic.LogDebugWithRequestRetries | alertlogic.LogDebugWithRequestErrors)
		options.Config.Logger = DebugLogger{}
	}

	return options, nil
}

// GetSession attempts to return valid AWS Go SDK session
func GetSession(c *Config) (*session.Session, error) {
	options, err := GetSessionOptions(c)

	if err != nil {
		return nil, err
	}

	sess, err := session.NewSessionWithOptions(*options)
	if err != nil {
		if IsALErr(err, "NoCredentialProviders", "") {
			return nil, errors.New(`no valid credential sources found for AL Provider`)
		}
		return nil, fmt.Errorf("error creating AIMS session: %s", err)
	}

	creds := aimscreds.NewCredentials(sess)
	sess = sess.Copy(&alertlogic.Config{Credentials: creds})

	if c.MaxRetries > 0 {
		sess = sess.Copy(&alertlogic.Config{MaxRetries: alertlogic.Int(c.MaxRetries)})
	}

	//for _, product := range c.UserAgentProducts {
	//	sess.Handlers.Build.PushBack(request.MakeAddToUserAgentHandler(product.Name, product.Version, product.Extra...))
	//}

	// Generally, we want to configure a lower retry theshold for networking issues
	// as the session retry threshold is very high by default and can mask permanent
	// networking failures, such as a non-existent service endpoint.
	// MaxRetries will override this logic if it has a lower retry threshold.
	// NOTE: This logic can be fooled by other request errors raising the retry count
	//       before any networking error occurs
	sess.Handlers.Retry.PushBack(func(r *request.Request) {
		// We currently depend on the DefaultRetryer exponential backoff here.
		// ~10 retries gives a fair backoff of a few seconds.
		if r.RetryCount < 9 {
			return
		}
		// RequestError: send request failed
		// caused by: Post https://FQDN/: dial tcp: lookup FQDN: no such host
		if IsALErrExtended(r.Error, "RequestError", "send request failed", "no such host") {
			log.Printf("[WARN] Disabling retries after next request due to networking issue")
			r.Retryable = alertlogic.Bool(false)
		}
		// RequestError: send request failed
		// caused by: Post https://FQDN/: dial tcp IPADDRESS:443: connect: connection refused
		if IsALErrExtended(r.Error, "RequestError", "send request failed", "connection refused") {
			log.Printf("[WARN] Disabling retries after next request due to networking issue")
			r.Retryable = alertlogic.Bool(false)
		}
	})

	if !c.SkipCredsValidation {
		aimsClient := aims.New(sess.Copy(&alertlogic.Config{Endpoint: alertlogic.String(c.YarpEndpoint)}))
		if _, err := GetAccountIDAFromAIMSGetTokenInfo(aimsClient); err != nil {
			return nil, fmt.Errorf("error validating provider credentials: %s", err)
		}
	}

	return sess, nil
}

// GetSessionWithAccountID attempts to return valid AlertLogic Go SDK session
// along with account ID information if available
func GetSessionWithAccountID(c *Config) (*session.Session, string, error) {
	sess, err := GetSession(c)

	if err != nil {
		return nil, "", err
	}

	aimsClient := aims.New(sess.Copy(&alertlogic.Config{Endpoint: alertlogic.String(c.YarpEndpoint)}))

	accountID, err := GetAccountID(aimsClient)

	if err == nil {
		return sess, accountID, nil
	}

	return nil, "", fmt.Errorf(
		"AlertLogic account ID not previously found and failed retrieving via all available methods. Errors: %s",
		err,
	)
}
