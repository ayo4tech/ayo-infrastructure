package guardduty_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/guardduty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func testAccAWSGuarddutyDetectorDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { acctest.PreCheck(t) },
		ErrorCheck:                acctest.ErrorCheck(t, guardduty.EndpointsID),
		Providers:                 acctest.Providers,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsGuarddutyDetectorBasicResourceConfig(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config: testAccAwsGuarddutyDetectorBasicResourceDataConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.aws_guardduty_detector.test", "id", "aws_guardduty_detector.test", "id"),
					resource.TestCheckResourceAttr("data.aws_guardduty_detector.test", "status", "ENABLED"),
					acctest.CheckResourceAttrGlobalARN("data.aws_guardduty_detector.test", "service_role_arn", "iam", "role/aws-service-role/guardduty.amazonaws.com/AWSServiceRoleForAmazonGuardDuty"),
					resource.TestCheckResourceAttrPair("data.aws_guardduty_detector.test", "finding_publishing_frequency", "aws_guardduty_detector.test", "finding_publishing_frequency"),
				),
			},
		},
	})
}

func testAccAWSGuarddutyDetectorDataSource_Id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, guardduty.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsGuarddutyDetectorExplicitConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.aws_guardduty_detector.test", "id", "aws_guardduty_detector.test", "id"),
					resource.TestCheckResourceAttr("data.aws_guardduty_detector.test", "status", "ENABLED"),
					acctest.CheckResourceAttrGlobalARN("data.aws_guardduty_detector.test", "service_role_arn", "iam", "role/aws-service-role/guardduty.amazonaws.com/AWSServiceRoleForAmazonGuardDuty"),
					resource.TestCheckResourceAttrPair("data.aws_guardduty_detector.test", "finding_publishing_frequency", "aws_guardduty_detector.test", "finding_publishing_frequency"),
				),
			},
		},
	})
}

func testAccAwsGuarddutyDetectorBasicResourceConfig() string {
	return `
resource "aws_guardduty_detector" "test" {}
`
}

func testAccAwsGuarddutyDetectorBasicResourceDataConfig() string {
	return `
resource "aws_guardduty_detector" "test" {}

data "aws_guardduty_detector" "test" {}
`
}

func testAccAwsGuarddutyDetectorExplicitConfig() string {
	return `
resource "aws_guardduty_detector" "test" {}

data "aws_guardduty_detector" "test" {
  id = aws_guardduty_detector.test.id
}
`
}