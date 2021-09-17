package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/transfer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tftransfer "github.com/terraform-providers/terraform-provider-aws/aws/internal/service/transfer"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/transfer/finder"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/tfresource"
)

func testAccAWSTransferAccess_s3_basic(t *testing.T) {
	var conf transfer.DescribedAccess
	resourceName := "aws_transfer_access.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAWSTransfer(t)
			testAccPreCheckAWSDirectoryService(t)
			testAccPreCheckAWSDirectoryServiceSimpleDirectory(t)
		},
		ErrorCheck:   testAccErrorCheck(t, transfer.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSTransferAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSTransferAccessS3BasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSTransferAccessExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "external_id", "S-1-1-12-1234567890-123456789-1234567890-1234"),
					resource.TestCheckResourceAttr(resourceName, "home_directory", "/"+rName+"/"),
					resource.TestCheckResourceAttr(resourceName, "home_directory_type", "PATH"),
					resource.TestCheckResourceAttrSet(resourceName, "role"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"role"},
			},
			{
				Config: testAccAWSTransferAccessS3UpdatedConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSTransferAccessExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "external_id", "S-1-1-12-1234567890-123456789-1234567890-1234"),
					resource.TestCheckResourceAttr(resourceName, "home_directory", "/"+rName+"/test"),
					resource.TestCheckResourceAttr(resourceName, "home_directory_type", "PATH"),
				),
			},
		},
	})
}

func testAccAWSTransferAccess_efs_basic(t *testing.T) {
	var conf transfer.DescribedAccess
	resourceName := "aws_transfer_access.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAWSTransfer(t)
			testAccPreCheckAWSDirectoryService(t)
			testAccPreCheckAWSDirectoryServiceSimpleDirectory(t)
		},
		ErrorCheck:   testAccErrorCheck(t, transfer.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSTransferAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSTransferAccessEFSBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSTransferAccessExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "external_id", "S-1-1-12-1234567890-123456789-1234567890-1234"),
					resource.TestCheckResourceAttrSet(resourceName, "home_directory"),
					resource.TestCheckResourceAttr(resourceName, "home_directory_type", "PATH"),
					resource.TestCheckResourceAttrSet(resourceName, "role"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"role"},
			},
			{
				Config: testAccAWSTransferAccessEFSUpdatedConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSTransferAccessExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "external_id", "S-1-1-12-1234567890-123456789-1234567890-1234"),
					resource.TestCheckResourceAttrSet(resourceName, "home_directory"),
					resource.TestCheckResourceAttr(resourceName, "home_directory_type", "PATH"),
					resource.TestCheckResourceAttrSet(resourceName, "role"),
				),
			},
		},
	})
}

func testAccAWSTransferAccess_disappears(t *testing.T) {
	var conf transfer.DescribedAccess
	resourceName := "aws_transfer_access.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAWSTransfer(t)
			testAccPreCheckAWSDirectoryService(t)
			testAccPreCheckAWSDirectoryServiceSimpleDirectory(t)
		},
		ErrorCheck:   testAccErrorCheck(t, transfer.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSTransferAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSTransferAccessS3BasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSTransferAccessExists(resourceName, &conf),
					testAccCheckResourceDisappears(testAccProvider, resourceAwsTransferAccess(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAWSTransferAccess_s3_policy(t *testing.T) {
	var conf transfer.DescribedAccess
	resourceName := "aws_transfer_access.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAWSTransfer(t)
			testAccPreCheckAWSDirectoryService(t)
			testAccPreCheckAWSDirectoryServiceSimpleDirectory(t)
		},
		ErrorCheck:   testAccErrorCheck(t, transfer.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSTransferAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSTransferAccessS3ScopeDownPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSTransferAccessExists(resourceName, &conf),
					resource.TestCheckResourceAttrSet(resourceName, "policy"),
				),
			},
		},
	})
}

func testAccCheckAWSTransferAccessExists(n string, v *transfer.DescribedAccess) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Transfer Access ID is set")
		}

		serverID, externalID, err := tftransfer.AccessParseResourceID(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error parsing Transfer Access ID: %w", err)
		}

		conn := testAccProvider.Meta().(*AWSClient).transferconn

		output, err := finder.AccessByServerIDAndExternalID(conn, serverID, externalID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckAWSTransferAccessDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).transferconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_transfer_access" {
			continue
		}

		serverID, externalID, err := tftransfer.AccessParseResourceID(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error parsing Transfer Access ID: %w", err)
		}
		_, err = finder.AccessByServerIDAndExternalID(conn, serverID, externalID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("Transfer Access %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccAWSTransferAccessConfigBase(rName string) string {
	return composeConfig(testAccAvailableAZsNoOptInConfig(), fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  vpc_id                  = aws_vpc.test.id
  cidr_block              = "10.0.0.0/24"
  map_public_ip_on_launch = true
  availability_zone       = data.aws_availability_zones.available.names[0]

  tags = {
    Name = %[1]q
  }

  depends_on = [aws_internet_gateway.test]
}

resource "aws_subnet" "test2" {
  vpc_id                  = aws_vpc.test.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = data.aws_availability_zones.available.names[1]

  tags = {
    Name = %[1]q
  }

  depends_on = [aws_internet_gateway.test]
}

resource "aws_directory_service_directory" "test" {
  name     = "corp.notexample.com"
  password = "SuperSecretPassw0rd"

  vpc_settings {
    vpc_id = aws_vpc.test.id

    subnet_ids = [
      aws_subnet.test.id,
      aws_subnet.test2.id
    ]
  }
}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {
      "Service": "transfer.amazonaws.com"
    },
    "Action": "sts:AssumeRole"
  }]
}
EOF
}`, rName))
}

func testAccAWSTransferAccessConfigBase_S3(rName string) string {
	return fmt.Sprintf(`
