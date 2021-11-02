package kafka_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kafka"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfkafka "github.com/hashicorp/terraform-provider-aws/internal/service/kafka"
)

func TestAccKafkaConfiguration_basic(t *testing.T) {
	var configuration1 kafka.DescribeConfigurationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kafka_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, kafka.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckMskConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMskConfigurationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMskConfigurationExists(resourceName, &configuration1),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "kafka", regexp.MustCompile(`configuration/.+`)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "kafka_versions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "latest_revision", "1"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestMatchResourceAttr(resourceName, "server_properties", regexp.MustCompile(`auto.create.topics.enable = true`)),
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

func TestAccKafkaConfiguration_disappears(t *testing.T) {
	var configuration1 kafka.DescribeConfigurationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kafka_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, kafka.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckMskConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMskConfigurationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMskConfigurationExists(resourceName, &configuration1),
					acctest.CheckResourceDisappears(acctest.Provider, tfkafka.ResourceConfiguration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccKafkaConfiguration_description(t *testing.T) {
	var configuration1, configuration2 kafka.DescribeConfigurationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kafka_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, kafka.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckMskConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMskConfigurationConfigDescription(rName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMskConfigurationExists(resourceName, &configuration1),
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMskConfigurationConfigDescription(rName, "description2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMskConfigurationExists(resourceName, &configuration2),
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
					resource.TestCheckResourceAttr(resourceName, "latest_revision", "2"),
				),
			},
		},
	})
}

func TestAccKafkaConfiguration_kafkaVersions(t *testing.T) {
	var configuration1 kafka.DescribeConfigurationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kafka_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, kafka.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckMskConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMskConfigurationConfigKafkaVersions(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMskConfigurationExists(resourceName, &configuration1),
					resource.TestCheckResourceAttr(resourceName, "kafka_versions.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "kafka_versions.*", "2.6.0"),
					resource.TestCheckTypeSetElemAttr(resourceName, "kafka_versions.*", "2.7.0"),
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

func TestAccKafkaConfiguration_serverProperties(t *testing.T) {
	var configuration1, configuration2 kafka.DescribeConfigurationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kafka_configuration.test"
	serverProperty1 := "auto.create.topics.enable = false"
	serverProperty2 := "auto.create.topics.enable = true"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, kafka.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckMskConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMskConfigurationConfigServerProperties(rName, serverProperty1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMskConfigurationExists(resourceName, &configuration1),
					resource.TestMatchResourceAttr(resourceName, "server_properties", regexp.MustCompile(serverProperty1)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMskConfigurationConfigServerProperties(rName, serverProperty2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMskConfigurationExists(resourceName, &configuration2),
					resource.TestCheckResourceAttr(resourceName, "latest_revision", "2"),
					resource.TestMatchResourceAttr(resourceName, "server_properties", regexp.MustCompile(serverProperty2)),
				),
			},
		},
	})
}

func testAccCheckMskConfigurationDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).KafkaConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_kafka_configuration" {
			continue
		}

		input := &kafka.DescribeConfigurationInput{
			Arn: aws.String(rs.Primary.ID),
		}

		output, err := conn.DescribeConfiguration(input)

		if tfawserr.ErrMessageContains(err, kafka.ErrCodeBadRequestException, "Configuration ARN does not exist") {
			continue
		}

		if err != nil {
			return err
		}

		if output != nil {
			return fmt.Errorf("MSK Configuration (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMskConfigurationExists(resourceName string, configuration *kafka.DescribeConfigurationOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource ID not set: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).KafkaConn

		input := &kafka.DescribeConfigurationInput{
			Arn: aws.String(rs.Primary.ID),
		}

		output, err := conn.DescribeConfiguration(input)

		if err != nil {
			return fmt.Errorf("error describing MSK Cluster (%s): %s", rs.Primary.ID, err)
		}

		*configuration = *output

		return nil
	}
}

func testAccMskConfigurationConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_kafka_configuration" "test" {
  name = %[1]q

  server_properties = <<PROPERTIES
auto.create.topics.enable = true
delete.topic.enable = true
PROPERTIES
}
`, rName)
}

func testAccMskConfigurationConfigDescription(rName, description string) string {
	return fmt.Sprintf(`
resource "aws_kafka_configuration" "test" {
  description = %[2]q
  name        = %[1]q

  server_properties = <<PROPERTIES
auto.create.topics.enable = true
PROPERTIES
}
`, rName, description)
}

func testAccMskConfigurationConfigKafkaVersions(rName string) string {
	return fmt.Sprintf(`
resource "aws_kafka_configuration" "test" {
  kafka_versions = ["2.6.0", "2.7.0"]
  name           = %[1]q

  server_properties = <<PROPERTIES
auto.create.topics.enable = true
PROPERTIES
}
`, rName)
}

func testAccMskConfigurationConfigServerProperties(rName string, serverProperty string) string {
	return fmt.Sprintf(`
resource "aws_kafka_configuration" "test" {
  name = %[1]q

  server_properties = <<PROPERTIES
%[2]s
PROPERTIES
}
`, rName, serverProperty)
}
