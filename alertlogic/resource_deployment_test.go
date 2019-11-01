package alertlogic

import (
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/alertlogic"
	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/service/deployments"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccAlertLogicDeployment_basic(t *testing.T) {
	var v deployments.Deployment

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccDefaultAWSDeployment,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDeploymentExists("alertlogic_deployment.foo", &v),
					resource.TestCheckResourceAttr(
						"alertlogic_deployment.foo",
						"name",
						"terraform-testacc",
					),
					resource.TestCheckResourceAttr(
						"alertlogic_deployment.foo",
						"platform_type",
						"aws",
					),
					resource.TestCheckResourceAttr(
						"alertlogic_deployment.foo",
						"mode",
						"manual",
					),
				),
			},
		},
	})
}

const testAccDefaultAWSDeployment = `
resource "alertlogic_deployment" "foo" {
  name          = "terraform-testacc"
  platform_type = "aws"
  mode          = "manual"
}
`

func testAccCheckDeploymentExists(n string, v *deployments.Deployment) resource.TestCheckFunc {
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

		conn := testAccProvider.Meta().(*ALClient).deploymentsconn
		resp, err := conn.GetDeployment(&deployments.GetDeploymentInput{
			AccountId: alertlogic.String(accId),
			Id:        alertlogic.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		*v = *resp.Deployment

		return nil
	}
}