resource "aws_transfer_server" "test" {
  identity_provider_type = "AWS_DIRECTORY_SERVICE"
  directory_id           = aws_directory_service_directory.test.id
  logging_role           = aws_iam_role.test.arn
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
  acl    = "private"
}

resource "aws_iam_role_policy" "test" {
  name = %[1]q
  role = aws_iam_role.test.id

  policy = <<POLICY
{
   "Version":"2012-10-17",
   "Statement":[
      {
         "Sid":"AllowFullAccesstoCloudWatchLogs",
         "Effect":"Allow",
         "Action":[
            "logs:*"
         ],
         "Resource":"*"
      },
      {
         "Sid":"AllowFullAccesstoS3",
         "Effect":"Allow",
         "Action":[
            "s3:*"
         ],
         "Resource":"*"
      }
   ]
}
POLICY
}
`, rName)
}

func testAccAWSTransferAccessS3BasicConfig(rName string) string {
	return composeConfig(
		testAccAWSTransferAccessConfigBase(rName),
		testAccAWSTransferAccessConfigBase_S3(rName),
		`
resource "aws_transfer_access" "test" {
  external_id = "S-1-1-12-1234567890-123456789-1234567890-1234"
  server_id   = aws_transfer_server.test.id
  role        = aws_iam_role.test.arn

  home_directory      = "/${aws_s3_bucket.test.id}/"
  home_directory_type = "PATH"
}
`)
}

func testAccAWSTransferAccessS3UpdatedConfig(rName string) string {
	return composeConfig(
		testAccAWSTransferAccessConfigBase(rName),
		testAccAWSTransferAccessConfigBase_S3(rName),
		`
resource "aws_transfer_access" "test" {
  external_id = "S-1-1-12-1234567890-123456789-1234567890-1234"
  server_id   = aws_transfer_server.test.id
  role        = aws_iam_role.test.arn

  home_directory      = "/${aws_s3_bucket.test.id}/test"
  home_directory_type = "PATH"
}
`)
}

func testAccAWSTransferAccessS3ScopeDownPolicyConfig(rName string) string {
	return composeConfig(
		testAccAWSTransferAccessConfigBase(rName),
		testAccAWSTransferAccessConfigBase_S3(rName),
		`
resource "aws_transfer_access" "test" {
  external_id = "S-1-1-12-1234567890-123456789-1234567890-1234"
  server_id   = aws_transfer_server.test.id
  role        = aws_iam_role.test.arn

  home_directory      = "/${aws_s3_bucket.test.id}/"
  home_directory_type = "PATH"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowListingOfUserFolder",
            "Action": [
                "s3:ListBucket"
            ],
            "Effect": "Allow",
            "Resource": [
                "arn:${data.aws_partition.current.partition}:s3:::$${transfer:HomeBucket}"
            ]
        },
        {
            "Sid": "HomeDirObjectAccess",
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:GetObject",
                "s3:DeleteObject",
                "s3:DeleteObjectVersion",
                "s3:GetObjectVersion",
                "s3:GetObjectACL",
                "s3:PutObjectACL"
            ],
            "Resource": "arn:${data.aws_partition.current.partition}:s3:::$${transfer:HomeDirectory}/*"
        }
    ]
}
EOF
}`)
}

func testAccAWSTransferAccessConfigBase_efs(rName string) string {
	return fmt.Sprintf(`
resource "aws_transfer_server" "test" {
  identity_provider_type = "AWS_DIRECTORY_SERVICE"
  directory_id           = aws_directory_service_directory.test.id
  logging_role           = aws_iam_role.test.arn
  domain                 = "EFS"
}

resource "aws_efs_file_system" "test" {}

resource "aws_iam_role_policy" "test" {
  name = %[1]q
  role = aws_iam_role.test.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [{
    "Sid": "AllowFullAccesstoCloudWatchLogs",
    "Effect": "Allow",
    "Action": [
      "logs:*"
    ],
    "Resource": "*"
  }]
}
POLICY
}
`, rName)
}

func testAccAWSTransferAccessEFSBasicConfig(rName string) string {
	return composeConfig(
		testAccAWSTransferAccessConfigBase(rName),
		testAccAWSTransferAccessConfigBase_efs(rName),
		`
resource "aws_transfer_access" "test" {
  external_id = "S-1-1-12-1234567890-123456789-1234567890-1234"
  server_id   = aws_transfer_server.test.id
  role        = aws_iam_role.test.arn

  home_directory      = "/${aws_efs_file_system.test.id}/"
  home_directory_type = "PATH"

  posix_profile {
    gid = 1000
    uid = 1000
  }
}
`)
}

func testAccAWSTransferAccessEFSUpdatedConfig(rName string) string {
	return composeConfig(
		testAccAWSTransferAccessConfigBase(rName),
		testAccAWSTransferAccessConfigBase_efs(rName),
		`
resource "aws_transfer_access" "test" {
  external_id = "S-1-1-12-1234567890-123456789-1234567890-1234"
  server_id   = aws_transfer_server.test.id
  role        = aws_iam_role.test.arn

  home_directory      = "/${aws_efs_file_system.test.id}/test"
  home_directory_type = "PATH"

  posix_profile {
    gid = 1000
    uid = 1000
  }
}
`)
}