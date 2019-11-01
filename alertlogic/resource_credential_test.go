package alertlogic

import (
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/alertlogic"
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/service/credentials"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccAlertLogicAzureCredential_basic(t *testing.T) {
	var v credentials.Credential

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccDefaultAzureCredential,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCredentialExists("alertlogic_credential.foo", &v),
					resource.TestCheckResourceAttr(
						"alertlogic_credential.foo",
						"name",
						"foo_credential",
					),
					resource.TestCheckResourceAttr(
						"alertlogic_credential.foo",
						"secret_type",
						"azure_ad_client",
					),
					resource.TestCheckResourceAttr(
						"alertlogic_credential.foo",
						"secret_ad_id",
						"test_ad_id",
					),
					resource.TestCheckResourceAttr(
						"alertlogic_credential.foo",
						"secret_client_id",
						"test_client_id",
					),
					resource.TestCheckResourceAttr(
						"alertlogic_credential.foo",
						"secret_client_secret",
						"test_client_secret",
					),
				),
			},
		},
	})
}

const testAccDefaultAzureCredential = `
resource "alertlogic_credential" "foo" {
 name                 = "foo_credential"
 secret_type          = "azure_ad_client"
 secret_ad_id         = "test_ad_id" 
 secret_client_id     = "test_client_id"
 secret_client_secret = "test_client_secret"
}
`

func testAccCheckCredentialExists(n string, v *credentials.Credential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[ERROR] Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("[ERROR] No ID is set")
		}

		accId, ok := rs.Primary.Attributes["account_id"]

		if !ok {
			return fmt.Errorf("[ERROR] No Account ID is set")
		}

		conn := testAccProvider.Meta().(*ALClient).credentialsconn
		resp, err := conn.GetCredential(&credentials.GetCredentialInput{
			AccountId: alertlogic.String(accId),
			Id:        alertlogic.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		*v = *resp.Credential

		return nil
	}
}
