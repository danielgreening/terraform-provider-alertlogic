package alertlogic

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"alertlogic": testAccProvider,
	}

}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("ALERTLOGIC_PROFILE") == "" && os.Getenv("ALERTLOGIC_ACCESS_KEY_ID") == "" {
		t.Fatal("ALERTLOGIC_ACCESS_KEY_ID or ALERTLOGIC_PROFILE must be set for acceptance tests")
	}

	if os.Getenv("ALERTLOGIC_ACCESS_KEY_ID") != "" && os.Getenv("ALERTLOGIC_SECRET_ACCESS_KEY") == "" {
		t.Fatal("ALERTLOGIC_SECRET_ACCESS_KEY must be set for acceptance tests")
	}

	endpoint := testAccGetEndpoint()
	log.Printf("[INFO] Test: Using %s as test endpoint", endpoint)
	os.Setenv("ALERTLOGIC_ENDPOINT", endpoint)

	err := testAccProvider.Configure(terraform.NewResourceConfig(nil))
	if err != nil {
		t.Fatal(err)
	}
}

// testAccAwsProviderAccountID returns the account ID of an AWS provider
func testAccAlertLogicProviderAccountID(provider *schema.Provider) string {
	if provider == nil {
		log.Print("[DEBUG] Unable to read account ID from test provider: empty provider")
		return ""
	}
	if provider.Meta() == nil {
		log.Print("[DEBUG] Unable to read account ID from test provider: unconfigured provider")
		return ""
	}
	client, ok := provider.Meta().(*ALClient)
	if !ok {
		log.Print("[DEBUG] Unable to read account ID from test provider: non-AWS or unconfigured AWS provider")
		return ""
	}
	return client.accountid
}

// testAccCheckResourceAttrAccountID ensures the Terraform state exactly matches the account ID
func testAccCheckResourceAttrAccountID(resourceName, attributeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return resource.TestCheckResourceAttr(resourceName, attributeName, testAccGetAccountID())(s)
	}
}

// testAccGetAccountID returns the account ID of testAccProvider
// Must be used returned within a resource.TestCheckFunc
func testAccGetAccountID() string {
	return testAccAlertLogicProviderAccountID(testAccProvider)
}

func testAccGetEndpoint() string {
	v := os.Getenv("ALERTLOGIC_ENDPOINT")
	if v == "" {
		return "https://api.cloudinsight.alertlogic.com"
	}
	return v
}
