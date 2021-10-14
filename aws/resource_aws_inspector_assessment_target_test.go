package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/inspector"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
)

func TestAccAWSInspectorTarget_basic(t *testing.T) {
	var assessmentTarget1 inspector.AssessmentTarget
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_inspector_assessment_target.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, inspector.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSInspectorTargetAssessmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSInspectorTargetAssessmentConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSInspectorTargetExists(resourceName, &assessmentTarget1),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "inspector", regexp.MustCompile(`target/.+`)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "resource_group_arn", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSInspectorTarget_disappears(t *testing.T) {
	var assessmentTarget1 inspector.AssessmentTarget
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_inspector_assessment_target.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, inspector.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSInspectorTargetAssessmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSInspectorTargetAssessmentConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSInspectorTargetExists(resourceName, &assessmentTarget1),
					testAccCheckAWSInspectorTargetDisappears(&assessmentTarget1),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSInspectorTarget_Name(t *testing.T) {
	var assessmentTarget1, assessmentTarget2 inspector.AssessmentTarget
	rName1 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_inspector_assessment_target.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, inspector.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSInspectorTargetAssessmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSInspectorTargetAssessmentConfig(rName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSInspectorTargetExists(resourceName, &assessmentTarget1),
					resource.TestCheckResourceAttr(resourceName, "name", rName1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSInspectorTargetAssessmentConfig(rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSInspectorTargetExists(resourceName, &assessmentTarget2),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
				),
			},
		},
	})
}

func TestAccAWSInspectorTarget_ResourceGroupArn(t *testing.T) {
	var assessmentTarget1, assessmentTarget2, assessmentTarget3, assessmentTarget4 inspector.AssessmentTarget
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	inspectorResourceGroupResourceName1 := "aws_inspector_resource_group.test1"
	inspectorResourceGroupResourceName2 := "aws_inspector_resource_group.test2"
	resourceName := "aws_inspector_assessment_target.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, inspector.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSInspectorTargetAssessmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSInspectorTargetAssessmentConfigResourceGroupArn(rName, inspectorResourceGroupResourceName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSInspectorTargetExists(resourceName, &assessmentTarget1),
					resource.TestCheckResourceAttrPair(resourceName, "resource_group_arn", inspectorResourceGroupResourceName1, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSInspectorTargetAssessmentConfigResourceGroupArn(rName, inspectorResourceGroupResourceName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSInspectorTargetExists(resourceName, &assessmentTarget2),
					resource.TestCheckResourceAttrPair(resourceName, "resource_group_arn", inspectorResourceGroupResourceName2, "arn"),
				),
			},
			{
				Config: testAccAWSInspectorTargetAssessmentConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSInspectorTargetExists(resourceName, &assessmentTarget3),
					resource.TestCheckResourceAttr(resourceName, "resource_group_arn", ""),
				),
			},
			{
				Config: testAccAWSInspectorTargetAssessmentConfigResourceGroupArn(rName, inspectorResourceGroupResourceName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSInspectorTargetExists(resourceName, &assessmentTarget4),
					resource.TestCheckResourceAttrPair(resourceName, "resource_group_arn", inspectorResourceGroupResourceName1, "arn"),
				),
			},
		},
	})
}

func testAccCheckAWSInspectorTargetAssessmentDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).InspectorConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_inspector_assessment_target" {
			continue
		}

		assessmentTarget, err := describeInspectorAssessmentTarget(conn, rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error finding Inspector Assessment Target: %s", err)
		}

		if assessmentTarget != nil {
			return fmt.Errorf("Inspector Assessment Target (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckAWSInspectorTargetExists(name string, target *inspector.AssessmentTarget) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).InspectorConn

		assessmentTarget, err := describeInspectorAssessmentTarget(conn, rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("Error finding Inspector Assessment Target: %s", err)
		}

		if assessmentTarget == nil {
			return fmt.Errorf("Inspector Assessment Target (%s) not found", rs.Primary.ID)
		}

		*target = *assessmentTarget

		return nil
	}
}

func testAccCheckAWSInspectorTargetDisappears(assessmentTarget *inspector.AssessmentTarget) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).InspectorConn

		input := &inspector.DeleteAssessmentTargetInput{
			AssessmentTargetArn: assessmentTarget.Arn,
		}

		_, err := conn.DeleteAssessmentTarget(input)

		return err
	}
}

func testAccAWSInspectorTargetAssessmentConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_inspector_assessment_target" "test" {
  name = %q
}
`, rName)
}

func testAccAWSInspectorTargetAssessmentConfigResourceGroupArn(rName, inspectorResourceGroupResourceName string) string {
	return fmt.Sprintf(`
resource "aws_inspector_resource_group" "test1" {
  tags = {
    Name = "%s1"
  }
}

resource "aws_inspector_resource_group" "test2" {
  tags = {
    Name = "%s2"
  }
}

resource "aws_inspector_assessment_target" "test" {
  name               = %q
  resource_group_arn = %s.arn
}
`, rName, rName, rName, inspectorResourceGroupResourceName)
}
