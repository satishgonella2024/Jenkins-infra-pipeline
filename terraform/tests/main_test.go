package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestS3Bucket(t *testing.T) {
	t.Parallel()

	// Define Terraform options
	terraformOptions := &terraform.Options{
		TerraformDir: "../", // Path to Terraform code
	}

	// Ensure resources are destroyed after tests
	defer terraform.Destroy(t, terraformOptions)

	// Run Terraform Init and Apply
	terraform.InitAndApply(t, terraformOptions)

	// Get the S3 bucket name output
	bucketName := terraform.Output(t, terraformOptions, "bucket_name")

	// Validate that the S3 bucket exists
	aws.AssertS3BucketExists(t, "eu-west-2", bucketName)

	// Assert bucket name is as expected
	assert.Equal(t, "jenkins-infra-pipeline-example-london", bucketName)
}
